package applib

import (
	"configlib"
	"errors"
	"path"
	"shelllib"
)

// NPM application
type NPM struct {
	Application
}

// NewNPM create a new NPM app
func NewNPM(rootPath string) *NPM {
	app := new(NPM)
	app.rootPath = rootPath
	app.config = configlib.NewConfig(ConfigFileName)
	return app
}

// GetPublicPath return the public path
func (app NPM) GetPublicPath() string {
	return path.Join(app.rootPath, "/dist")
}

// GetType return current app type
func (app NPM) GetType() string {
	return "npm"
}

// Build application
func (app NPM) Build() error {
	err := shelllib.ExecuteCommand("npm", "install")
	if err != nil {
		return errors.New("Error executing npm install command: " + err.Error())
	}

	return nil
}
