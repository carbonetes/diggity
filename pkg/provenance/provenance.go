package provenance

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/carbonetes/diggity/internal/helper"
)

func Parse(file string) (map[string]interface{}, error) {
	found, err := helper.IsFileExists(file)
	if err != nil {
		return nil, err
	}

	if !found {
		return nil, fmt.Errorf("File not found! - %s", file)
	}

	bytes, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	if !json.Valid(bytes) {
		return nil, fmt.Errorf("Invalid metdata! - %s", file)
	}
	var metadata map[string]interface{}
	if err = json.Unmarshal(bytes, &metadata); err != nil {
		return nil, fmt.Errorf("Failed to unmarshal json file! - %s", file)
	}

	if !validate(metadata) {
		return nil, fmt.Errorf("%s contains invalid provenance metadata.", file)
	}

	return metadata, nil

}

func validate(metadata map[string]interface{}) bool {
	if _, ok := metadata["_type"]; !ok {
		return false
	}
	if _, ok := metadata["subject"]; !ok {
		return false
	}
	if _, ok := metadata["predicateType"]; !ok {
		return false
	}
	if _, ok := metadata["predicate"]; !ok {
		return false
	}
	if !strings.Contains(metadata["_type"].(string), "https://in-toto.io/Statement") {
		return false
	}
	return true
}
