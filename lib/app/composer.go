package applib

import (
	"configlib"
	"errors"
	"shelllib"
)

// Composer application
type Composer struct {
	Application
}

// NewComposer create a new Composer app
func NewComposer(rootPath string) *Composer {
	app := new(Composer)
	app.rootPath = rootPath

	// Create config object
	var config Config
	configlib.LoadConfig("", &config)
	app.config = &config

	return app
}

// GetType return current app type
func (app Composer) GetType() string {
	return "composer"
}

// Build application
func (app Composer) Build() error {
	err := shelllib.ExecuteCommand("composer", "install")
	if err != nil {
		return errors.New("Error executing composer install command: " + err.Error())
	}

	return nil
}
