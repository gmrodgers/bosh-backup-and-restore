package orchestrator

import (
	"github.com/cloudfoundry-incubator/bosh-backup-and-restore/executor"
)

func NewLocker(logger Logger, deploymentManager DeploymentManager,
	lockOrderer LockOrderer, executor executor.Executor) *Locker {

	findDeploymentStep := NewFindDeploymentStep(deploymentManager, logger)
	backupable := NewBackupableStep(lockOrderer, logger)
	lock := NewLockStep(lockOrderer, executor)
	unlock := NewPostBackupUnlockStep(lockOrderer, executor)

	workflow := NewWorkflow()
	workflow.StartWith(findDeploymentStep).OnSuccess(backupable)
	workflow.Add(backupable).OnSuccess(lock)
	workflow.Add(lock).OnFailure(unlock)
	workflow.Add(unlock)

	return &Locker{
		workflow: workflow,
	}
}

type Locker struct {
	workflow *Workflow
}

type AuthInfo struct {
	Type   string
	UaaUrl string
}

//Backup checks if a deployment has backupable instances and backs them up.
func (b Locker) Lock(deploymentName string) Error {
	session := NewSession(deploymentName)

	err := b.workflow.Run(session)

	return err
}
