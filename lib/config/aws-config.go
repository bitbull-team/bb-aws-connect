package configlib

// AWSConfig is struct for AWS
type AWSConfig struct {
	Profile string `yaml:"profile" json:"profile"`
	Region  string `yaml:"region" json:"region"`
}
