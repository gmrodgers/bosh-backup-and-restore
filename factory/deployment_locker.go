package factory

import (
	"github.com/cloudfoundry-incubator/bosh-backup-and-restore/executor"
	"github.com/cloudfoundry-incubator/bosh-backup-and-restore/orchestrator"
	"github.com/cloudfoundry-incubator/bosh-backup-and-restore/orderer"
)

func BuildDeploymentLocker(target, username, password, caCert string, hasDebug bool) (*orchestrator.Locker, error) {
	logger := BuildLogger(hasDebug)
	deploymentManager, err := BuildBoshDeploymentManager(
		target,
		username,
		password,
		caCert,
		logger,
	)
	if err != nil {
		return nil, err
	}
	execr := executor.NewParallelExecutor()

	return orchestrator.NewLocker(logger, deploymentManager, orderer.NewKahnBackupLockOrderer(), execr), nil
}
