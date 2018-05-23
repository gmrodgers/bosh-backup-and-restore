package factory

import (
	"github.com/cloudfoundry-incubator/bosh-backup-and-restore/bosh"
	"github.com/cloudfoundry-incubator/bosh-backup-and-restore/orchestrator"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

func BuildBoshDeploymentManager(targetUrl, username, password, caCert string, logger boshlog.Logger) (orchestrator.DeploymentManager, error) {
	boshClient, err := bosh.BuildClient(targetUrl, username, password, caCert, logger)
	if err != nil {
		return nil, err
	}

	return bosh.NewDeploymentManager(boshClient, logger), nil
}
