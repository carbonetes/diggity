package slsa

import (
	"encoding/json"
	"os"
)

var ProvenanceMetadata map[string]interface{}

func GetProvenanceMetadata(location string) (map[string]interface{}, error) {
	file, err := os.ReadFile(location)
	if err != nil {
		return nil, err
	}

	// Validate if file is valid JSON
	if !json.Valid(file) {
		return nil, nil
	}

	if err = json.Unmarshal(file, &ProvenanceMetadata); err != nil {
		return nil, err
	}

	return ProvenanceMetadata, nil
}
