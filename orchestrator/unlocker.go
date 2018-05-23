package orchestrator

import (
	"github.com/cloudfoundry-incubator/bosh-backup-and-restore/executor"
)

func NewUnlocker(logger Logger, deploymentManager DeploymentManager,
	lockOrderer LockOrderer, executor executor.Executor) *Unlocker {

	findDeploymentStep := NewFindDeploymentStep(deploymentManager, logger)
	backupable := NewBackupableStep(lockOrderer, logger)
	unlock := NewPostBackupUnlockStep(lockOrderer, executor)

	workflow := NewWorkflow()
	workflow.StartWith(findDeploymentStep).OnSuccess(backupable)
	workflow.Add(backupable).OnSuccess(unlock)
	workflow.Add(unlock)

	return &Unlocker{
		workflow: workflow,
	}
}

type Unlocker struct {
	workflow *Workflow
}

//Backup checks if a deployment has backupable instances and backs them up.
func (b Unlocker) Unlock(deploymentName string) Error {
	session := NewSession(deploymentName)

	err := b.workflow.Run(session)

	return err
}
