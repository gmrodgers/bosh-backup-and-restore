package instance

import (
	"fmt"

	"github.com/cloudfoundry-incubator/bosh-backup-and-restore/orchestrator"
	"github.com/cloudfoundry-incubator/bosh-backup-and-restore/ssh"
	"github.com/pkg/errors"
)

func NewJob(remoteRunner ssh.RemoteRunner, instanceIdentifier string, logger Logger, release string,
	jobScripts BackupAndRestoreScripts, metadata Metadata) Job {
	jobName := jobScripts[0].JobName()
	return Job{
		Logger:             logger,
		remoteRunner:       remoteRunner,
		instanceIdentifier: instanceIdentifier,
		name:               jobName,
		release:            release,
		metadata:           metadata,
		backupScript:       jobScripts.BackupOnly().firstOrBlank(),
		preBackupScript:    jobScripts.PreBackupLockOnly().firstOrBlank(),
		postBackupScript:   jobScripts.PostBackupUnlockOnly().firstOrBlank(),
	}
}

type Job struct {
	Logger             Logger
	name               string
	release            string
	metadata           Metadata
	backupScript       Script
	preBackupScript    Script
	postBackupScript   Script
	remoteRunner       ssh.RemoteRunner
	instanceIdentifier string
}

func (j Job) Name() string {
	return j.name
}

func (j Job) Release() string {
	return j.release
}

func (j Job) InstanceIdentifier() string {
	return j.instanceIdentifier
}

func (j Job) HasBackup() bool {
	return j.backupScript != ""
}

func (j Job) PreBackupLock() error {
	if j.preBackupScript != "" {
		j.Logger.Debug("db-lock", "> %s", j.preBackupScript)
		j.Logger.Info("db-lock", "Locking %s on %s for backup...", j.name, j.instanceIdentifier)

		_, err := j.remoteRunner.RunScript(
			string(j.preBackupScript),
			fmt.Sprintf("pre-backup lock %s on %s", j.name, j.instanceIdentifier),
		)
		if err != nil {
			j.Logger.Error("db-lock", "Error locking %s on %s.", j.name, j.instanceIdentifier)

			return errors.Wrap(err, fmt.Sprintf(
				"Error attempting to run pre-backup-lock for job %s on %s",
				j.Name(),
				j.instanceIdentifier,
			))
		}

		j.Logger.Info("db-lock", "Finished locking %s on %s for backup.", j.name, j.instanceIdentifier)
	}

	return nil
}

func (j Job) PostBackupUnlock() error {
	if j.postBackupScript != "" {
		j.Logger.Debug("db-lock", "> %s", j.postBackupScript)
		j.Logger.Info("db-lock", "Unlocking %s on %s...", j.name, j.instanceIdentifier)

		_, err := j.remoteRunner.RunScript(
			string(j.postBackupScript),
			fmt.Sprintf("post-backup unlock %s on %s", j.name, j.instanceIdentifier),
		)
		if err != nil {
			j.Logger.Error("db-lock", "Error unlocking %s on %s.", j.name, j.instanceIdentifier)

			return errors.Wrap(err, fmt.Sprintf(
				"Error attempting to run post-backup-unlock for job %s on %s",
				j.Name(),
				j.instanceIdentifier,
			))
		}

		j.Logger.Info("db-lock", "Finished unlocking %s on %s.", j.name, j.instanceIdentifier)
	}

	return nil
}

func (j Job) handleErrs(jobName, label string, err error, exitCode int, stdout, stderr []byte) error {
	var foundErrors []error

	if err != nil {
		j.Logger.Error("db-lock", fmt.Sprintf(
			"Error attempting to run %s script for job %s on %s. Error: %s",
			label,
			jobName,
			j.instanceIdentifier,
			err.Error(),
		))
		foundErrors = append(foundErrors, err)
	} else if exitCode != 0 {
		errorString := fmt.Sprintf(
			"%s script for job %s failed on %s.\nStdout: %s\nStderr: %s",
			label,
			jobName,
			j.instanceIdentifier,
			stdout,
			stderr,
		)

		foundErrors = append(foundErrors, errors.New(errorString))

		j.Logger.Error("db-lock", errorString)
	}

	return orchestrator.ConvertErrors(foundErrors)
}

func (j Job) BackupShouldBeLockedBefore() []orchestrator.JobSpecifier {
	jobSpecifiers := []orchestrator.JobSpecifier{}

	for _, lockBefore := range j.metadata.BackupShouldBeLockedBefore {
		jobSpecifiers = append(jobSpecifiers, orchestrator.JobSpecifier{
			Name: lockBefore.JobName, Release: lockBefore.Release,
		})
	}

	return jobSpecifiers
}
