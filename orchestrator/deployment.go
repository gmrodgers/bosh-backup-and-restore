package orchestrator

import (
	"github.com/cloudfoundry-incubator/bosh-backup-and-restore/executor"
)

//go:generate counterfeiter -o fakes/fake_deployment.go . Deployment
type Deployment interface {
	IsBackupable() bool
	BackupableInstances() []Instance
	PreBackupLock(LockOrderer, executor.Executor) error
	PostBackupUnlock(LockOrderer, executor.Executor) error
	Instances() []Instance
	ValidateLockingDependencies(orderer LockOrderer) error
}

//go:generate counterfeiter -o fakes/fake_lock_orderer.go . LockOrderer
type LockOrderer interface {
	Order(jobs []Job) ([][]Job, error)
}

type deployment struct {
	Logger
	instances instances
}

func NewDeployment(logger Logger, instancesArray []Instance) Deployment {
	return &deployment{Logger: logger, instances: instances(instancesArray)}
}

func (bd *deployment) IsBackupable() bool {
	backupableInstances := bd.instances.AllBackupable()
	return !backupableInstances.IsEmpty()
}

func (bd *deployment) BackupableInstances() []Instance {
	return bd.instances.AllBackupable()
}

func (bd *deployment) ValidateLockingDependencies(lockOrderer LockOrderer) error {
	jobs := bd.instances.Jobs()
	_, err := lockOrderer.Order(jobs)
	return err
}

func (bd *deployment) PreBackupLock(lockOrderer LockOrderer, executor executor.Executor) error {
	bd.Logger.Info("db-lock", "Running pre-backup-lock scripts...")

	jobs := bd.instances.Jobs()

	orderedJobs, err := lockOrderer.Order(jobs)
	if err != nil {
		return err
	}

	preBackupLockErrors := executor.Run(newJobExecutables(orderedJobs, NewJobPreBackupLockExecutable))

	bd.Logger.Info("db-lock", "Finished running pre-backup-lock scripts.")
	return ConvertErrors(preBackupLockErrors)
}

func (bd *deployment) PostBackupUnlock(lockOrderer LockOrderer, executor executor.Executor) error {
	bd.Logger.Info("db-lock", "Running post-backup-unlock scripts...")

	jobs := bd.instances.Jobs()

	orderedJobs, err := lockOrderer.Order(jobs)
	if err != nil {
		return err
	}
	reversedJobs := Reverse(orderedJobs)

	postBackupUnlockErrors := executor.Run(newJobExecutables(reversedJobs, NewJobPostBackupUnlockExecutable))
	bd.Logger.Info("db-lock", "Finished running post-backup-unlock scripts.")
	return ConvertErrors(postBackupUnlockErrors)
}

func newJobExecutables(jobsList [][]Job, newJobExecutable func(Job) executor.Executable) [][]executor.Executable {
	var executablesList [][]executor.Executable
	for _, jobs := range jobsList {
		var executables []executor.Executable
		for _, job := range jobs {
			executables = append(executables, newJobExecutable(job))
		}
		executablesList = append(executablesList, executables)
	}
	return executablesList
}

func (bd *deployment) Instances() []Instance {
	return bd.instances
}

func getFirstTen(input []string) (output []string) {
	for i := 0; i < len(input); i++ {
		if i == 10 {
			break
		}
		output = append(output, input[i])
	}
	return output
}
