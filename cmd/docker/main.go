package docker

import (
	"github.com/urfave/cli/v2"
)

// Commands - Return all commands
func Commands() []*cli.Command {
	return []*cli.Command{
		{
			Name:   "docker:deploy",
			Action: DeployCommand,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "username",
					Aliases: []string{"u"},
					Usage:   "username for Docker repository",
				},
				&cli.StringFlag{
					Name:    "password",
					Aliases: []string{"p"},
					Usage:   "password for Docker repository",
				},
				&cli.StringFlag{
					Name:  "repository",
					Usage: "Docker repository URL",
				},
				&cli.StringSliceFlag{
					Name:  "build-arg",
					Usage: "Docker build arguments",
				},
				&cli.StringFlag{
					Name:    "tag",
					Aliases: []string{"t"},
					Usage:   "Docker image tag",
				},
			},
		},
	}
}
