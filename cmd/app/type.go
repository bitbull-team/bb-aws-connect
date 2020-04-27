package app

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

// Type print application type
func Type(c *cli.Context) error {

	if c.Bool("raw") {
		fmt.Print(app.GetType())
	} else {
		fmt.Println(app.GetType())
	}

	return nil
}
