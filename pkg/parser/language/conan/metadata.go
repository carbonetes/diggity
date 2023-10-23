package conan

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/carbonetes/diggity/pkg/model/metadata"
)

func parseConanFileMetadata(attribute string) *[]metadata.ConanMetadata {
	metadataList := new([]metadata.ConanMetadata)
	if len(attribute) == 0 {
		return nil
	}

	packages := strings.Split(attribute, "\n")
	for _, pkg := range packages[1:] {
		if !strings.Contains(pkg, "/") {
			continue
		}

		if strings.Contains(pkg, "[") {
			continue
		}

		properties := strings.Split(pkg, "/")
		metadata := new(metadata.ConanMetadata)

		if len(properties) < 1 {
			continue
		}

		name, version := properties[0], properties[1]

		if name == "" || version == "" {
			continue
		}

		metadata.Name = name
		metadata.Version = version

		*metadataList = append(*metadataList, *metadata)
	}

	return metadataList
}

func parseConanLockMetadata(path string) (*metadata.ConanLockMetadata, error) {
	var metadata metadata.ConanLockMetadata
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(file, &metadata); err != nil {
		return nil, err
	}

	return &metadata, nil
}
