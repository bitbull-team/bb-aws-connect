package docker

import (
	"dockerlib"
	"fmt"

	"github.com/urfave/cli/v2"
)

// DeployCommand - Deploy docker image
func DeployCommand(c *cli.Context) error {

	if c.String("repository") != "" {
		dockerlib.LoginToRepository(c.String("repository"), c.String("username"), c.String("password"))
	}

	fmt.Println("Building docker image..")
	err := dockerlib.BuildImage(c.String("cwd"), c.String("tag"), c.StringSlice("build-arg"))
	if err != nil {
		return cli.Exit("Error during Docker image build: "+err.Error(), -1)
	}
	fmt.Println("Docker image built successfully!")
	return nil
}
