package composer

import (
	"encoding/json"
	"os"

	"github.com/carbonetes/diggity/pkg/model/metadata"
)

// Parse composer package metadata
// Metadata model based from https://github.com/composer/composer/blob/main/src/Composer/Package/Locker.php
func parseMetadata(path string) (*metadata.ComposerMetadata, error) {
	metadata := new(metadata.ComposerMetadata)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, metadata); err != nil {
		return nil, err
	}
	return metadata, nil
}
