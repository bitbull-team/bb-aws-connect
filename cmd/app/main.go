package app

import (
	"applib"
	"os"

	"github.com/urfave/cli/v2"
)

var app applib.ApplicationInterface

// Commands - Return all commands
func Commands() []*cli.Command {
	return []*cli.Command{
		{
			Name: "app",
			Before: func(c *cli.Context) error {
				cwd, _ := os.Getwd()
				app = applib.NewApplication(cwd)
				return nil
			},
			Subcommands: []*cli.Command{
				{
					Name:  "type",
					Usage: "Return the current app type",
					Flags: []cli.Flag{
						&cli.BoolFlag{
							Name:    "raw",
							Aliases: []string{"r"},
							Usage:   "Print value without newline",
						},
					},
					Action: Type,
				},
				{
					Name:  "config",
					Usage: "Dump config",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:    "format",
							Aliases: []string{"f"},
							Usage:   "Output format (example: json, yml, go)",
							Value:   "yml",
						},
					},
					Action: DumpConfig,
				},
				{
					Name:   "install",
					Usage:  "Install application",
					Action: Install,
				},
				{
					Name:   "build",
					Usage:  "Build application",
					Action: Build,
				},
			},
		},
	}
}
