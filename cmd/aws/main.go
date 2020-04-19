package aws

import (
	"github.com/urfave/cli/v2"
)

// Commands - Return all commands
func Commands() []*cli.Command {
	globalFlags := []cli.Flag{
		&cli.StringFlag{
			Name:     "profile",
			Aliases:  []string{"p"},
			Usage:    "AWS profile name",
			Required: true,
		},
		&cli.StringFlag{
			Name:    "region",
			Aliases: []string{"r"},
			Usage:   "AWS region",
			Value:   "eu-west-1",
		},
	}

	return []*cli.Command{
		{
			Name: "aws",
			Subcommands: []*cli.Command{
				{
					Name:   "connect",
					Usage:  "Connect to an EC2 instance using SSM session",
					Action: SSMListAndStartSession,
					Flags: append(globalFlags, []cli.Flag{
						&cli.StringFlag{
							Name:    "service",
							Aliases: []string{"s"},
							Usage:   "Service Type (example: bastion, frontend, varnish)",
						},
						&cli.StringFlag{
							Name:    "env",
							Aliases: []string{"e"},
							Usage:   "Environment (example: test, stage, prod)",
						},
						&cli.StringFlag{
							Name:    "instance",
							Aliases: []string{"i"},
							Usage:   "Instace ID (example: i-xxxxxxxxxxxxxxxxx)",
						},
					}...),
				},
			},
		},
	}
}
