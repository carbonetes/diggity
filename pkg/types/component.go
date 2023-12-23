package types

import "github.com/google/uuid"

type Component struct {
	ID              string        `json:"id"`
	Name            string        `json:"name"`
	Version         string        `json:"version"`
	Type            string        `json:"type"`
	PURL            string        `json:"purl,omitempty"`
	Description     string        `json:"description,omitempty"`
	Origin          string        `json:"origin,omitempty"`
	Licenses        []string      `json:"licenses,omitempty"`
	CPEs            []string      `json:"cpes,omitempty"`
	Metadata        interface{}   `json:"metadata,omitempty"`
	Vulnerabilities []interface{} `json:"vulnerabilities,omitempty"`
}

func NewComponent(name, version, category, origin, desc string, metadata ...interface{}) Component {
	return Component{
		ID:          uuid.New().String(),
		Name:        name,
		Version:     version,
		Type:        category,
		Description: desc,
		Origin:      origin,
		PURL:        makePURL(category, name, version),
		Metadata:    metadata,
	}
}

func makePURL(category, name, version string) string {
	return "pkg:" + category + "/" + name + "@" + version
}
