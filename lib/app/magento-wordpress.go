package applib

import (
	"configlib"
	"path"
)

// MagentoWordpressConfig is struct for Magento&Wordpress application config
type MagentoWordpressConfig struct {
	Config

	Magento   MagentoConfig
	Wordpress WordpressConfig
}

// MagentoWordpress application
type MagentoWordpress struct {
	Composer
}

// NewMagentoWordpress create a new MagentoWordpress app
func NewMagentoWordpress(rootPath string, configPath string) *MagentoWordpress {
	app := new(MagentoWordpress)
	app.rootPath = rootPath

	// Create config object
	var config Config
	configlib.LoadConfig(configPath, &config)
	app.config = &config

	return app
}

// GetPublicPath return the public path
func (app MagentoWordpress) GetPublicPath() string {
	return path.Join(app.rootPath, "/public")
}

// GetType return current app type
func (app MagentoWordpress) GetType() string {
	return "magento-wordpress"
}
