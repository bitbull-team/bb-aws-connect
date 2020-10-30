package configlib

import (
	"errors"
	"net/url"
	"path/filepath"
)

// Config contain CLI configuration
type Config struct{}

// Provider load configuration
type Provider interface {
	Load() (string, error)
}

// Parser parse content
type Parser interface {
	Parse(config interface{}) error
}

// LoadConfig load config into existing object
func LoadConfig(path string, config interface{}) error {
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
	content, errLoad := provider.Load()
	if errLoad != nil {
		return errLoad
	}

	var parser Parser

	// Check for yml file config
	switch extension {
	case ".yml":
		parser = NewYAMLParser(content)
	case ".yaml":
		parser = NewYAMLParser(content)
	case ".json":
		parser = NewJSONParser(content)
	default:
		return errors.New("No parser found for config file extension " + extension)
	}

	// Parse content
	errParse := parser.Parse(config)
	if errParse != nil {
		return errParse
	}

	return nil
}
