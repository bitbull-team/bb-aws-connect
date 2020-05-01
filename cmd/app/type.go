package app

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

// NewTypeCommand return "type" command
func NewTypeCommand() *cli.Command {
	return &cli.Command{
		Name:  "type",
		Usage: "Return the current app type",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "raw",
				Aliases: []string{"r"},
				Usage:   "Print value without newline",
			},
		},
		Action: Type,
	}
}

// Type print application type
func Type(c *cli.Context) error {

	if c.Bool("raw") {
		fmt.Print(app.GetType())
	} else {
		fmt.Println(app.GetType())
	}

	return nil
}
