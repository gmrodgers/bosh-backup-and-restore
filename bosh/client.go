package bosh

import (
	"strconv"

	"github.com/cloudfoundry-incubator/bosh-backup-and-restore/instance"
	"github.com/cloudfoundry-incubator/bosh-backup-and-restore/orchestrator"
	"github.com/cloudfoundry-incubator/bosh-backup-and-restore/ssh"
	"github.com/cloudfoundry/bosh-cli/director"
	"github.com/cloudfoundry/bosh-utils/uuid"

	"fmt"

	"github.com/pkg/errors"
	gossh "golang.org/x/crypto/ssh"
)

//go:generate counterfeiter -o fakes/fake_bosh_client.go . BoshClient
type BoshClient interface {
	FindInstances(deploymentName string) ([]orchestrator.Instance, error)
	GetManifest(deploymentName string) (string, error)
}

func NewClient(boshDirector director.Director,
	sshOptsGenerator ssh.SSHOptsGenerator,
	remoteRunnerFactory ssh.RemoteRunnerFactory,
	logger Logger,
	jobFinder instance.JobFinder,
	manifestQuerierCreator instance.ManifestQuerierCreator) Client {
	return Client{
		Director:               boshDirector,
		SSHOptsGenerator:       sshOptsGenerator,
		RemoteRunnerFactory:    remoteRunnerFactory,
		Logger:                 logger,
		jobFinder:              jobFinder,
		manifestQuerierCreator: manifestQuerierCreator,
	}
}

type Client struct {
	director.Director
	ssh.SSHOptsGenerator
	ssh.RemoteRunnerFactory
	Logger
	jobFinder              instance.JobFinder
	manifestQuerierCreator instance.ManifestQuerierCreator
}

//go:generate counterfeiter -o fakes/fake_logger.go . Logger
type Logger interface {
	Debug(tag, msg string, args ...interface{})
	Info(tag, msg string, args ...interface{})
	Warn(tag, msg string, args ...interface{})
	Error(tag, msg string, args ...interface{})
}

func (c Client) FindInstances(deploymentName string) ([]orchestrator.Instance, error) {
	deployment, err := c.Director.FindDeployment(deploymentName)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't find deployment "+deploymentName)
	}

	c.Logger.Debug("bbr", "Finding VMs...")

	vms, err := deployment.VMInfos()
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get vm infos")
	}

	var instances []orchestrator.Instance
	var instancesConnectedTo []director.AllOrInstanceGroupOrInstanceSlug

	manifest, err := deployment.Manifest()
	if err != nil {
		return nil, errors.Wrap(err, "couldn't find manifest for deployment "+deploymentName)
	}

	manifestQuerier, err := c.manifestQuerierCreator(manifest)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't generate manifest querier for deployment "+deploymentName)
	}

	sshOpts, privateKey, err := c.SSHOptsGenerator(uuid.NewGenerator())
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate ssh options")
	}

	c.Logger.Debug("bbr", "SSH user generated: %s", sshOpts.Username)

	for _, instanceGroupName := range uniqueInstanceGroupNamesFromVMs(vms) {
		c.Logger.Debug("bbr", "Setting up SSH for job %s", instanceGroupName)

		instanceGroupInstances, err := director.NewAllOrInstanceGroupOrInstanceSlugFromString(instanceGroupName)
		if err != nil {
			cleanupAlreadyMadeConnections(deployment, instancesConnectedTo, sshOpts)
			return nil, errors.Wrap(err, "invalid instance group name: "+instanceGroupName)
		}

		sshRes, err := deployment.SetUpSSH(instanceGroupInstances, sshOpts)
		if err != nil {
			cleanupAlreadyMadeConnections(deployment, instancesConnectedTo, sshOpts)
			return nil, errors.Wrap(err, "failed to set up ssh")
		}

		instancesConnectedTo = append(instancesConnectedTo, instanceGroupInstances)

		for _, host := range sshRes.Hosts {
			instance, hasNoJobs, isWindows, err := c.getInstance(
				instance.InstanceIdentifier{
					InstanceGroupName: instanceGroupName,
					InstanceId:        host.IndexOrID,
					Bootstrap:         isInstanceABootstrapNode(instanceGroupName, host.Host, vms),
				},
				host,
				privateKey,
				deployment,
				instanceGroupName,
				vms,
				manifestQuerier,
				func() { cleanupAlreadyMadeConnections(deployment, instancesConnectedTo, sshOpts) },
			)

			if err != nil {
				return nil, err
			}

			if isWindows {
				c.Logger.Warn("bbr", "skipping Windows instance %s/%s", instanceGroupName, host.IndexOrID)
				continue
			}

			instances = append(instances, instance)

			if hasNoJobs {
				c.Logger.Debug("bbr", "no scripts found on instance %s/%s, skipping rest of the instances for %s", instanceGroupName, host.IndexOrID, instanceGroupName)
				break
			}
		}
	}

	return instances, nil
}

func (c Client) getInstance(
	instanceIdentifier instance.InstanceIdentifier,
	host director.Host,
	privateKey string,
	deployment director.Deployment,
	instanceGroupName string,
	vms []director.VMInfo,
	manifestQuerier instance.ManifestQuerier,
	cleanUpFunc func(),
) (orchestrator.Instance, bool, bool, error) {
	c.Logger.Debug("bbr", "Attempting to SSH onto %s, %s", host.Host, host.IndexOrID)

	remoteRunner, err := c.createSSHRemoteRunner(host.Host, host.Username, host.HostPublicKey, privateKey)
	if err != nil {
		cleanUpFunc()
		return nil, false, false, errors.Wrap(err, "failed to connect using ssh")
	}

	isWindows, err := remoteRunner.IsWindows()
	if err != nil {
		cleanUpFunc()
		return nil, false, true, errors.Wrap(err, "failed to check os")
	}

	if isWindows {
		return nil, false, true, nil
	}

	jobs, err := c.jobFinder.FindJobs(
		instanceIdentifier,
		remoteRunner,
		manifestQuerier,
	)
	if err != nil {
		cleanUpFunc()
		return nil, false, false, errors.Wrap(err, "couldn't find jobs")
	}

	vmIndex, err := findInstanceIndexByID(vms, host.IndexOrID)
	if err != nil {
		cleanUpFunc()
		return nil, false, false, errors.Wrap(err, "couldn't find instance index")
	}

	return NewBoshDeployedInstance(
		instanceGroupName,
		vmIndex,
		host.IndexOrID,
		remoteRunner,
		deployment,
		false,
		c.Logger,
		jobs,
	), len(jobs) == 0, false, nil
}

func (c Client) createSSHRemoteRunner(host, username, publicKey, privateKey string) (ssh.RemoteRunner, error) {
	hostPublicKey, _, _, _, err := gossh.ParseAuthorizedKey([]byte(publicKey))
	if err != nil {
		return nil, errors.Wrap(err, "ssh.NewConnection.ParseAuthorizedKey failed")
	}

	remoteRunner, err := c.RemoteRunnerFactory(
		host,
		username,
		privateKey,
		gossh.FixedHostKey(hostPublicKey),
		[]string{hostPublicKey.Type()},
		c.Logger,
	)

	return remoteRunner, err
}

func isInstanceABootstrapNode(jobName, ip string, vms []director.VMInfo) bool {
	for _, vm := range vms {
		if jobName != vm.JobName || !vm.Bootstrap {
			continue
		}

		for _, vmIP := range vm.IPs {
			if ip == vmIP {
				return true
			}
		}
	}

	return false
}

func (c Client) GetManifest(deploymentName string) (string, error) {
	deployment, err := c.Director.FindDeployment(deploymentName)
	if err != nil {
		return "", errors.Wrap(err, "couldn't find deployment "+deploymentName)
	}

	return deployment.Manifest()
}

type JobVMInfo struct {
	JobName string
	VMInfo  director.VMInfo
}

func uniqueInstanceGroupNamesFromVMs(vms []director.VMInfo) []string {
	var jobs []string

	for _, vm := range vms {
		if !contains(jobs, vm.JobName) {
			jobs = append(jobs, vm.JobName)
		}
	}

	return jobs
}

func findInstanceIndexByID(vmInfos []director.VMInfo, vmID string) (string, error) {
	for _, vmInfo := range vmInfos {
		if vmID == vmInfo.ID {
			if vmInfo.Index == nil {
				return "", fmt.Errorf("vmInfo index is nil, given vmID %s", vmID)
			}

			return strconv.Itoa(*vmInfo.Index), nil
		}
	}

	return "", fmt.Errorf("vmInfo does not contain given vmID %s", vmID)
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}

	return false
}

func cleanupAlreadyMadeConnections(
	deployment director.Deployment,
	slugs []director.AllOrInstanceGroupOrInstanceSlug,
	opts director.SSHOpts,
) {
	for _, slug := range slugs {
		deployment.CleanUpSSH(slug, director.SSHOpts{Username: opts.Username})
	}
}
