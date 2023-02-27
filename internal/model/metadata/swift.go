package metadata

//Swift and Objective-C Podfile Lock Metadata
type PodFileLockMetadata struct {
	Pods          []interface{}     `yaml:"PODS"`
	Dependencies  []string          `yaml:"DEPENDENCIES"`
	SpecChecksums map[string]string `yaml:"SPEC CHECKSUMS"`
	Cocoapods     string            `yaml:"COCOAPODS"`
}

type PodFileLockMetadataCheckSums struct {
	Checksums string `json:"checksum"`
}
