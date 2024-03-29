package ecs

import (
	"github.com/bitbull-team/bb-aws-connect/internal/config"
	"github.com/urfave/cli/v2"
)

// Config is struct for AWS
type Config struct {
	Profile string
	Region  string
	ECS     struct {
		Cluster string
		Session struct {
			User    string
			Shell   string
			Workdir string
		}
	}
}

// Global config
var globalConfig Config

// Commands - Return all commands
func Commands(globalFlags []cli.Flag) []*cli.Command {
	return []*cli.Command{
		{
			Name:  "ecs",
			Usage: "AWS ECS Commands",
			Before: func(c *cli.Context) error {
				config.LoadConfig(c.String("config"), &globalConfig)
				// Store args into config, child commands cannot access it
				globalConfig.Region = c.String("region")
				globalConfig.Profile = c.String("profile")
				return nil
			},
			Flags: globalFlags,
			Subcommands: []*cli.Command{
				NewConnectCommand(globalFlags),
			},
		},
	}
}
