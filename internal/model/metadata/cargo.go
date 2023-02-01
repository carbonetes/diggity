package metadata

// CargoMetadata rust metadata
type CargoMetadata struct {
	Name         string   `json:"name,omitempty"`
	Version      string   `json:"version,omitempty"`
	Source       string   `json:"source,omitempty"`
	Checksum     string   `json:"checksum,omitempty"`
	Dependencies []string `json:"dependencies,omitempty"`
}
