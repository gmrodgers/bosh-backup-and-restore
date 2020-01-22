package instance

import (
	"fmt"
	"strconv"

	"github.com/cloudfoundry-incubator/bosh-backup-and-restore/orchestrator"
	"github.com/cloudfoundry-incubator/bosh-backup-and-restore/ssh"
	"github.com/pkg/errors"
)

type JobContext struct {
	Logger             Logger
	RemoteRunner       ssh.RemoteRunner
	InstanceIdentifier string
	Release            string
	Metadata           Metadata
	OnBootstrapNode    bool
}

func NewJob(jobScripts BackupAndRestoreScripts, backupOneRestoreAll bool, ctx JobContext) Job {
	jobName := jobScripts[0].JobName()
	return Job{
		name:                jobName,
		backupScript:        jobScripts.BackupOnly().firstOrBlank(),
		restoreScript:       jobScripts.RestoreOnly().firstOrBlank(),
		preBackupScript:     jobScripts.PreBackupLockOnly().firstOrBlank(),
		preRestoreScript:    jobScripts.PreRestoreLockOnly().firstOrBlank(),
		postBackupScript:    jobScripts.PostBackupUnlockOnly().firstOrBlank(),
		postRestoreScript:   jobScripts.SinglePostRestoreUnlockScript(),
		backupOneRestoreAll: backupOneRestoreAll,
		ctx:                 ctx,
	}
}

type Job struct {
	name                string
	backupScript        Script
	preBackupScript     Script
	postBackupScript    Script
	preRestoreScript    Script
	restoreScript       Script
	postRestoreScript   Script
	backupOneRestoreAll bool
	ctx                 JobContext
}

func (j Job) Name() string {
	return j.name
}

func (j Job) Release() string {
	return j.ctx.Release
}

func (j Job) InstanceIdentifier() string {
	return j.ctx.InstanceIdentifier
}

func (j Job) BackupArtifactName() string {
	if j.backupOneRestoreAll && j.ctx.OnBootstrapNode {
		return j.backupOneRestoreAllArtifactName()
	}

	return j.ctx.Metadata.BackupName
}

func (j Job) backupOneRestoreAllArtifactName() string {
	return fmt.Sprintf("%s-%s-backup-one-restore-all", j.name, j.ctx.Release)
}

func (j Job) HasMetadataRestoreName() bool {
	if j.ctx.Metadata.RestoreName != "" {
		return true
	}
	return false
}

func (j Job) RestoreArtifactName() string {
	if j.backupOneRestoreAll {
		return j.backupOneRestoreAllArtifactName()
	}

	return j.ctx.Metadata.RestoreName
}

func (j Job) BackupArtifactDirectory() string {
	return fmt.Sprintf("%s/%s", orchestrator.ArtifactDirectory, j.backupArtifactOrJobName())
}

func (j Job) RestoreArtifactDirectory() string {
	return fmt.Sprintf("%s/%s", orchestrator.ArtifactDirectory, j.restoreArtifactOrJobName())
}

func (j Job) RestoreScript() Script {
	return j.restoreScript
}

func (j Job) HasBackup() bool {
	return j.backupScript != ""
}

func (j Job) HasRestore() bool {
	return j.RestoreScript() != ""
}

func (j Job) HasNamedBackupArtifact() bool {
	return j.backupOneRestoreAll && j.ctx.OnBootstrapNode
}

func (j Job) HasNamedRestoreArtifact() bool {
	return j.backupOneRestoreAll || j.ctx.Metadata.RestoreName != ""
}

func (j Job) Backup() error {
	if j.backupScript != "" {
		j.ctx.Logger.Debug("bbr", "> %s", j.backupScript)
		j.ctx.Logger.Info("bbr", "Backing up %s on %s...", j.name, j.ctx.InstanceIdentifier)

		err := j.ctx.RemoteRunner.CreateDirectory(j.BackupArtifactDirectory())
		if err != nil {
			return err
		}

		env := artifactDirectoryVariables(j.BackupArtifactDirectory())
		_, err = j.ctx.RemoteRunner.RunScriptWithEnv(
			string(j.backupScript),
			env,
			fmt.Sprintf("backup %s on %s", j.name, j.ctx.InstanceIdentifier),
		)

		if err != nil {
			j.ctx.Logger.Error("bbr", "Error backing up %s on %s.", j.name, j.ctx.InstanceIdentifier)

			return errors.Wrap(err, fmt.Sprintf(
				"Error attempting to run backup for job %s on %s",
				j.Name(),
				j.ctx.InstanceIdentifier,
			))
		}

		j.ctx.Logger.Info("bbr", "Finished backing up %s on %s.", j.name, j.ctx.InstanceIdentifier)
	}

	return nil
}

func (j Job) PreBackupLock() error {
	if j.preBackupScript != "" {
		j.ctx.Logger.Debug("bbr", "> %s", j.preBackupScript)
		j.ctx.Logger.Info("bbr", "Locking %s on %s for backup...", j.name, j.ctx.InstanceIdentifier)

		_, err := j.ctx.RemoteRunner.RunScript(
			string(j.preBackupScript),
			fmt.Sprintf("pre-backup lock %s on %s", j.name, j.ctx.InstanceIdentifier),
		)
		if err != nil {
			j.ctx.Logger.Error("bbr", "Error locking %s on %s.", j.name, j.ctx.InstanceIdentifier)

			return errors.Wrap(err, fmt.Sprintf(
				"Error attempting to run pre-backup-lock for job %s on %s",
				j.Name(),
				j.ctx.InstanceIdentifier,
			))
		}

		j.ctx.Logger.Info("bbr", "Finished locking %s on %s for backup.", j.name, j.ctx.InstanceIdentifier)
	}

	return nil
}

func (j Job) PostBackupUnlock(afterSuccessfulBackup bool) error {
	if j.postBackupScript != "" {
		j.ctx.Logger.Debug("bbr", "> %s", j.postBackupScript)
		j.ctx.Logger.Info("bbr", "Unlocking %s on %s...", j.name, j.ctx.InstanceIdentifier)
		env := map[string]string{
			"BBR_AFTER_BACKUP_SCRIPTS_SUCCESSFUL": strconv.FormatBool(afterSuccessfulBackup),
		}
		_, err := j.ctx.RemoteRunner.RunScriptWithEnv(
			string(j.postBackupScript),
			env,
			fmt.Sprintf("post-backup unlock %s on %s", j.name, j.ctx.InstanceIdentifier),
		)
		if err != nil {
			j.ctx.Logger.Error("bbr", "Error unlocking %s on %s.", j.name, j.ctx.InstanceIdentifier)

			return errors.Wrap(err, fmt.Sprintf(
				"Error attempting to run post-backup-unlock for job %s on %s",
				j.Name(),
				j.ctx.InstanceIdentifier,
			))
		}

		j.ctx.Logger.Info("bbr", "Finished unlocking %s on %s.", j.name, j.ctx.InstanceIdentifier)
	}

	return nil
}

func (j Job) PreRestoreLock() error {
	if j.preRestoreScript != "" {
		j.ctx.Logger.Debug("bbr", "> %s", j.preRestoreScript)
		j.ctx.Logger.Info("bbr", "Locking %s on %s for restore...", j.name, j.ctx.InstanceIdentifier)

		_, err := j.ctx.RemoteRunner.RunScript(
			string(j.preRestoreScript),
			fmt.Sprintf("pre-restore lock %s on %s", j.name, j.ctx.InstanceIdentifier),
		)
		if err != nil {
			j.ctx.Logger.Error("bbr", "Error locking %s on %s.", j.name, j.ctx.InstanceIdentifier)

			return errors.Wrap(err, fmt.Sprintf(
				"Error attempting to run pre-restore-lock for job %s on %s",
				j.Name(),
				j.ctx.InstanceIdentifier,
			))
		}

		j.ctx.Logger.Info("bbr", "Finished locking %s on %s for restore.", j.name, j.ctx.InstanceIdentifier)
	}

	return nil
}

func (j Job) Restore() error {
	if j.restoreScript != "" {
		j.ctx.Logger.Debug("bbr", "> %s", j.restoreScript)
		j.ctx.Logger.Info("bbr", "Restoring %s on %s...", j.name, j.ctx.InstanceIdentifier)

		env := artifactDirectoryVariables(j.RestoreArtifactDirectory())
		_, err := j.ctx.RemoteRunner.RunScriptWithEnv(
			string(j.restoreScript), env,
			fmt.Sprintf("restore %s on %s", j.name, j.ctx.InstanceIdentifier),
		)
		if err != nil {
			j.ctx.Logger.Error("bbr", "Error restoring %s on %s.", j.name, j.ctx.InstanceIdentifier)
			return errors.Wrap(err, fmt.Sprintf(
				"Error attempting to run restore for job %s on %s",
				j.Name(),
				j.ctx.InstanceIdentifier,
			))
		}

		j.ctx.Logger.Info("bbr", "Finished restoring %s on %s.", j.name, j.ctx.InstanceIdentifier)
	}

	return nil
}

func (j Job) PostRestoreUnlock() error {
	if j.postRestoreScript != "" {
		j.ctx.Logger.Debug("bbr", "> %s", j.postRestoreScript)
		j.ctx.Logger.Info("bbr", "Unlocking %s on %s...", j.name, j.ctx.InstanceIdentifier)

		_, err := j.ctx.RemoteRunner.RunScript(
			string(j.postRestoreScript),
			fmt.Sprintf("post-restore unlock %s on %s", j.name, j.ctx.InstanceIdentifier),
		)
		if err != nil {
			j.ctx.Logger.Error("bbr", "Error unlocking %s on %s.", j.name, j.ctx.InstanceIdentifier)

			return errors.Wrap(err, fmt.Sprintf(
				"Error attempting to run post-restore-unlock for job %s on %s",
				j.Name(),
				j.ctx.InstanceIdentifier,
			))
		}

		j.ctx.Logger.Info("bbr", "Finished unlocking %s on %s.", j.name, j.ctx.InstanceIdentifier)
	}

	return nil
}

func (j Job) backupArtifactOrJobName() string {
	if j.HasNamedBackupArtifact() {
		return j.BackupArtifactName()
	}

	return j.name
}

func (j Job) restoreArtifactOrJobName() string {
	if j.HasNamedRestoreArtifact() {
		return j.RestoreArtifactName()
	}

	return j.name
}

func (j Job) handleErrs(jobName, label string, err error, exitCode int, stdout, stderr []byte) error {
	var foundErrors []error

	if err != nil {
		j.ctx.Logger.Error("bbr", fmt.Sprintf(
			"Error attempting to run %s script for job %s on %s. Error: %s",
			label,
			jobName,
			j.ctx.InstanceIdentifier,
			err.Error(),
		))
		foundErrors = append(foundErrors, err)
	} else if exitCode != 0 {
		errorString := fmt.Sprintf(
			"%s script for job %s failed on %s.\nStdout: %s\nStderr: %s",
			label,
			jobName,
			j.ctx.InstanceIdentifier,
			stdout,
			stderr,
		)

		foundErrors = append(foundErrors, errors.New(errorString))

		j.ctx.Logger.Error("bbr", errorString)
	}

	return orchestrator.ConvertErrors(foundErrors)
}

func (j Job) BackupShouldBeLockedBefore() []orchestrator.JobSpecifier {
	jobSpecifiers := []orchestrator.JobSpecifier{}

	for _, lockBefore := range j.ctx.Metadata.BackupShouldBeLockedBefore {
		jobSpecifiers = append(jobSpecifiers, orchestrator.JobSpecifier{
			Name: lockBefore.JobName, Release: lockBefore.Release,
		})
	}

	return jobSpecifiers
}

func (j Job) RestoreShouldBeLockedBefore() []orchestrator.JobSpecifier {
	jobSpecifiers := []orchestrator.JobSpecifier{}

	for _, lockBefore := range j.ctx.Metadata.RestoreShouldBeLockedBefore {
		jobSpecifiers = append(jobSpecifiers, orchestrator.JobSpecifier{
			Name: lockBefore.JobName, Release: lockBefore.Release,
		})
	}

	return jobSpecifiers
}
