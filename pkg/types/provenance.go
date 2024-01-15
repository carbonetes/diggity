package types

type SLSA struct {
	Provenance map[string]interface{} `json:"provenance,omitempty"`
}
