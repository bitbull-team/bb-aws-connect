package main

import (
	"fmt"
	"log"
	"os"

	"docker"

	"github.com/urfave/cli/v2"
)

func main() {
	cmds := []*cli.Command{}
	cmds = append(cmds, docker.Commands()...)

	cwd, _ := os.Getwd()
	app := &cli.App{
		Name:        "bb-cli",
		Description: "Bitbull CLI",
		Version:     "VERSION", // this will be replace during build phase
		Commands:    cmds,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "cwd",
				Value: cwd,
				Usage: "Current working directory",
			},
		},
		Before: func(c *cli.Context) error {
			if c.String("cwd") != "" {
				os.Chdir(c.String("cwd"))
				cwd, err := os.Getwd()
				if err == nil {
					fmt.Println("Current directory is now: " + cwd)
				}
			}
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
