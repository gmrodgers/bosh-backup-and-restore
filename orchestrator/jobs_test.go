package orchestrator_test

import (
	"log"

	"github.com/cloudfoundry-incubator/bosh-backup-and-restore/orchestrator"
	orchestratorFakes "github.com/cloudfoundry-incubator/bosh-backup-and-restore/orchestrator/fakes"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Jobs", func() {
	var jobs orchestrator.Jobs
	var logger boshlog.Logger

	BeforeEach(func() {
		combinedLog := log.New(GinkgoWriter, "[instance-test] ", log.Lshortfile)
		logger = boshlog.New(boshlog.LevelDebug, combinedLog, combinedLog)
	})

	Context("contains jobs with backup script", func() {
		var backupableJob *orchestratorFakes.FakeJob
		var nonBackupableJob *orchestratorFakes.FakeJob

		BeforeEach(func() {
			backupableJob = new(orchestratorFakes.FakeJob)
			backupableJob.HasBackupReturns(true)

			nonBackupableJob = new(orchestratorFakes.FakeJob)
			nonBackupableJob.HasBackupReturns(false)

			jobs = orchestrator.Jobs([]orchestrator.Job{
				backupableJob,
				nonBackupableJob,
			})
		})

		Describe("Backupable", func() {
			It("returns the backupable job", func() {
				Expect(jobs.Backupable()).To(ConsistOf(backupableJob))
			})
		})

		Describe("AnyAreBackupable", func() {
			It("returns true", func() {
				Expect(jobs.AnyAreBackupable()).To(BeTrue())
			})
		})
	})

	Context("contains no jobs with backup script", func() {
		var nonBackupableJob *orchestratorFakes.FakeJob

		BeforeEach(func() {
			nonBackupableJob = new(orchestratorFakes.FakeJob)
			nonBackupableJob.HasBackupReturns(false)

			jobs = orchestrator.Jobs([]orchestrator.Job{
				nonBackupableJob,
			})
		})

		Describe("Backupable", func() {
			It("returns empty", func() {
				Expect(jobs.Backupable()).To(BeEmpty())
			})
		})

		Describe("AnyAreBackupable", func() {
			It("returns false", func() {
				Expect(jobs.AnyAreBackupable()).To(BeFalse())
			})
		})
	})
})
