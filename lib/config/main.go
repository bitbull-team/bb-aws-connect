package configlib

import (
	"net/url"
	"path/filepath"
)

// Config contain CLI configuration
type Config struct {
	AppType string `yaml:"appType" json:"appType"`
}

// Provider load configuration
type Provider interface {
	Load() (string, error)
}

// Parser parse content
type Parser interface {
	Parse(*Config) error
}

// NewConfig create a new Config
func NewConfig(path string) *Config {
	url, _ := url.Parse(path)
	extension := filepath.Ext(path)

	var provider Provider

	// Check for file
	switch url.Scheme {
	case "file":
	default:
		provider = NewFileProvier(path)
	}

	// Load config
	config := new(Config)
	content, errLoad := provider.Load()
	if errLoad != nil {
		return config
	}

	var parser Parser

	// Check for yml file config
	switch extension {
	case ".yml":
		parser = NewYAMLParser(content)
		break
	case ".json":
		parser = NewJSONParser(content)
		break
	default:
		return config
	}

	// Parse content
	errParse := parser.Parse(config)
	if errParse != nil {
		return config
	}

	return config
}
