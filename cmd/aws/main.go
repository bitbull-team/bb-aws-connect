package aws

import (
	"github.com/urfave/cli/v2"
)

// Commands - Return all commands
func Commands() []*cli.Command {
	return []*cli.Command{
		{
			Name: "aws",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "profile",
					Aliases: []string{"p"},
					Usage:   "AWS profile",
				},
				&cli.StringFlag{
					Name:    "region",
					Aliases: []string{"r"},
					Usage:   "AWS region",
					Value:   "eu-west-1",
				},
				&cli.StringFlag{
					Name:    "env",
					Aliases: []string{"e"},
					Usage:   "Target environment",
				},
			},
			Subcommands: []*cli.Command{
				{
					Name:   "ssm",
					Usage:  "AWS profile",
					Action: SSMListAvailableInstances,
				},
			},
		},
	}
}
