package command

import (
	"github.com/cloudfoundry-incubator/bosh-backup-and-restore/factory"
	"github.com/cloudfoundry-incubator/bosh-backup-and-restore/orchestrator"
	"github.com/urfave/cli"
)

type DeploymentLockCommand struct {
}

func NewDeploymentLockCommand() DeploymentLockCommand {
	return DeploymentLockCommand{}
}

func (d DeploymentLockCommand) Cli() cli.Command {
	return cli.Command{
		Name:    "lock",
		Aliases: []string{"l"},
		Usage:   "Lock a deployment",
		Action:  d.Action,
		Flags:   []cli.Flag{},
	}
}

func (d DeploymentLockCommand) Action(c *cli.Context) error {
	trapSigint(true)

	backuper, err := factory.BuildDeploymentBackuper(c.Parent().String("target"),
		c.Parent().String("username"),
		c.Parent().String("password"),
		c.Parent().String("ca-cert"),
		c.Bool("with-manifest"),
		c.GlobalBool("debug"),
	)

	if err != nil {
		return processError(orchestrator.NewError(err))
	}

	deployment := c.Parent().String("deployment")
	backupErr := backuper.Backup(deployment, c.String("artifact-path"))

	if backupErr.ContainsUnlockOrCleanup() {
		return processErrorWithFooter(backupErr, backupCleanupAdvisedNotice)
	} else {
		return processError(backupErr)
	}
}
