package app

import (
	"applib"

	"github.com/urfave/cli/v2"
)

// Global application
var app applib.ApplicationInterface

// Commands - Return all commands
func Commands() []*cli.Command {
	return []*cli.Command{
		{
			Name: "app",
			Before: func(c *cli.Context) error {
				app = applib.NewApplication(c.String("root"), c.String("config"))
				return nil
			},
			Subcommands: []*cli.Command{
				NewTypeCommand(),
				NewConfigCommand(),
				NewInstallCommand(),
				NewBuildCommand(),
			},
		},
	}
}
