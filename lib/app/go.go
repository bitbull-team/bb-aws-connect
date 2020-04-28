package applib

import (
	"configlib"
)

// GO application
type GO struct {
	Application
}

// NewGO create a new GO app
func NewGO(rootPath string, configPath string) *GO {
	app := new(GO)
	app.rootPath = rootPath

	// Create config object
	var config Config
	configlib.LoadConfig(configPath, &config)
	app.config = &config

	return app
}

// GetType return current app type
func (app GO) GetType() string {
	return "go"
}
