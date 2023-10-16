package pub

import (
	"os"

	"github.com/carbonetes/diggity/pkg/model/metadata"
	"gopkg.in/yaml.v3"
)

type Metadata map[string]interface{}

func parseMetadata(path string) (*Metadata, error) {
	metadata := make(Metadata)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal([]byte(data), &metadata); err != nil {
		return nil, err
	}

	return &metadata, nil
}

func parseLockfileMetadata(path string) (*metadata.PubspecLockPackage, error) {
	var metadata metadata.PubspecLockPackage
	byteValue, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal([]byte(byteValue), &metadata); err != nil {
		return nil, err
	}
	return &metadata, nil
}
