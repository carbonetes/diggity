package cocoapods

import (
	"github.com/carbonetes/diggity/internal/log"
	"gopkg.in/yaml.v3"
)

// Swift and Objective-C Podfile Lock Metadata
type FileLockMetadata struct {
	Pods          []interface{}     `yaml:"PODS"`
	Dependencies  []string          `yaml:"DEPENDENCIES"`
	SpecChecksums map[string]string `yaml:"SPEC CHECKSUMS"`
	Cocoapods     string            `yaml:"COCOAPODS"`
}

type FileLockMetadataCheckSums struct {
	Checksums string `json:"checksum"`
}

func readManifestFile(content []byte) FileLockMetadata {
	var metadata FileLockMetadata
	err := yaml.Unmarshal(content, &metadata)
	if err != nil {
		log.Error("Failed to unmarshal Podfile.lock")
	}
	return metadata
}
