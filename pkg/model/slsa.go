package model

// SLSA - SLSA metadata
type SLSA struct {
	Provenance map[string]interface{} `json:"provenance,omitempty"`
}
