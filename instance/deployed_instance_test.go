package instance_test

import (
	"log"
	"github.com/cloudfoundry-incubator/bosh-backup-and-restore/instance"
	"github.com/cloudfoundry-incubator/bosh-backup-and-restore/orchestrator"
	sshfakes "github.com/cloudfoundry-incubator/bosh-backup-and-restore/ssh/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

var _ = Describe("DeployedInstance", func() {
	var boshLogger boshlog.Logger
	var stdout, stderr *gbytes.Buffer
	var instanceGroupName, instanceIndex, instanceID, expectedStdout, expectedStderr string
	var jobs orchestrator.Jobs
	var remoteRunner *sshfakes.FakeRemoteRunner

	var deployedInstance *instance.DeployedInstance
	BeforeEach(func() {
		instanceGroupName = "instance-group-name"
		instanceIndex = "instance-index"
		instanceID = "instance-id"
		expectedStdout = "i'm a stdout"
		expectedStderr = "i'm a stderr"
		stdout = gbytes.NewBuffer()
		stderr = gbytes.NewBuffer()
		boshLogger = boshlog.New(boshlog.LevelDebug, log.New(stdout, "[bosh-package] ", log.Lshortfile), log.New(stderr, "[bosh-package] ", log.Lshortfile))
		remoteRunner = new(sshfakes.FakeRemoteRunner)
	})

	JustBeforeEach(func() {
		remoteRunner.ConnectedUsernameReturns("sshUsername")
		deployedInstance = instance.NewDeployedInstance(
			instanceIndex,
			instanceGroupName,
			instanceID,
			false,
			remoteRunner,
			boshLogger,
			jobs)
	})

	Describe("IsBackupable", func() {
		var actualBackupable bool

		JustBeforeEach(func() {
			actualBackupable = deployedInstance.IsBackupable()
		})

		Describe("there are backup scripts in the job directories", func() {
			BeforeEach(func() {
				jobs = orchestrator.Jobs([]orchestrator.Job{
					instance.NewJob(remoteRunner, instanceGroupName+"/"+instanceID, boshLogger, "", instance.BackupAndRestoreScripts{
						"/var/vcap/jobs/dave/bin/bbr/backup",
					}, instance.Metadata{}),
				})
			})

			It("returns true", func() {
				Expect(actualBackupable).To(BeTrue())
			})
		})

		Describe("there are no backup scripts in the job directories", func() {
			BeforeEach(func() {
				jobs = orchestrator.Jobs([]orchestrator.Job{
					instance.NewJob(remoteRunner, instanceGroupName+"/"+instanceID, boshLogger, "", instance.BackupAndRestoreScripts{
						"/var/vcap/jobs/dave/bin/foo",
					}, instance.Metadata{}),
				})
			})

			It("returns false", func() {
				Expect(actualBackupable).To(BeFalse())
			})
		})
	})

	Describe("Jobs", func() {
		BeforeEach(func() {
			jobs = orchestrator.Jobs([]orchestrator.Job{
				instance.NewJob(remoteRunner, instanceGroupName+"/"+instanceID, boshLogger, "", instance.BackupAndRestoreScripts{
					"/var/vcap/jobs/dave/bin/foo",
				}, instance.Metadata{}),
			})
		})

		It("returns the instance's jobs", func() {
			Expect(deployedInstance.Jobs()).To(HaveLen(1))
			Expect(deployedInstance.Jobs()[0].Name()).To(Equal("dave"))
		})
	})

	Describe("Name", func() {
		It("returns the instance name", func() {
			Expect(deployedInstance.Name()).To(Equal("instance-group-name"))
		})
	})

	Describe("Index", func() {
		It("returns the instance Index", func() {
			Expect(deployedInstance.Index()).To(Equal("instance-index"))
		})
	})
})
