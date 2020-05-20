package applib

import (
	"configlib"
)

// MagentoConfig is struct for Magento application config
type MagentoConfig struct {
	Key string
}

// Magento application
type Magento struct {
	Composer
}

// NewMagento create a new Magento app
func NewMagento(rootPath string, configPath string) *Magento {
	app := new(Magento)
	app.rootPath = rootPath

	// Create config object
	var config Config
	configlib.LoadConfig(configPath, &config)
	app.config = &config

	return app
}

// GetType return current app type
func (app Magento) GetType() string {
	return "magento"
}
