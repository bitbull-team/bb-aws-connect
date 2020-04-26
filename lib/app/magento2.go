package applib

import (
	"configlib"
	"path"
)

// Magento2 application
type Magento2 struct {
	Composer
}

// NewMagento2 create a new Magento2 app
func NewMagento2(rootPath string) *Magento2 {
	app := new(Magento2)
	app.rootPath = rootPath
	app.config = configlib.NewConfig(ConfigFileName)
	return app
}

// GetPublicPath return the public path
func (app Magento2) GetPublicPath() string {
	return path.Join(app.rootPath, "/pub")
}

// GetType return current app type
func (app Magento2) GetType() string {
	return "magento2"
}
