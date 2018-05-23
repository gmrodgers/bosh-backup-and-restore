package bosh

import (
	"github.com/cloudfoundry-incubator/bosh-backup-and-restore/orchestrator"
	"github.com/pkg/errors"
)

func NewDeploymentManager(boshDirector BoshClient, logger Logger) *DeploymentManager {
	return &DeploymentManager{BoshClient: boshDirector, Logger: logger}
}

type DeploymentManager struct {
	BoshClient
	Logger
}

func (b *DeploymentManager) Find(deploymentName string) (orchestrator.Deployment, error) {
	instances, err := b.FindInstances(deploymentName)
	return orchestrator.NewDeployment(b.Logger, instances), errors.Wrap(err, "failed to find instances for deployment "+deploymentName)
}
