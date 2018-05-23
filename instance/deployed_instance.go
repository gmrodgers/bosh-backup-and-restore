package instance

import (
	"fmt"

	"github.com/cloudfoundry-incubator/bosh-backup-and-restore/orchestrator"
	"github.com/cloudfoundry-incubator/bosh-backup-and-restore/ssh"
	"github.com/pkg/errors"
)

type Logger interface {
	Debug(tag, msg string, args ...interface{})
	Info(tag, msg string, args ...interface{})
	Warn(tag, msg string, args ...interface{})
	Error(tag, msg string, args ...interface{})
}

type DeployedInstance struct {
	backupAndRestoreInstanceIndex string
	instanceID                    string
	instanceGroupName             string
	artifactDirCreated            bool
	Logger                        Logger
	jobs                          orchestrator.Jobs
	remoteRunner                  ssh.RemoteRunner
}

func NewDeployedInstance(instanceIndex string, instanceGroupName string, instanceID string, remoteRunner ssh.RemoteRunner, logger Logger, jobs orchestrator.Jobs) *DeployedInstance {
	return &DeployedInstance{
		backupAndRestoreInstanceIndex: instanceIndex,
		instanceGroupName:             instanceGroupName,
		instanceID:                    instanceID,
		Logger:                        logger,
		jobs:                          jobs,
		remoteRunner:                  remoteRunner,
	}
}

func (i *DeployedInstance) IsBackupable() bool {
	return i.jobs.AnyAreBackupable()
}

func (i *DeployedInstance) Jobs() []orchestrator.Job {
	return i.jobs
}

func (i *DeployedInstance) Backup() error {
	var backupErrors []error
	//for _, job := range i.jobs {
	//	if err := job.Backup(); err != nil {
	//		backupErrors = append(backupErrors, err)
	//	}
	//}

	if i.IsBackupable() {
		//i.artifactDirCreated = true
	}

	return orchestrator.ConvertErrors(backupErrors)
}

func (i *DeployedInstance) Name() string {
	return i.instanceGroupName
}

func (i *DeployedInstance) Index() string {
	return i.backupAndRestoreInstanceIndex
}

func (i *DeployedInstance) ID() string {
	return i.instanceID
}

func (i *DeployedInstance) ConnectedUsername() string {
	return i.remoteRunner.ConnectedUsername()
}

func (i *DeployedInstance) handleErrs(jobName, label string, err error, exitCode int, stdout, stderr []byte) error {
	var foundErrors []error

	if err != nil {
		i.Logger.Error("bbr", fmt.Sprintf(
			"Error attempting to run %s script for job %s on %s/%s. Error: %s",
			label,
			jobName,
			i.instanceGroupName,
			i.instanceID,
			err.Error(),
		))
		foundErrors = append(foundErrors, err)
	}

	if exitCode != 0 {
		errorString := fmt.Sprintf(
			"%s script for job %s failed on %s/%s.\nStdout: %s\nStderr: %s",
			label,
			jobName,
			i.instanceGroupName,
			i.instanceID,
			stdout,
			stderr,
		)

		foundErrors = append(foundErrors, errors.New(errorString))

		i.Logger.Error("bbr", errorString)
	}

	return orchestrator.ConvertErrors(foundErrors)
}
