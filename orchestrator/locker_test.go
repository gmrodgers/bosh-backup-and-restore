package orchestrator_test

import (
	"fmt"

	"github.com/cloudfoundry-incubator/bosh-backup-and-restore/executor"
	"github.com/cloudfoundry-incubator/bosh-backup-and-restore/orchestrator"
	"github.com/cloudfoundry-incubator/bosh-backup-and-restore/orchestrator/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Backup", func() {
	var (
		b                 *orchestrator.Locker
		deployment        *fakes.FakeDeployment
		deploymentManager *fakes.FakeDeploymentManager
		logger            *fakes.FakeLogger
		lockOrderer       *fakes.FakeLockOrderer
		deploymentName    = "foobarbaz"
		actualLockError   error
	)

	BeforeEach(func() {
		deployment = new(fakes.FakeDeployment)
		deploymentManager = new(fakes.FakeDeploymentManager)
		logger = new(fakes.FakeLogger)

		b = orchestrator.NewLocker(logger, deploymentManager, lockOrderer, executor.NewParallelExecutor())
	})

	JustBeforeEach(func() {
		actualLockError = b.Lock(deploymentName)
	})

	Context("locks up a deployment", func() {
		BeforeEach(func() {
			deploymentManager.FindReturns(deployment, nil)
			deployment.IsBackupableReturns(true)
		})

		It("does not fail", func() {
			Expect(actualLockError).NotTo(HaveOccurred())
		})

		It("finds the deployment", func() {
			Expect(deploymentManager.FindCallCount()).To(Equal(1))
			Expect(deploymentManager.FindArgsForCall(0)).To(Equal(deploymentName))
		})

		It("checks if the deployment is backupable", func() {
			Expect(deployment.IsBackupableCallCount()).To(Equal(1))
		})

		It("runs pre-backup-lock scripts on the deployment", func() {
			Expect(deployment.PreBackupLockCallCount()).To(Equal(1))
		})
	})

	Describe("failures", func() {
		var expectedError = fmt.Errorf("Profanity")

		Context("fails to find deployment", func() {
			BeforeEach(func() {
				deploymentManager.FindReturns(nil, expectedError)
			})

			It("fails the backup process", func() {
				expectErrorMatch(actualLockError, expectedError)
			})
		})

		Context("fails if the deployment is not backupable", func() {
			BeforeEach(func() {
				deploymentManager.FindReturns(deployment, nil)
				deployment.IsBackupableReturns(false)
			})

			It("finds a deployment with the deployment name", func() {
				Expect(deploymentManager.FindCallCount()).To(Equal(1))
				Expect(deploymentManager.FindArgsForCall(0)).To(Equal(deploymentName))
			})

			It("checks if the deployment is backupable", func() {
				Expect(deployment.IsBackupableCallCount()).To(Equal(1))
			})

			It("fails the backup process", func() {
				Expect(actualLockError).To(ConsistOf(MatchError("Deployment '" + deploymentName + "' has no backup scripts")))
			})
		})

		Context("fails if pre-backup-lock fails", func() {
			var lockError = orchestrator.NewLockError("smoooooooth jazz")

			BeforeEach(func() {
				deploymentManager.FindReturns(deployment, nil)
				deployment.IsBackupableReturns(true)
				deployment.PreBackupLockReturns(lockError)
			})

			It("fails the lock process", func() {
				expectErrorMatch(actualLockError, lockError)
			})
		})
	})
})

func expectErrorMatch(actual error, expected ...error) {
	if actualErrors, isErrorList := actual.(orchestrator.Error); isErrorList {
		for _, err := range actualErrors {
			Expect(actual).To(MatchError(ContainSubstring(err.Error())))
		}
		Expect(len(actualErrors)).To(Equal(len(expected)))
	} else {
		Expect(actual).To(MatchError(expected))
	}
}
