package ssm

import (
	"github.com/bitbull-team/bb-aws-connect/internal/config"
	"github.com/urfave/cli/v2"
)

// Config is struct for AWS
type Config struct {
	Profile string
	Region  string
	SSM     struct {
		User  string
		Shell string
		Cwd   string
	}
}

var globalConfig Config

// Commands - Return all commands
func Commands(globalFlags []cli.Flag) []*cli.Command {
	return []*cli.Command{
		{
			Name:  "ssm",
			Usage: "AWS SSM Commands",
			Before: func(c *cli.Context) error {
				config.LoadConfig(c.String("config"), &globalConfig)
				if len(c.String("region")) > 0 {
					globalConfig.Region = c.String("region")
				}
				if len(c.String("profile")) > 0 {
					globalConfig.Profile = c.String("profile")
				}
				return nil
			},
			Flags: globalFlags,
			Subcommands: []*cli.Command{
				NewConnectCommand(globalFlags),
				NewRunCommand(globalFlags),
				NewRunDocumentCommand(globalFlags),
				NewTunnelCommand(globalFlags),
			},
		},
	}
}
