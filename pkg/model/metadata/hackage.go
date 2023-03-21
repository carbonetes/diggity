package metadata

// StackConfig stack.yaml metadata containing extra-deps
type StackConfig struct {
	ExtraDeps []string `yaml:"extra-deps"`
}

// StackLockConfig stack.yaml.lock metadata containing packages
type StackLockConfig struct {
	Packages  []StackPackages `yaml:"packages"`
	Snapshots []interface{}   `yaml:"snapshots"`
}

// StackPackages stack.yaml.lock packages metadata
type StackPackages struct {
	Original Hackage `yaml:"original"`
}

// Hackage metadata
type Hackage struct {
	Hackage string `yaml:"hackage"`
}

// HackageMetadata haskell packages metadata
type HackageMetadata struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	PkgHash     string `json:"pkgHash,omitempty"`
	Size        string `json:"size,omitempty"`
	Revision    string `json:"revision,omitempty"`
	SnapshotURL string `json:"snapshotURL,omitempty"`
}
