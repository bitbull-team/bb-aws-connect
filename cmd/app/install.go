package app

import (
	"github.com/urfave/cli/v2"
)

// NewInstallCommand return "build" command
func NewInstallCommand() *cli.Command {
	return &cli.Command{
		Name:   "install",
		Usage:  "Install application",
		Action: Install,
	}
}

// Install application
func Install(c *cli.Context) error {
	err := app.Install()
	if err != nil {
		return cli.Exit("Error during install process: "+err.Error(), 1)
	}

	return nil
}
