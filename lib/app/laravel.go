package applib

import (
	"configlib"
	"path"
)

// Laravel application
type Laravel struct {
	Composer
}

// NewLaravel create a new Laravel app
func NewLaravel(rootPath string) *Laravel {
	app := new(Laravel)
	app.rootPath = rootPath
	app.config = configlib.NewConfig(ConfigFileName)
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
