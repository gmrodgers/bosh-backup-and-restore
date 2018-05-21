package main

import (
	"os"

	"github.com/urfave/cli"

	"github.com/cloudfoundry-incubator/bosh-backup-and-restore/cli/command"
	"github.com/cloudfoundry-incubator/bosh-backup-and-restore/cli/flags"
)

var version string

func main() {
	cli.AppHelpTemplate = helpTextTemplate

	app := cli.NewApp()

	app.Version = version
	app.Name = "db-lock"
	app.Usage = "BOSH Deployment Lock for Database Upgrades"
	app.HideHelp = true

	app.Commands = []cli.Command{
		{
			Name:   "deployment",
			Usage:  "Lock/Unlock BOSH deployments",
			Flags:  availableDeploymentFlags(),
			Before: validateDeploymentFlags,
			Subcommands: []cli.Command{
				command.NewDeploymentLockCommand().Cli(),
				command.NewDeploymentUnlockCommand().Cli(),
			},
		},
		{
			Name:    "help",
			Aliases: []string{"h"},
			Usage:   "Shows a list of commands or help for one command",
			Action:  versionAction,
		},
		{
			Name:    "version",
			Aliases: []string{"v"},
			Usage:   "Shows the version",
			Action: func(c *cli.Context) error {
				cli.ShowVersion(c)
				return nil
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		os.Exit(1)
	}
}

func versionAction(c *cli.Context) error {
	cli.ShowAppHelp(c)
	return nil
}

func validateDeploymentFlags(c *cli.Context) error {
	return flags.Validate([]string{"target", "username", "password", "deployment"}, c)
}

func availableDeploymentFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "target, t",
			Value: "",
			Usage: "Target BOSH Director URL",
		},
		cli.StringFlag{
			Name:  "username, u",
			Value: "",
			Usage: "BOSH Director username",
		},
		cli.StringFlag{
			Name:   "password, p",
			Value:  "",
			EnvVar: "BOSH_CLIENT_SECRET",
			Usage:  "BOSH Director password",
		},
		cli.StringFlag{
			Name:  "deployment, d",
			Value: "",
			Usage: "Name of BOSH deployment",
		},
		cli.BoolFlag{
			Name:  "debug",
			Usage: "Enable debug logs",
		},
		cli.StringFlag{
			Name:   "ca-cert",
			Value:  "",
			EnvVar: "CA_CERT",
			Usage:  "Custom CA certificate",
		},
	}
}
