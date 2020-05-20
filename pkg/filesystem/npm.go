package filesystemlib

import (
	"encoding/json"
	"io/ioutil"
)

// NPMPackageFile is a parsed package.json
type NPMPackageFile struct {
	Name            string            `json:"name"`
	Description     string            `json:"description"`
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

// HasDependency return true o false if depency exist
func (npm NPMPackageFile) HasDependency(dependency string) bool {
	return npm.Dependencies[dependency] != ""
}

// LoadNPMPackageFile load package.json file
func LoadNPMPackageFile(filepath string) (NPMPackageFile, error) {
	data := NPMPackageFile{}

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
