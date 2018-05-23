package standalone_test

import (
	"fmt"
	"os"

	"github.com/cloudfoundry-incubator/bosh-backup-and-restore/instance"
	. "github.com/cloudfoundry-incubator/bosh-backup-and-restore/standalone"

	"io/ioutil"

	instancefakes "github.com/cloudfoundry-incubator/bosh-backup-and-restore/instance/fakes"
	"github.com/cloudfoundry-incubator/bosh-backup-and-restore/orchestrator"
	"github.com/cloudfoundry-incubator/bosh-backup-and-restore/orchestrator/fakes"
	sshfakes "github.com/cloudfoundry-incubator/bosh-backup-and-restore/ssh/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DeploymentManager", func() {
	var deploymentManager DeploymentManager
	var deploymentName = "bosh"
	var logger *fakes.FakeLogger
	var hostName = "hostname"
	var username = "username"
	var privateKey string
	var fakeJobFinder *instancefakes.FakeJobFinder
	var remoteRunnerFactory *sshfakes.FakeRemoteRunnerFactory
	var remoteRunner *sshfakes.FakeRemoteRunner

	BeforeEach(func() {
		privateKey = createTempFile("privateKey")
		logger = new(fakes.FakeLogger)
		remoteRunnerFactory = new(sshfakes.FakeRemoteRunnerFactory)
		fakeJobFinder = new(instancefakes.FakeJobFinder)
		remoteRunner = new(sshfakes.FakeRemoteRunner)

		deploymentManager = NewDeploymentManager(logger, hostName, username, privateKey, fakeJobFinder, remoteRunnerFactory.Spy)
	})

	AfterEach(func() {
		os.Remove(privateKey)
	})

	Describe("Find", func() {
		var actualDeployment orchestrator.Deployment
		var actualError error
		var fakeJobs orchestrator.Jobs

		JustBeforeEach(func() {
			actualDeployment, actualError = deploymentManager.Find(deploymentName)
		})

		Context("success", func() {
			BeforeEach(func() {
				fakeJobs = orchestrator.Jobs{instance.NewJob(nil, "", nil, "", instance.BackupAndRestoreScripts{"foo"}, instance.Metadata{})}
				remoteRunnerFactory.Returns(remoteRunner, nil)
				fakeJobFinder.FindJobsReturns(fakeJobs, nil)
			})
			It("does not fail", func() {
				Expect(actualError).NotTo(HaveOccurred())
			})

			It("invokes connection creator", func() {
				Expect(remoteRunnerFactory.CallCount()).To(Equal(1))
			})

			It("invokes job finder", func() {
				Expect(fakeJobFinder.FindJobsCallCount()).To(Equal(1))
			})

			It("returns a deployment", func() {
				Expect(actualDeployment).To(Equal(orchestrator.NewDeployment(logger, []orchestrator.Instance{
					NewDeployedInstance("bosh", remoteRunner, logger, fakeJobs),
				})))
			})
		})

		Context("can't read private key", func() {
			BeforeEach(func() {
				os.Remove(privateKey)
			})

			It("should fail", func() {
				Expect(actualError).To(MatchError(ContainSubstring("failed reading private key")))
			})

			It("should not invoke connection creator", func() {
				Expect(remoteRunnerFactory.CallCount()).To(BeZero())
			})
		})

		Context("can't create SSH connection", func() {
			connError := fmt.Errorf("error")

			BeforeEach(func() {
				remoteRunnerFactory.Returns(nil, connError)
			})

			It("should fail", func() {
				Expect(actualError).To(MatchError(connError))
			})

			It("should invoke connection creator", func() {
				Expect(remoteRunnerFactory.CallCount()).To(Equal(1))
			})

			It("should not invoke job finder", func() {
				Expect(fakeJobFinder.FindJobsCallCount()).To(BeZero())
			})

		})

		Context("can't find jobs", func() {
			findJobsErr := fmt.Errorf("error")

			BeforeEach(func() {
				remoteRunnerFactory.Returns(remoteRunner, nil)
				fakeJobFinder.FindJobsReturns(nil, findJobsErr)
			})

			It("should fail", func() {
				Expect(actualError).To(MatchError(findJobsErr))
			})

			It("should invoke connection creator", func() {
				Expect(remoteRunnerFactory.CallCount()).To(Equal(1))
			})

			It("should not invoke job finder", func() {
				Expect(fakeJobFinder.FindJobsCallCount()).To(Equal(1))
			})
		})

	})
})

func createTempFile(contents string) string {
	tempFile, err := ioutil.TempFile("", "")
	Expect(err).NotTo(HaveOccurred())
	tempFile.Write([]byte(contents))
	tempFile.Close()
	return tempFile.Name()
}
