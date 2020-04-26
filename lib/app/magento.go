package applib

import "configlib"

// Magento application
type Magento struct {
	Composer
}

// NewMagento create a new Magento app
func NewMagento(rootPath string) *Magento {
	app := new(Magento)
	app.rootPath = rootPath
	app.config = configlib.NewConfig(ConfigFileName)
	return app
}

// GetType return current app type
func (app Magento) GetType() string {
	return "magento"
}
