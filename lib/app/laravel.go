package applib

import (
	"configlib"
	"path"
)

// LaravelConfig is struct for laravel application config
type LaravelConfig struct {
	Key string
}

// Laravel application
type Laravel struct {
	Composer
}

// NewLaravel create a new Laravel app
func NewLaravel(rootPath string) *Laravel {
	app := new(Laravel)
	app.rootPath = rootPath

	// Create config object
	var config Config
	configlib.LoadConfig("", &config)
	app.config = &config

	return app
}

// GetPublicPath return the public path
func (app Laravel) GetPublicPath() string {
	return path.Join(app.rootPath, "/public")
}

// GetType return current app type
func (app Laravel) GetType() string {
	return "laravel"
}
