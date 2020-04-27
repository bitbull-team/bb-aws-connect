package applib

import "configlib"

// Nuxt application
type Nuxt struct {
	NPM
}

// NewNuxt create a new Nuxt app
func NewNuxt(rootPath string) *Nuxt {
	app := new(Nuxt)
	app.rootPath = rootPath

	// Create config object
	var config Config
	configlib.LoadConfig("", &config)
	app.config = &config

	return app
}

// GetType return current app type
func (app Nuxt) GetType() string {
	return "nuxt"
}
