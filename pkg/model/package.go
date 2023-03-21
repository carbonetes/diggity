package model

// Package actual package found
type Package struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Type        string      `json:"type"`
	Version     string      `json:"version"`
	Path        string      `json:"path"`
	Locations   []Location  `json:"locations"`
	Description string      `json:"description,omitempty"`
	Licenses    []string    `json:"licenses,omitempty"`
	CPEs        []string    `json:"cpes"`
	PURL        PURL        `json:"purl"`
	Metadata    interface{} `json:"metadata"`
}

// PURL - Package URL
type PURL string
