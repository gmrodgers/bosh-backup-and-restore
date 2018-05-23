package orchestrator_test

import (
	"fmt"
	"github.com/cloudfoundry-incubator/bosh-backup-and-restore/executor"
	"github.com/cloudfoundry-incubator/bosh-backup-and-restore/orchestrator"
	"github.com/cloudfoundry-incubator/bosh-backup-and-restore/orchestrator/fakes"
	executorFakes "github.com/cloudfoundry-incubator/bosh-backup-and-restore/executor/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Deployment", func() {
	var (
		deployment orchestrator.Deployment
		logger     *fakes.FakeLogger

		instances []orchestrator.Instance
		instance1 *fakes.FakeInstance
		instance2 *fakes.FakeInstance
		instance3 *fakes.FakeInstance

		job1a *fakes.FakeJob
		job1b *fakes.FakeJob
		job2a *fakes.FakeJob
		job3a *fakes.FakeJob
	)

	BeforeEach(func() {
		logger = new(fakes.FakeLogger)

		instance1 = new(fakes.FakeInstance)
		instance2 = new(fakes.FakeInstance)
		instance3 = new(fakes.FakeInstance)

		job1a = new(fakes.FakeJob)
		job1b = new(fakes.FakeJob)
		job2a = new(fakes.FakeJob)
		job3a = new(fakes.FakeJob)

		instance1.JobsReturns([]orchestrator.Job{job1a, job1b})
		instance2.JobsReturns([]orchestrator.Job{job2a})
		instance3.JobsReturns([]orchestrator.Job{job3a})
	})

	JustBeforeEach(func() {
		deployment = orchestrator.NewDeployment(logger, instances)
	})

	Context("PreBackupLock", func() {
		var (
			lockError    error
			lockOrderer  *fakes.FakeLockOrderer
			fakeExecutor *executorFakes.FakeExecutor
		)

		BeforeEach(func() {
			lockOrderer = new(fakes.FakeLockOrderer)
			fakeExecutor = new(executorFakes.FakeExecutor)
			instances = []orchestrator.Instance{instance1, instance2, instance3}
			lockOrderer.OrderReturns([][]orchestrator.Job{{job2a}, {job3a, job1a}, {job1b}}, nil)
		})

		JustBeforeEach(func() {
			lockError = deployment.PreBackupLock(lockOrderer, fakeExecutor)
		})

		It("delegates the execution to the executor", func() {
			Expect(lockError).NotTo(HaveOccurred())
			Expect(lockOrderer.OrderArgsForCall(0)).To(ConsistOf(job1a, job1b, job2a, job3a))
			Expect(fakeExecutor.RunArgsForCall(0)).To(Equal([][]executor.Executable{
				{orchestrator.NewJobPreBackupLockExecutable(job2a)},
				{orchestrator.NewJobPreBackupLockExecutable(job3a), orchestrator.NewJobPreBackupLockExecutable(job1a)},
				{orchestrator.NewJobPreBackupLockExecutable(job1b)},
			}))
		})

		Context("if the pre-backup-lock fails", func() {
			BeforeEach(func() {
				fakeExecutor.RunReturns([]error{
					fmt.Errorf("job1b failed"),
					fmt.Errorf("job2a failed"),
				})
			})

			It("fails", func() {
				Expect(lockError).To(MatchError(SatisfyAll(
					ContainSubstring("job1b failed"),
					ContainSubstring("job2a failed"),
				)))
			})
		})

		Context("if the lockOrderer returns an error", func() {
			BeforeEach(func() {
				lockOrderer.OrderReturns(nil, fmt.Errorf("test lock orderer error"))
			})

			It("fails", func() {
				Expect(lockError).To(MatchError(ContainSubstring("test lock orderer error")))
			})
		})
	})

	Context("PostBackupUnlock", func() {
		var (
			lockError    error
			lockOrderer  *fakes.FakeLockOrderer
			fakeExecutor *executorFakes.FakeExecutor
		)

		BeforeEach(func() {
			lockOrderer = new(fakes.FakeLockOrderer)
			fakeExecutor = new(executorFakes.FakeExecutor)
			instances = []orchestrator.Instance{instance1, instance2, instance3}
			lockOrderer.OrderReturns([][]orchestrator.Job{{job2a}, {job3a, job1a}, {job1b}}, nil)
		})

		JustBeforeEach(func() {
			lockError = deployment.PostBackupUnlock(lockOrderer, fakeExecutor)
		})

		It("delegates the execution to the executor", func() {
			Expect(lockError).NotTo(HaveOccurred())
			Expect(lockOrderer.OrderArgsForCall(0)).To(ConsistOf(job1a, job1b, job2a, job3a))
			Expect(fakeExecutor.RunArgsForCall(0)).To(Equal([][]executor.Executable{
				{orchestrator.NewJobPostBackupUnlockExecutable(job2a)},
				{orchestrator.NewJobPostBackupUnlockExecutable(job3a), orchestrator.NewJobPostBackupUnlockExecutable(job1a)},
				{orchestrator.NewJobPostBackupUnlockExecutable(job1b)},
			}))
		})

		Context("if the post-backup-unlock fails", func() {
			BeforeEach(func() {
				fakeExecutor.RunReturns([]error{
					fmt.Errorf("job1b failed"),
					fmt.Errorf("job2a failed"),
				})
			})

			It("fails", func() {
				Expect(lockError).To(MatchError(SatisfyAll(
					ContainSubstring("job1b failed"),
					ContainSubstring("job2a failed"),
				)))
			})
		})

		Context("if the lockOrderer returns an error", func() {
			BeforeEach(func() {
				lockOrderer.OrderReturns(nil, fmt.Errorf("test lock orderer error"))
			})

			It("fails", func() {
				Expect(lockError).To(MatchError(ContainSubstring("test lock orderer error")))
			})
		})
	})

	Context("IsBackupable", func() {
		Context("when at least one instance is backupable", func() {
			BeforeEach(func() {
				instance1.IsBackupableReturns(false)
				instance2.IsBackupableReturns(true)
				instances = []orchestrator.Instance{instance1, instance2}
			})

			It("returns true", func() {
				Expect(deployment.IsBackupable()).To(BeTrue())
			})
		})

		Context("when no instances are backupable", func() {
			BeforeEach(func() {
				instance1.IsBackupableReturns(false)
				instance2.IsBackupableReturns(false)
				instances = []orchestrator.Instance{instance1, instance2}
			})

			It("returns false", func() {
				Expect(deployment.IsBackupable()).To(BeFalse())
			})
		})
	})

	Context("BackupableInstances", func() {
		BeforeEach(func() {
			instance1.IsBackupableReturns(true)
			instance2.IsBackupableReturns(false)
			instance3.IsBackupableReturns(true)
			instances = []orchestrator.Instance{instance1, instance2, instance3}
		})

		It("returns a list of all backupable instances", func() {
			Expect(deployment.BackupableInstances()).To(ConsistOf(instance1, instance3))
		})
	})

	Context("Instances", func() {
		BeforeEach(func() {
			instances = []orchestrator.Instance{instance1, instance2, instance3}
		})

		It("returns instances for the deployment", func() {
			Expect(deployment.Instances()).To(ConsistOf(instance1, instance2, instance3))
		})
	})
})
