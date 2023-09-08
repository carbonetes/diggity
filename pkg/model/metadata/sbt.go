package metadata

type SbtMetadata struct {
	Vendor  string `json:"vendor"`
	Name    string `json:"name"`
	Version string `json:"version"`
	Config  string `json:"config,omitempty"`
}