package docker

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

// DeployCommand - Deploy docker image
func DeployCommand(c *cli.Context) error {
	fmt.Println("Hi from Docker Build")
	return nil
}
