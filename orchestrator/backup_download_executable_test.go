package orchestrator_test

import (
	"fmt"
	"io"
	"time"

	"github.com/cloudfoundry-incubator/bosh-backup-and-restore/counter"
	"github.com/cloudfoundry-incubator/bosh-backup-and-restore/orchestrator"
	"github.com/cloudfoundry-incubator/bosh-backup-and-restore/orchestrator/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("BackupDownloadExecutable", func() {
	var (
		executable                orchestrator.BackupDownloadExecutable
		localBackup               *fakes.FakeBackup
		remoteArtifact            *fakes.FakeBackupArtifact
		logger                    *fakes.FakeLogger
		localBackupArtifactWriter *fakes.FakeWriteCloser
		fakeClock                 *fakes.FakeClock
		actualError               error
	)

	BeforeEach(func() {
		fakeClock = new(fakes.FakeClock)
		localBackup = new(fakes.FakeBackup)
		remoteArtifact = new(fakes.FakeBackupArtifact)
		remoteArtifact.SizeReturns("1K", nil)
		logger = new(fakes.FakeLogger)
		localBackupArtifactWriter = new(fakes.FakeWriteCloser)

		localBackup.CreateArtifactReturns(localBackupArtifactWriter, nil)
		remoteArtifact.StreamFromRemoteStub = func(w io.Writer) error {
			w.Write([]byte(""))
			for fakeClock.SleepCallCount() < 1 {
			}
			w.Write([]byte(""))
			return nil
		}

		localBackupArtifactWriter.WriteStub = func(p []byte) (int, error) {
			return 512, nil
		}
	})

	JustBeforeEach(func() {
		executable = orchestrator.NewBackupDownloadExecutable(localBackup, remoteArtifact, logger)
		executable.Clock = fakeClock
		actualError = executable.Execute()
	})

	It("downloads the artifact", func() {
		By("not failing", func() {
			Expect(actualError).NotTo(HaveOccurred())
		})

		By("creating a local artifact", func() {
			Expect(localBackup.CreateArtifactCallCount()).To(Equal(1))
			Expect(localBackup.CreateArtifactArgsForCall(0)).To(Equal(remoteArtifact))
		})

		By("calculating the size of the remote artifact", func() {
			Expect(remoteArtifact.SizeCallCount()).To(Equal(1))
		})

		By("streaming from the remote artifact", func() {
			Expect(remoteArtifact.StreamFromRemoteCallCount()).To(Equal(1))
			Expect(remoteArtifact.StreamFromRemoteArgsForCall(0)).To(BeAssignableToTypeOf(&counter.CountWriter{}))
		})

		By("closing the local backup artifact writer", func() {
			Expect(localBackupArtifactWriter.CloseCallCount()).To(Equal(1))
		})

		By("calculating the local checksum", func() {
			Expect(localBackup.CalculateChecksumCallCount()).To(Equal(1))
		})

		By("calculating the remote checksum", func() {
			Expect(remoteArtifact.ChecksumCallCount()).To(Equal(1))
		})

		By("logging the download", func() {
			By("waiting 5s")
			Expect(fakeClock.SleepCallCount()).To(BeNumerically(">", 0))
			Expect(fakeClock.SleepArgsForCall(0)).To(Equal(time.Second * 5))

			Expect(logger.InfoCallCount()).To(BeNumerically(">", 0))
			By("logging 0% when starting")
			prefix, message, args := logger.InfoArgsForCall(0)
			Expect(prefix).To(Equal("bbr"))
			Expect(message).To(ContainSubstring("Copying backup -- %s of %s complete -- for job %s on %s"))
			Expect(args[0]).To(Equal("0%"))

			By("logging 50% when halfway complete")
			prefix, message, args = logger.InfoArgsForCall(1)
			Expect(prefix).To(Equal("bbr"))
			Expect(message).To(ContainSubstring("Copying backup -- %s of %s complete -- for job %s on %s"))
			Expect(args[0]).To(Equal("50%"))

			By("waiting 5s")
			Expect(fakeClock.SleepCallCount()).To(Equal(2))
			Expect(fakeClock.SleepArgsForCall(1)).To(Equal(time.Second * 5))

			By("logging 100% when complete")
			prefix, message, args = logger.InfoArgsForCall(2)
			Expect(prefix).To(Equal("bbr"))
			Expect(message).To(ContainSubstring("Copying backup -- %s of %s complete -- for job %s on %s"))
			Expect(args[0]).To(Equal("100%"))
		})
	})

	Context("When the local artifact cannot be created", func() {
		BeforeEach(func() {
			localBackup.CreateArtifactReturns(nil, fmt.Errorf("create artifact error"))
		})

		It("should fail", func() {
			Expect(actualError).To(MatchError("create artifact error"))
		})
	})

	Context("When the remote artifact size cannot be determined", func() {
		BeforeEach(func() {
			remoteArtifact.SizeReturns("", fmt.Errorf("size error"))
		})

		It("should fail", func() {
			Expect(actualError).To(MatchError("size error"))
		})
	})

	Context("When streaming the remote artifact fails", func() {
		BeforeEach(func() {
			remoteArtifact.StreamFromRemoteReturns(fmt.Errorf("stream error"))
		})

		It("should fail", func() {
			Expect(actualError).To(MatchError("stream error"))
		})
	})

	Context("When the local backup writer cannot be closed", func() {
		BeforeEach(func() {
			localBackupArtifactWriter.CloseReturns(fmt.Errorf("close error"))
		})

		It("should fail", func() {
			Expect(actualError).To(MatchError("close error"))
		})
	})

	Context("When the local backup checksum cannot be calculated", func() {
		BeforeEach(func() {
			localBackup.CalculateChecksumReturns(nil, fmt.Errorf("local checksum error"))
		})

		It("should fail", func() {
			Expect(actualError).To(MatchError("local checksum error"))
		})
	})

	Context("When the remote artifact checksum cannot be calculated", func() {
		BeforeEach(func() {
			localBackup.CreateArtifactReturns(nil, fmt.Errorf("remote checksum error"))
		})

		It("should fail", func() {
			Expect(actualError).To(MatchError("remote checksum error"))
		})
	})

	Context("When the checksums are mismatched", func() {
		BeforeEach(func() {
			localBackup.CalculateChecksumReturns(orchestrator.BackupChecksum{"file1": "abcd", "file2": "not matching"}, nil)
			remoteArtifact.ChecksumReturns(orchestrator.BackupChecksum{"file1": "abcd", "file2": "efgh"}, nil)
		})

		It("should fail", func() {
			Expect(actualError).To(MatchError(ContainSubstring("Backup is corrupted, checksum failed")))
		})
	})

	Context("When the checksum cannot be added to the local backup", func() {
		BeforeEach(func() {
			localBackup.AddChecksumReturns(fmt.Errorf("add checksum error"))
		})

		It("should fail", func() {
			Expect(actualError).To(MatchError("add checksum error"))
		})
	})

	Context("When the remote artifact cannot be deleted", func() {
		BeforeEach(func() {
			remoteArtifact.DeleteReturns(fmt.Errorf("remote artifact deletion error"))
		})

		It("should fail", func() {
			Expect(actualError).To(MatchError("remote artifact deletion error"))
		})
	})
})
