package applib

import (
	"configlib"
	"filesystemlib"
	"fmt"
)

// ApplicationInterface is the base application type
type ApplicationInterface interface {
	GetPublicPath() string
	GetType() string
	GetConfig() *configlib.Config
	Install() error
	Build() error
}

// Application is the base application type
type Application struct {
	rootPath string
	config   *configlib.Config
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
func (app Application) GetConfig() *configlib.Config {
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
func NewApplication(rootPath string) ApplicationInterface {
	config := configlib.NewConfig("")

	// Check if config override app type
	if len(config.App.Type) != 0 {
		switch config.App.Type {
		case "wordpress":
			return *NewWordpress(rootPath)
		case "magento":
			return *NewMagento(rootPath)
		case "laravel":
			return *NewLaravel(rootPath)
		case "magento2":
			return *NewMagento2(rootPath)
		case "magentowp":
			return *NewMagentoWordpress(rootPath)
		case "nuxt":
			return *NewNuxt(rootPath)
		case "go":
			return *NewGO(rootPath)
		case "composer":
			return *NewComposer(rootPath)
		case "yarn":
			return *NewYARN(rootPath)
		case "npm":
			return *NewNPM(rootPath)
		}
	}

	// Discover app type

	if filesystemlib.FileExist("wp-load.php") {
		return *NewWordpress(rootPath)
	}

	if filesystemlib.FileExist("app/Mage.php") {
		if filesystemlib.FileExist("wp/wp-load.php") || filesystemlib.FileExist("blog/wp-load.php") {
			return NewMagentoWordpress(rootPath)
		}
		return *NewMagento(rootPath)
	}

	if filesystemlib.FileExist("composer.json") {
		composer, _ := filesystemlib.LoadComposerFile("composer.json")

		if composer.HasDependency("laravel/framework") {
			return *NewLaravel(rootPath)
		}

		if composer.HasDependency("laravel/lumen-framework") {
			return *NewLaravel(rootPath)
		}

		if composer.HasDependency("magento/product-community-edition") {
			return *NewMagento2(rootPath)
		}

		if composer.HasDependency("magento/product-enterprise-edition") {
			return *NewMagento2(rootPath)
		}

		if composer.HasDependency("magento-hackathon/magento-composer-installer") {
			if composer.HasDependency("wordpress") || composer.HasDependency("wordpress/core") {
				return *NewMagentoWordpress(rootPath)
			}

			if filesystemlib.FileExist("wp/wp-load.php") || filesystemlib.FileExist("blog/wp-load.php") {
				return *NewMagentoWordpress(rootPath)
			}

			return *NewMagento(rootPath)
		}
	}

	if filesystemlib.FileExist("package.json") {
		npm, _ := filesystemlib.LoadNPMPackageFile("package.json")

		if npm.HasDependency("nuxt") {
			return *NewNuxt(rootPath)
		}
	}

	if filesystemlib.FileExist("go.mod") {
		return *NewGO(rootPath)
	}

	if filesystemlib.FileExist("composer.json") {
		return *NewComposer(rootPath)
	}

	if filesystemlib.FileExist("package.json") {
		if filesystemlib.FileExist("yarn.lock") {
			return *NewYARN(rootPath)
		}
		return *NewNPM(rootPath)
	}

	var app Application
	app.rootPath = rootPath
	app.config = configlib.NewConfig("")
	return app
}
