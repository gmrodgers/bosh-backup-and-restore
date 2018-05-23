package instance_test

import (
	"github.com/cloudfoundry-incubator/bosh-backup-and-restore/instance"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"

	"log"

	"fmt"

	sshfakes "github.com/cloudfoundry-incubator/bosh-backup-and-restore/ssh/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("Job", func() {
	var job instance.Job
	var jobScripts instance.BackupAndRestoreScripts
	var metadata instance.Metadata
	var stdout, stderr *gbytes.Buffer
	var logger boshlog.Logger
	var releaseName string
	var remoteRunner *sshfakes.FakeRemoteRunner
	var instanceIdentifier = "instance/identifier"

	BeforeEach(func() {
		jobScripts = instance.BackupAndRestoreScripts{
			"/var/vcap/jobs/jobname/bin/bbr/restore",
			"/var/vcap/jobs/jobname/bin/bbr/backup",
			"/var/vcap/jobs/jobname/bin/bbr/pre-backup-lock",
			"/var/vcap/jobs/jobname/bin/bbr/post-backup-unlock",
		}
		metadata = instance.Metadata{}
		stdout = gbytes.NewBuffer()
		stderr = gbytes.NewBuffer()
		stdoutLog := log.New(stdout, "[instance-test] ", log.Lshortfile)
		stderrLog := log.New(stderr, "[instance-test] ", log.Lshortfile)
		logger = boshlog.New(boshlog.LevelDebug, stdoutLog, stderrLog)
		releaseName = "redis"
		remoteRunner = new(sshfakes.FakeRemoteRunner)
	})

	JustBeforeEach(func() {
		job = instance.NewJob(remoteRunner, instanceIdentifier, logger, releaseName, jobScripts, metadata)
	})

	Describe("HasBackup", func() {
		It("returns true", func() {
			Expect(job.HasBackup()).To(BeTrue())
		})

		Context("no backup scripts exist", func() {
			BeforeEach(func() {
				jobScripts = instance.BackupAndRestoreScripts{"/var/vcap/jobs/jobname/bin/bbr/restore"}
			})

			It("returns false", func() {
				Expect(job.HasBackup()).To(BeFalse())
			})
		})
	})

	Describe("PreBackupLock", func() {
		var preBackupLockError error

		JustBeforeEach(func() {
			preBackupLockError = job.PreBackupLock()
		})

		Context("job has no pre-backup-lock script", func() {
			BeforeEach(func() {
				jobScripts = instance.BackupAndRestoreScripts{
					"/var/vcap/jobs/jobname/bin/bbr/restore",
				}
			})

			It("should not call the remote runner", func() {
				Expect(remoteRunner.Invocations()).To(HaveLen(0))
			})
		})

		Context("job has a pre-backup-lock script", func() {
			BeforeEach(func() {
				jobScripts = instance.BackupAndRestoreScripts{
					"/var/vcap/jobs/jobname/bin/bbr/pre-backup-lock",
				}
			})

			It("runs the script", func() {
				By("calling the remote runner", func() {
					Expect(remoteRunner.RunScriptCallCount()).To(Equal(1))
					cmd, _ := remoteRunner.RunScriptArgsForCall(0)
					Expect(cmd).To(Equal("/var/vcap/jobs/jobname/bin/bbr/pre-backup-lock"))
				})

				By("logging the script path", func() {
					Expect(string(stdout.Contents())).To(ContainSubstring(`> /var/vcap/jobs/jobname/bin/bbr/pre-backup-lock`))
				})

				By("logging the job name that it has locked", func() {
					Expect(string(stdout.Contents())).To(ContainSubstring(fmt.Sprintf(
						"INFO - Locking jobname on %s",
						instanceIdentifier,
					)))
					Expect(string(stdout.Contents())).To(ContainSubstring(fmt.Sprintf(
						"INFO - Finished locking jobname on %s",
						instanceIdentifier,
					)))
				})
			})

			Context("pre-backup-lock script runs successfully", func() {
				BeforeEach(func() {
					remoteRunner.RunScriptReturns("stdout", nil)
				})

				It("succeeds", func() {
					Expect(preBackupLockError).NotTo(HaveOccurred())
				})
			})

			Context("pre-backup-lock script errors", func() {
				BeforeEach(func() {
					remoteRunner.RunScriptReturns("", fmt.Errorf("some strange error"))
				})

				It("fails", func() {
					By("including the error in the returned error", func() {
						Expect(preBackupLockError).To(MatchError(ContainSubstring("some strange error")))
					})
				})
			})
		})
	})

	Describe("PostBackupUnlock", func() {
		var postBackupUnlockError error

		JustBeforeEach(func() {
			postBackupUnlockError = job.PostBackupUnlock()
		})

		Context("job has no post-backup-unlock script", func() {
			BeforeEach(func() {
				jobScripts = instance.BackupAndRestoreScripts{
					"/var/vcap/jobs/jobname/bin/bbr/restore",
				}
			})

			It("should not run anything on the remote runner", func() {
				Expect(remoteRunner.Invocations()).To(HaveLen(0))
			})
		})

		Context("job has a post-backup-unlock script", func() {
			BeforeEach(func() {
				jobScripts = instance.BackupAndRestoreScripts{
					"/var/vcap/jobs/jobname/bin/bbr/post-backup-unlock",
				}
			})

			It("uses remote runner to run the script", func() {
				Expect(remoteRunner.RunScriptCallCount()).To(Equal(1))
				cmd, _ := remoteRunner.RunScriptArgsForCall(0)
				Expect(cmd).To(Equal("/var/vcap/jobs/jobname/bin/bbr/post-backup-unlock"))
			})

			Context("post-backup-unlock script runs successfully", func() {
				BeforeEach(func() {
					remoteRunner.RunScriptReturns("stdout", nil)
				})

				It("succeeds", func() {
					Expect(postBackupUnlockError).NotTo(HaveOccurred())
				})
			})

			Context("post-backup-unlock script fails", func() {
				BeforeEach(func() {
					remoteRunner.RunScriptReturns("", fmt.Errorf("it failed"))
				})

				It("fails", func() {
					Expect(postBackupUnlockError).To(MatchError(ContainSubstring("it failed")))
				})
			})
		})
	})

	Describe("Release", func() {
		It("returns the job's release name", func() {
			Expect(job.Release()).To(Equal("redis"))
		})
	})
})
