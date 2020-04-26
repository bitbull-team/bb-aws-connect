package filesystemlib

import (
	"encoding/json"
	"io/ioutil"
)

// ComposerFile is a parsed composer.json
type ComposerFile struct {
	Name            string            `json:"name"`
	Description     string            `json:"description"`
	Dependencies    map[string]string `json:"require"`
	DevDependencies map[string]string `json:"require-dev"`
	Extra           map[string]string `json:"extra"`
}

// HasDependency return true o false if depency exist
func (composer ComposerFile) HasDependency(dependency string) bool {
	return composer.Dependencies[dependency] != ""
}

// LoadComposerFile load composer.json file
func LoadComposerFile(filepath string) (ComposerFile, error) {
	data := ComposerFile{}

	// Read file
	file, readErr := ioutil.ReadFile(filepath)
	if readErr != nil {
		return data, readErr
	}

	// Parse content
	unmarshalErr := json.Unmarshal([]byte(file), &data)
	if unmarshalErr != nil {
		return data, unmarshalErr
	}

	return data, nil
}
