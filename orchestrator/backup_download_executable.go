package orchestrator

import (
	"context"
	"fmt"
	"time"

	"code.cloudfoundry.org/bytefmt"

	"github.com/machinebox/progress"
	"github.com/pkg/errors"
)

type BackupDownloadExecutable struct {
	localBackup    Backup
	remoteArtifact BackupArtifact
	Logger
}

func NewBackupDownloadExecutable(localBackup Backup, remoteArtifact BackupArtifact, logger Logger) BackupDownloadExecutable {
	return BackupDownloadExecutable{
		localBackup:    localBackup,
		remoteArtifact: remoteArtifact,
		Logger:         logger,
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

	progressWriter := progress.NewWriter(localBackupArtifactWriter)

	size, err := remoteBackupArtifact.Size()
	if err != nil {
		return err
	}

	length, err := bytefmt.ToBytes(size)
	if err != nil {
		return err
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		progressChan := progress.NewTicker(ctx, progressWriter, int64(length), 30*time.Second)

		for p := range progressChan {
			seconds := p.Remaining().Round(time.Second)
			if seconds < 1 {
				break
			}

			e.Logger.Info("bbr", "Copying backup -- %v remaining -- for job %s on %s/%s...", seconds, remoteBackupArtifact.Name(), remoteBackupArtifact.InstanceName(), remoteBackupArtifact.InstanceID())
		}
	}()

	e.Logger.Info("bbr", "Copying backup -- %s uncompressed -- for job %s on %s/%s...", size, remoteBackupArtifact.Name(), remoteBackupArtifact.InstanceName(), remoteBackupArtifact.InstanceID())
	err = remoteBackupArtifact.StreamFromRemote(progressWriter)
	cancel()
	if err != nil {
		return err
	}

	err = localBackupArtifactWriter.Close()
	if err != nil {
		return err
	}

	e.Logger.Info("bbr", "Finished copying backup -- for job %s on %s/%s...", remoteBackupArtifact.Name(), remoteBackupArtifact.InstanceName(), remoteBackupArtifact.InstanceID())
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
