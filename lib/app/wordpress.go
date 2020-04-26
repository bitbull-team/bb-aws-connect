package applib

import "configlib"

// Wordpress application
type Wordpress struct {
	Application
}

// NewWordpress create a new Wordpress app
func NewWordpress(rootPath string) *Wordpress {
	app := new(Wordpress)
	app.rootPath = rootPath
	app.config = configlib.NewConfig("")
	return app
}

// GetType return current app type
func (app Wordpress) GetType() string {
	return "wordpress"
}
