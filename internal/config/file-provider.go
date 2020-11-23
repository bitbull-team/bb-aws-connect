package config

import (
	"io/ioutil"
)

// FileProvier load config from file
type FileProvier struct {
	path string
}

// Load read file and load content
func (provider FileProvier) Load() (string, error) {
	data, err := ioutil.ReadFile(provider.path)
	if err != nil {
		return string(data), err
	}

	return string(data), nil
}

// NewFileProvier create a FileProvier
func NewFileProvier(path string) *FileProvier {
	provider := new(FileProvier)
	provider.path = path
	return provider
}
