package orchestrator

import "github.com/cloudfoundry-incubator/bosh-backup-and-restore/executor"

type JobPreBackupLockExecutor struct {
	Job
}

func NewJobPreBackupLockExecutable(job Job) executor.Executable {
	return JobPreBackupLockExecutor{job}
}

func (j JobPreBackupLockExecutor) Execute() error {
	return j.PreBackupLock()
}

type JobPostBackupUnlockExecutor struct {
	Job
}

func NewJobPostBackupUnlockExecutable(job Job) executor.Executable {
	return JobPostBackupUnlockExecutor{job}
}

func (j JobPostBackupUnlockExecutor) Execute() error {
	return j.PostBackupUnlock()
}
