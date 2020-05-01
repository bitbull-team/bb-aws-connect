package app

import (
	"encoding/json"
	"fmt"

	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

// NewConfigCommand return "config" command
func NewConfigCommand() *cli.Command {
	return &cli.Command{
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
	}
}

// DumpConfig print application configurations
func DumpConfig(c *cli.Context) error {
	config := app.GetConfig()
	var dump []byte

	format := c.String("format")
	switch format {
	case "json":
		dump, _ = json.Marshal(config)
	case "yml":
		dump, _ = yaml.Marshal(config)
	case "yaml":
		dump, _ = yaml.Marshal(config)
	case "go":
		dump = []byte(fmt.Sprintf("%+v", config))
	default:
		return cli.Exit("Format not recognized: "+format, 1)
	}

	fmt.Printf("%+s", dump)
	return nil
}
