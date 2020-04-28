package applib

import (
	"configlib"
	"filesystemlib"
	"fmt"
)

// Config is struct for application config
type Config struct {
	App struct {
		Type       string
		Env        string
		PublicPath string
	}

	Magento   MagentoConfig
	Wordpress WordpressConfig
	Laravel   LaravelConfig
}

// ApplicationInterface is the base application type
type ApplicationInterface interface {
	GetPublicPath() string
	GetType() string
	GetConfig() *Config
	Install() error
	Build() error
}

// Application is the base application type
type Application struct {
	rootPath string
	config   *Config
}

// GetPublicPath return the public path
func (app Application) GetPublicPath() string {
	return app.rootPath
}

// GetType return current app type
func (Application) GetType() string {
	return "base"
}

// GetConfig load config property
func (app Application) GetConfig() *Config {
	return app.config
}

// Install new application
func (Application) Install() error {
	fmt.Println("Install not implemented, skipping..")
	return nil
}

// Build application
func (Application) Build() error {
	fmt.Println("Build not implemented, skipping..")
	return nil
}

// NewApplication create an Application struct
func NewApplication(rootPath string, configPath string) ApplicationInterface {
	var config Config
	configlib.LoadConfig(configPath, &config)

	// Check if config override app type
	if len(config.App.Type) != 0 {
		switch config.App.Type {
		case "wordpress":
			return *NewWordpress(rootPath, configPath)
		case "magento":
			return *NewMagento(rootPath, configPath)
		case "laravel":
			return *NewLaravel(rootPath, configPath)
		case "magento2":
			return *NewMagento2(rootPath, configPath)
		case "magentowp":
			return *NewMagentoWordpress(rootPath, configPath)
		case "nuxt":
			return *NewNuxt(rootPath, configPath)
		case "go":
			return *NewGO(rootPath, configPath)
		case "composer":
			return *NewComposer(rootPath, configPath)
		case "yarn":
			return *NewYARN(rootPath, configPath)
		case "npm":
			return *NewNPM(rootPath, configPath)
		}
	}

	// Discover app type

	if filesystemlib.FileExist("wp-load.php") {
		return *NewWordpress(rootPath, configPath)
	}

	if filesystemlib.FileExist("app/Mage.php") {
		if filesystemlib.FileExist("wp/wp-load.php") || filesystemlib.FileExist("blog/wp-load.php") {
			return NewMagentoWordpress(rootPath, configPath)
		}
		return *NewMagento(rootPath, configPath)
	}

	if filesystemlib.FileExist("go.mod") {
		return *NewGO(rootPath, configPath)
	}

	// Elaborate composer.json dependencies
	if filesystemlib.FileExist("composer.json") {
		composer, _ := filesystemlib.LoadComposerFile("composer.json")

		if composer.HasDependency("laravel/framework") {
			return *NewLaravel(rootPath, configPath)
		}

		if composer.HasDependency("laravel/lumen-framework") {
			return *NewLaravel(rootPath, configPath)
		}

		if composer.HasDependency("magento/product-community-edition") {
			return *NewMagento2(rootPath, configPath)
		}

		if composer.HasDependency("magento/product-enterprise-edition") {
			return *NewMagento2(rootPath, configPath)
		}

		if composer.HasDependency("magento-hackathon/magento-composer-installer") {
			if composer.HasDependency("wordpress") || composer.HasDependency("wordpress/core") {
				return *NewMagentoWordpress(rootPath, configPath)
			}

			if filesystemlib.FileExist("wp/wp-load.php") || filesystemlib.FileExist("blog/wp-load.php") {
				return *NewMagentoWordpress(rootPath, configPath)
			}

			return *NewMagento(rootPath, configPath)
		}
	}

	// Elaborate package.json dependencies
	if filesystemlib.FileExist("package.json") {
		npm, _ := filesystemlib.LoadNPMPackageFile("package.json")

		if npm.HasDependency("nuxt") {
			return *NewNuxt(rootPath, configPath)
		}
	}

	// Defaults for composer and package.json

	if filesystemlib.FileExist("composer.json") {
		return *NewComposer(rootPath, configPath)
	}

	if filesystemlib.FileExist("package.json") {
		if filesystemlib.FileExist("yarn.lock") {
			return *NewYARN(rootPath, configPath)
		}
		return *NewNPM(rootPath, configPath)
	}

	// Base application
	var app Application
	app.rootPath = rootPath
	app.config = &config

	return app
}
