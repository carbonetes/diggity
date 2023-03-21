package metadata

// RPMMetadata rpm metadata
type RPMMetadata struct {
	Release         string `json:"release,omitempty"`
	Architecture    string `json:"architecture,omitempty"`
	SourceRpm       string `json:"sourceRpm,omitempty"`
	License         string `json:"license,omitempty"`
	Size            int    `json:"size,omitempty"`
	Name            string `json:"name,omitempty"`
	PGP             string `json:"pgp,omitempty"`
	ModularityLabel string `json:"modularityLabel,omitempty"`
	Summary         string `json:"summary,omitempty"`
	Vendor          string `json:"vendor,omitempty"`
	Version         string `json:"version,omitempty"`
	Epoch           int    `json:"epoch,omitempty"`
	DigestAlgorithm string `json:"digestAlgorithm,omitempty"`
}
