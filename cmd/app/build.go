package app

import (
	"github.com/urfave/cli/v2"
)

// Build application
func Build(c *cli.Context) error {
	err := app.Build()
	if err != nil {
		return cli.Exit("Error during build process: "+err.Error(), -1)
	}

	return nil
}
