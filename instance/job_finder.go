package instance

import (
	"fmt"
	"path/filepath"

	"github.com/cloudfoundry-incubator/bosh-backup-and-restore/orchestrator"
	"github.com/cloudfoundry-incubator/bosh-backup-and-restore/ssh"
	"github.com/pkg/errors"
)

type InstanceIdentifier struct {
	InstanceGroupName string
	InstanceId        string
	Bootstrap         bool
}

func (i InstanceIdentifier) String() string {
	return fmt.Sprintf("%s/%s", i.InstanceGroupName, i.InstanceId)
}

//go:generate counterfeiter -o fakes/fake_job_finder.go . JobFinder
type JobFinder interface {
	FindJobs(instanceIdentifier InstanceIdentifier, remoteRunner ssh.RemoteRunner, manifestQuerier ManifestQuerier) (orchestrator.Jobs, error)
}

type JobFinderFromScripts struct {
	bbrVersion       string
	Logger           Logger
	parseJobMetadata MetadataParserFunc
}

func NewJobFinder(bbrVersion string, logger Logger) *JobFinderFromScripts {
	return &JobFinderFromScripts{
		bbrVersion:       bbrVersion,
		Logger:           logger,
		parseJobMetadata: ParseJobMetadata,
	}
}

func NewJobFinderOmitMetadataReleases(bbrVersion string, logger Logger) *JobFinderFromScripts {
	return &JobFinderFromScripts{
		bbrVersion:       bbrVersion,
		Logger:           logger,
		parseJobMetadata: ParseJobMetadataOmitReleases,
	}
}

func (jf *JobFinderFromScripts) FindJobs(
	instanceIdentifier InstanceIdentifier,
	remoteRunner ssh.RemoteRunner,
	manifestQuerier ManifestQuerier,
) (orchestrator.Jobs, error) {
	scripts, err := jf.allScripts(instanceIdentifier, remoteRunner)
	if err != nil {
		return nil, err
	}
	scriptsByJob := groupScriptsByJobName(scripts)

	metadataScripts := onlyMetadataScripts(scripts)
	metadataByJob, err := jf.groupMetadataByJobName(metadataScripts, instanceIdentifier, remoteRunner)
	if err != nil {
		return nil, err
	}

	return jf.buildJobs(remoteRunner, instanceIdentifier, scriptsByJob, metadataByJob, manifestQuerier)
}

func onlyMetadataScripts(scripts BackupAndRestoreScripts) BackupAndRestoreScripts {
	var metadataScripts BackupAndRestoreScripts
	for _, script := range scripts {
		if script.isMetadata() {
			metadataScripts = append(metadataScripts, script)
		}
	}
	return metadataScripts
}

func (jf *JobFinderFromScripts) groupMetadataByJobName(
	scripts BackupAndRestoreScripts,
	instanceIdentifier InstanceIdentifier,
	remoteRunner ssh.RemoteRunner,
) (map[string]Metadata, error) {
	metadata := map[string]Metadata{}
	for _, script := range scripts {
		jobMetadata, err := jf.extractMetadataFromScript(instanceIdentifier, script, remoteRunner)
		if err != nil {
			return nil, err
		}

		jf.logMetadata(jobMetadata, script.JobName())

		jobMetadata.BackupName = ""
		metadata[script.JobName()] = *jobMetadata
	}
	return metadata, nil
}

func (jf *JobFinderFromScripts) logMetadata(jobMetadata *Metadata, jobName string) {
	for _, lockBefore := range jobMetadata.BackupShouldBeLockedBefore {
		jf.Logger.Info("bbr", "Detected order: %s should be locked before %s during backup", jobName, filepath.Join(lockBefore.Release, lockBefore.JobName))
	}
	for _, lockBefore := range jobMetadata.RestoreShouldBeLockedBefore {
		jf.Logger.Info("bbr", "Detected order: %s should be locked before %s during restore", jobName, filepath.Join(lockBefore.Release, lockBefore.JobName))
	}

	if jobMetadata.BackupName != "" {
		jf.Logger.Warn("bbr", "discontinued metadata keys backup_name/restore_name found in job %s. bbr will not be able to restore this backup artifact.", jobName)
	}
}

func (jf *JobFinderFromScripts) extractMetadataFromScript(
	instanceIdentifier InstanceIdentifier,
	script Script,
	remoteRunner ssh.RemoteRunner,
) (*Metadata, error) {
	metadataContent, err := remoteRunner.RunScriptWithEnv(
		string(script),
		map[string]string{"BBR_VERSION": jf.bbrVersion},
		fmt.Sprintf("find metadata for %s on %s", script.JobName(), instanceIdentifier),
	)
	if err != nil {
		return nil, errors.Wrapf(
			err,
			"An error occurred while running metadata script for job %s on %s",
			script.JobName(),
			instanceIdentifier,
		)
	}
	jobMetadata, err := jf.parseJobMetadata(metadataContent)
	if err != nil {
		return nil, errors.Wrapf(
			err,
			"Parsing metadata from job %s on %s failed",
			script.JobName(),
			instanceIdentifier,
		)
	}

	return jobMetadata, nil
}

func (jf *JobFinderFromScripts) allScripts(
	instanceIdentifierForLogging InstanceIdentifier,
	remoteRunner ssh.RemoteRunner,
) (BackupAndRestoreScripts, error) {
	jf.Logger.Debug("bbr", "Attempting to find scripts on %s", instanceIdentifierForLogging)

	scripts, err := remoteRunner.FindFiles("/var/vcap/jobs/*/bin/bbr/*")
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("finding scripts failed on %s", instanceIdentifierForLogging))
	}

	return NewBackupAndRestoreScripts(scripts), nil
}

func groupScriptsByJobName(scripts BackupAndRestoreScripts) map[string]BackupAndRestoreScripts {
	scriptsByJob := map[string]BackupAndRestoreScripts{}
	for _, script := range scripts {
		jobName := script.JobName()
		scriptsByJob[jobName] = append(scriptsByJob[jobName], script)
	}
	return scriptsByJob
}

func (jf *JobFinderFromScripts) buildJobs(
	remoteRunner ssh.RemoteRunner,
	instanceIdentifier InstanceIdentifier,
	scriptsByJob map[string]BackupAndRestoreScripts,
	metadataByJob map[string]Metadata,
	manifestQuerier ManifestQuerier,
) (orchestrator.Jobs, error) {
	var jobs orchestrator.Jobs

	for jobName, jobScripts := range scriptsByJob {
		if metadataByJob[jobName].SkipBBRScripts {
			continue
		}

		for _, jobScript := range jobScripts {
			jf.Logger.Info("bbr", "%s/%s/%s", instanceIdentifier, jobName, jobScript.Name())
		}

		releaseName, err := manifestQuerier.FindReleaseName(instanceIdentifier.InstanceGroupName, jobName)
		if err != nil {
			jf.Logger.Warn("bbr", "could not find release name for job %s", jobName)
			releaseName = ""
		}

		backupOneRestoreAll, _ := manifestQuerier.IsJobBackupOneRestoreAll(instanceIdentifier.InstanceGroupName, jobName)

		jobs = append(jobs, NewJob(
			remoteRunner,
			instanceIdentifier.String(),
			jf.Logger,
			releaseName,
			jobScripts,
			metadataByJob[jobName],
			backupOneRestoreAll,
			instanceIdentifier.Bootstrap,
		))
	}

	if len(getSkippedJobs(scriptsByJob, metadataByJob)) != 0 {
		jf.logSkippedJobs(getSkippedJobs(scriptsByJob, metadataByJob), instanceIdentifier)
	}
	return jobs, nil
}

func getSkippedJobs(scriptsByJob map[string]BackupAndRestoreScripts, metadata map[string]Metadata) []string {
	var skippedJobs []string
	for jobName := range scriptsByJob {
		if metadata[jobName].SkipBBRScripts {
			skippedJobs = append(skippedJobs, jobName)
		}
	}
	return skippedJobs
}

func (jf *JobFinderFromScripts) logSkippedJobs(skippedJobs []string, instanceIdentifier InstanceIdentifier) {
	var skippedJobsMsg = "Found disabled jobs on instance"
	skippedJobsMsg = fmt.Sprintf("%s %s jobs:", skippedJobsMsg, instanceIdentifier)
	for _, job := range skippedJobs {
		skippedJobsMsg = skippedJobsMsg + " " + job
	}

	jf.Logger.Debug("bbr", skippedJobsMsg)
}
