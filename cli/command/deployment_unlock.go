package command

import (
	"github.com/cloudfoundry-incubator/bosh-backup-and-restore/factory"
	"github.com/cloudfoundry-incubator/bosh-backup-and-restore/orchestrator"
	"github.com/urfave/cli"
)

type DeploymentUnlockCommand struct {
}

func NewDeploymentUnlockCommand() DeploymentUnlockCommand {
	return DeploymentUnlockCommand{}
}

func (d DeploymentUnlockCommand) Cli() cli.Command {
	return cli.Command{
		Name:    "unlock",
		Aliases: []string{"u"},
		Usage:   "Unlock a deployment",
		Action:  d.Action,
		Flags:   []cli.Flag{},
	}
}

func (d DeploymentUnlockCommand) Action(c *cli.Context) error {
	trapSigint(true)

	unlocker, err := factory.BuildDeploymentUnlocker(c.Parent().String("target"),
		c.Parent().String("username"),
		c.Parent().String("password"),
		c.Parent().String("ca-cert"),
		c.GlobalBool("debug"),
	)

	if err != nil {
		return processError(orchestrator.NewError(err))
	}

	deployment := c.Parent().String("deployment")
	unlockErr := unlocker.Unlock(deployment)

	if unlockErr.ContainsUnlockOrCleanup() {
		return processErrorWithFooter(unlockErr, backupCleanupAdvisedNotice)
	} else {
		return processError(unlockErr)
	}
}
