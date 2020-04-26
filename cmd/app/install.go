package app

import (
	"github.com/urfave/cli/v2"
)

// Install application
func Install(c *cli.Context) error {
	err := app.Install()
	if err != nil {
		return cli.Exit("Error during install process: "+err.Error(), -1)
	}

	return nil
}
