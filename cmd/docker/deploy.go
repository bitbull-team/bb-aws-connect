package docker

import (
	"dockerlib"
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
)

// NewDeployCommand return "deploy" command
func NewDeployCommand() *cli.Command {
	return &cli.Command{
		Name:   "docker:deploy",
		Action: Deploy,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "username",
				Aliases: []string{"u"},
				Usage:   "username for Docker registry",
			},
			&cli.StringFlag{
				Name:    "password",
				Aliases: []string{"p"},
				Usage:   "password for Docker registry",
			},
			&cli.StringFlag{
				Name:  "registry",
				Usage: "Docker registry URL",
			},
			&cli.StringSliceFlag{
				Name:  "build-arg",
				Usage: "Docker build arguments",
			},
			&cli.StringFlag{
				Name:     "image",
				Aliases:  []string{"i"},
				Usage:    "Docker image name, can contain :tag",
				Required: true,
			},
			&cli.StringFlag{
				Name:    "tag",
				Aliases: []string{"t"},
				Usage:   "Docker image tag, override tag through with image:tag",
			},
		},
	}
}

// Deploy docker image
func Deploy(c *cli.Context) error {
	var err error

	// Login to registry
	if c.String("registry") != "" {
		fmt.Println("Login to registry..")
		err = dockerlib.LoginToRegistry(c.String("registry"), c.String("username"), c.String("password"))
		if err != nil {
			return cli.Exit("Error during Docker registry login: "+err.Error(), -1)
		}
		fmt.Println("Logged to Docker registry!")
	}

	// Check if tag is provided through image name
	if len(c.String("tag")) == 0 {
		imagesParts := strings.Split(c.String("image"), ":")
		if len(imagesParts) == 2 {
			c.Set("tag", imagesParts[1])
		} else {
			c.Set("tag", "latest")
		}
	}

	// Build Docker image
	fmt.Println("Building docker image..")
	err = dockerlib.BuildImage(c.String("cwd"), c.String("image"), c.String("tag"), c.StringSlice("build-arg"))
	if err != nil {
		return cli.Exit("Error during Docker image build: "+err.Error(), -1)
	}
	fmt.Println("Docker image built successfully!")

	// Push Docker image
	fmt.Println("Pushing docker image..")
	err = dockerlib.PushImage(c.String("image"), c.String("tag"))
	if err != nil {
		return cli.Exit("Error during Docker push build: "+err.Error(), -1)
	}
	fmt.Println("Docker image pushed successfully!")

	return nil
}
