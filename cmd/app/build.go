package app

import (
	"github.com/urfave/cli/v2"
)

// NewBuildCommand return "build" command
func NewBuildCommand() *cli.Command {
	return &cli.Command{
		Name:   "build",
		Usage:  "Build application",
		Action: Build,
	}
}

// Build application
func Build(c *cli.Context) error {
	err := app.Build()
	if err != nil {
		return cli.Exit("Error during build process: "+err.Error(), -1)
	}

	return nil
}
