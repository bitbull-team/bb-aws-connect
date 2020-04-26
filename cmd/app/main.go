package app

import (
	"applib"
	"encoding/json"
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

var app applib.ApplicationInterface

// Commands - Return all commands
func Commands() []*cli.Command {
	return []*cli.Command{
		{
			Name: "app",
			Subcommands: []*cli.Command{
				{
					Name:  "type",
					Usage: "Return the current app type",
					Action: func(c *cli.Context) error {
						fmt.Println(app.GetType())
						return nil
					},
				},
				{
					Name:  "config",
					Usage: "Dump config",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:    "format",
							Aliases: []string{"f"},
							Usage:   "Output format (example: json, yml, go)",
							Value:   "json",
						},
					},
					Action: func(c *cli.Context) error {
						config := app.GetConfig()
						var dump []byte

						format := c.String("format")
						switch format {
						case "json":
							dump, _ = json.Marshal(config)
							break
						case "yml":
						case "yaml":
							dump, _ = yaml.Marshal(config)
							break
						case "go":
							dump = []byte(fmt.Sprintf("%+v", config))
							break
						default:
							return cli.Exit("Format not recognized: "+format, -1)
						}

						fmt.Printf("%+s", dump)
						return nil
					},
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
			Before: func(c *cli.Context) error {
				cwd, _ := os.Getwd()
				app = applib.NewApplication(cwd)
				return nil
			},
		},
	}
}
