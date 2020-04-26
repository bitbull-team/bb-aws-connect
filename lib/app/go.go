package applib

import "configlib"

// GO application
type GO struct {
	Application
}

// NewGO create a new GO app
func NewGO(rootPath string) *GO {
	app := new(GO)
	app.rootPath = rootPath
	app.config = configlib.NewConfig("")
	return app
}

// GetType return current app type
func (app GO) GetType() string {
	return "go"
}
