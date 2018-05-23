package standalone

import (
	"github.com/cloudfoundry-incubator/bosh-backup-and-restore/instance"
	"github.com/cloudfoundry-incubator/bosh-backup-and-restore/orchestrator"
	"github.com/cloudfoundry-incubator/bosh-backup-and-restore/ssh"
)

type DeployedInstance struct {
	*instance.DeployedInstance
}

func NewDeployedInstance(instanceGroupName string, remoteRunner ssh.RemoteRunner, logger instance.Logger, jobs orchestrator.Jobs) DeployedInstance {
	return DeployedInstance{
		DeployedInstance: instance.NewDeployedInstance("0", instanceGroupName, "0", remoteRunner, logger, jobs),
	}
}
