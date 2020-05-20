package applib

import "configlib"

// WordpressConfig is struct for WordpressConfig application config
type WordpressConfig struct {
	Keys struct {
		AuthKey        string
		SecureAuthKey  string
		LoggedInKey    string
		NonceKey       string
		AuthSalt       string
		SecureAuthSalt string
		LoggedInSalt   string
		NonceSalt      string
	}
}

// Wordpress application
type Wordpress struct {
	Application
}

// NewWordpress create a new Wordpress app
func NewWordpress(rootPath string, configPath string) *Wordpress {
	app := new(Wordpress)
	app.rootPath = rootPath

	// Create config object
	var config Config
	configlib.LoadConfig(configPath, &config)
	app.config = &config

	return app
}

// GetType return current app type
func (app Wordpress) GetType() string {
	return "wordpress"
}
