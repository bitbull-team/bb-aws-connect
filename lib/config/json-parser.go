package configlib

import "encoding/json"

// JSONParser parse YML format
type JSONParser struct {
	content string
}

// Parse convert string into config object
func (parser JSONParser) Parse(config *Config) error {
	err := json.Unmarshal([]byte(parser.content), &config)
	if err != nil {
		return err
	}

	return nil
}

// NewJSONParser create a JSONParser
func NewJSONParser(content string) *JSONParser {
	parser := new(JSONParser)
	parser.content = content
	return parser
}
