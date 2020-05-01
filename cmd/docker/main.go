package docker

import (
	"github.com/urfave/cli/v2"
)

// Commands - Return all commands
func Commands() []*cli.Command {
	return []*cli.Command{
		NewDeployCommand(),
	}
}
