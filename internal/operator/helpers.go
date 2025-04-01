package operator

import (
	"encoding/base64"
	"gopkg.in/yaml.v3"
)

// getValues unmarshals the base64 encoded YAML string into a map
func getValues(configuration string) (map[string]interface{}, error) {
	// Decode the base64 encoded YAML string
	decodedYAML, err := base64.StdEncoding.DecodeString(configuration)
	if err != nil {
		return nil, err
	}

	// Unmarshal the YAML into a map
	vals := map[string]interface{}{}

	err = yaml.Unmarshal(decodedYAML, &vals)
	if err != nil {
		return nil, err
	}

	return vals, nil
}
