package main

import (
	"fmt"
	"log"
	"os"

	"aws"
	"docker"

	"github.com/urfave/cli/v2"
)

func main() {
	cmds := []*cli.Command{}
	cmds = append(cmds, docker.Commands()...)
	cmds = append(cmds, aws.Commands()...)

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
				_, err := os.Getwd()
				if err != nil {
					fmt.Println("Cannot change CWD: " + err.Error())
				}
			}
			return nil
		},
		EnableBashCompletion: true,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
