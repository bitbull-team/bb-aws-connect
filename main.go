package main

import (
	"log"
	"os"

	"docker"

	"github.com/urfave/cli/v2"
)

func main() {
	cmds := []*cli.Command{}

	cmds = append(cmds, docker.Commands()...)

	app := &cli.App{
		Name:        "bb-cli",
		Description: "Bitbull CLI",
		Version:     "VERSION", // this will be replace during build phase
		Commands:    cmds,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
