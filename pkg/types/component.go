package types

type Component struct {
	ID              string        `json:"id"`
	Name            string        `json:"name"`
	Version         string        `json:"version"`
	Type            string        `json:"type"`
	PURL            string        `json:"purl,omitempty"`
	Description     string        `json:"description,omitempty"`
	Origin          string        `json:"origin,omitempty"`
	Licenses        []string      `json:"licenses,omitempty"`
	Metadata        interface{}   `json:"metadata,omitempty"`
	Dependencies    []Dependency  `json:"dependencies,omitempty"`
	Vulnerabilities []interface{} `json:"vulnerabilities,omitempty"`
}

type Dependency struct {
	ParentID string `json:"parent_id,omitempty"`
	ChildID  string `json:"child_id,omitempty"`
}

