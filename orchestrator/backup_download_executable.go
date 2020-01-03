package orchestrator

import (
	"fmt"
	"sync"
	"time"

	"code.cloudfoundry.org/bytefmt"
	"github.com/cloudfoundry-incubator/bosh-backup-and-restore/counter"
	"github.com/pkg/errors"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/number"
)

//go:generate counterfeiter -o fakes/fake_clock.go . Clock
type Clock interface {
	Sleep(time.Duration)
}

type BackupDownloadExecutable struct {
	localBackup    Backup
	remoteArtifact BackupArtifact
	Logger         Logger
	Clock          Clock
}

func NewBackupDownloadExecutable(localBackup Backup, remoteArtifact BackupArtifact, logger Logger) BackupDownloadExecutable {
	return BackupDownloadExecutable{
		localBackup:    localBackup,
		remoteArtifact: remoteArtifact,
		Logger:         logger,
		Clock:          clock{},
	}
}

func (e BackupDownloadExecutable) Execute() error {
	err := e.downloadBackupArtifact(e.localBackup, e.remoteArtifact)
	if err != nil {
		return err
	}

	checksum, err := e.compareChecksums(e.localBackup, e.remoteArtifact)
	if err != nil {
		return err
	}

	err = e.localBackup.AddChecksum(e.remoteArtifact, checksum)
	if err != nil {
		return err
	}

	err = e.remoteArtifact.Delete()
	if err != nil {
		return err
	}

	e.Logger.Info("bbr", "Finished validity checks -- for job %s on %s/%s...", e.remoteArtifact.Name(), e.remoteArtifact.InstanceName(), e.remoteArtifact.InstanceID())
	return nil
}

func (e BackupDownloadExecutable) downloadBackupArtifact(localBackup Backup, remoteBackupArtifact BackupArtifact) error {
	localBackupArtifactWriter, err := localBackup.CreateArtifact(remoteBackupArtifact)
	if err != nil {
		return err
	}

	size, err := remoteBackupArtifact.Size()
	if err != nil {
		return err
	}

	writerCounter := counter.NewCountWriter(localBackupArtifactWriter)
	e.logAmountTransfered(writerCounter, size, remoteBackupArtifact.Name(), remoteBackupArtifact.InstanceName(), remoteBackupArtifact.InstanceID())

	var wg sync.WaitGroup
	wg.Add(1)

	cancelLogging := make(chan struct{})
	go func() {
		e.Clock.Sleep(time.Second * 5)
		for {
			select {
			case <-cancelLogging:
				wg.Done()
				return
			default:
				e.logAmountTransfered(writerCounter, size, remoteBackupArtifact.Name(), remoteBackupArtifact.InstanceName(), remoteBackupArtifact.InstanceID())
				e.Clock.Sleep(time.Second * 5)
			}
		}
	}()

	err = remoteBackupArtifact.StreamFromRemote(writerCounter)
	close(cancelLogging)
	wg.Wait()
	if err != nil {
		return err
	}

	err = localBackupArtifactWriter.Close()
	if err != nil {
		return err
	}

	e.logAmountTransfered(writerCounter, size, remoteBackupArtifact.Name(), remoteBackupArtifact.InstanceName(), remoteBackupArtifact.InstanceID())
	return nil
}

func (e BackupDownloadExecutable) compareChecksums(localBackup Backup, remoteBackupArtifact BackupArtifact) (BackupChecksum, error) {
	e.Logger.Info("bbr", "Starting validity checks -- for job %s on %s/%s...", remoteBackupArtifact.Name(), remoteBackupArtifact.InstanceName(), remoteBackupArtifact.InstanceID())

	localChecksum, err := localBackup.CalculateChecksum(remoteBackupArtifact)
	if err != nil {
		return nil, err
	}

	remoteChecksum, err := remoteBackupArtifact.Checksum()
	if err != nil {
		return nil, err
	}

	e.Logger.Debug("bbr", "Comparing shasums")

	match, mismatchedFiles := localChecksum.Match(remoteChecksum)
	if !match {
		e.Logger.Debug("bbr", "Checksums didn't match for:")
		e.Logger.Debug("bbr", fmt.Sprintf("%v\n", mismatchedFiles))

		err = errors.Errorf(
			"Backup is corrupted, checksum failed for %s/%s %s - checksums don't match for %v. "+
				"Checksum failed for %d files in total",
			remoteBackupArtifact.InstanceName(), remoteBackupArtifact.InstanceID(), remoteBackupArtifact.Name(), getFirstTen(mismatchedFiles), len(mismatchedFiles))
		return nil, err
	}

	return localChecksum, nil
}

func (e BackupDownloadExecutable) logAmountTransfered(writerCounter *counter.CountWriter, size, backupName, instanceName, instanceID string) {
	byteSize, _ := bytefmt.ToBytes(size)
	decimalPercentageComplete := float64(writerCounter.Count()) / float64(byteSize)
	percentageComplete := number.Percent(decimalPercentageComplete, number.MaxFractionDigits(0))

	printer := message.NewPrinter(language.English)
	p := printer.Sprintf("%v", percentageComplete)

	e.Logger.Info("bbr", "Copying backup -- %s of %s complete -- for job %s on %s/%s...", p, size, backupName, instanceName, instanceID)
}
