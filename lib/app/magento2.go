package applib

import "configlib"

// Magento2 application
type Magento2 struct {
	Composer
}

// NewMagento2 create a new Magento2 app
func NewMagento2(rootPath string) *Magento2 {
	app := new(Magento2)
	app.rootPath = rootPath

	// Create config object
	var config Config
	configlib.LoadConfig("", &config)
	app.config = &config

	return app
}

// GetType return current app type
func (app Magento2) GetType() string {
	return "magento"
}
