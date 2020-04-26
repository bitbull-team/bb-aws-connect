package applib

import (
	"configlib"
	"errors"
	"path"
	"shelllib"
)

// YARN application
type YARN struct {
	Application
}

// NewYARN create a new YARN app
func NewYARN(rootPath string) *YARN {
	app := new(YARN)
	app.rootPath = rootPath
	app.config = configlib.NewConfig(ConfigFileName)
	return app
}

// GetPublicPath return the public path
func (app YARN) GetPublicPath() string {
	return path.Join(app.rootPath, "/dist")
}

// GetType return current app type
func (app YARN) GetType() string {
	return "yarn"
}

// Build application
func (app YARN) Build() error {
	err := shelllib.ExecuteCommand("yarn")
	if err != nil {
		return errors.New("Error executing yarn command: " + err.Error())
	}

	return nil
}
