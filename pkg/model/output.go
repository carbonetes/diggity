package model

type (
	// Output type
	Output string

	// Result - Final SBOM output content
	Result struct {
		Packages      []*Package     `json:"packages"`
		Relationships []Relationship `json:"relationships,omitempty"`
		Secret        *SecretResults `json:"secrets,omitempty"`
		SourceInfo    *SourceInfo    `json:"sourceInfo,omitempty"`
		ImageInfo     *ImageInfo     `json:"imageInfo,omitempty"`
		Distro        *Distro        `json:"distro"`
		SLSA          *SLSA          `json:"slsa,omitempty"`
	}
)

const (
	// JSON Output Type
	JSON Output = "json"
	// Table Output Type (Default)
	Table = "table"
	// CycloneDXXML Output Type
	CycloneDXXML = "cyclonedx-xml"
	// CycloneDXJSON Output Type
	CycloneDXJSON = "cyclonedx-json"
	// SPDXJSON Output Type
	SPDXJSON = "spdx-json"
	// SPDXTagValue Output Type
	SPDXTagValue = "spdx-tag-value"
	// SPDXYML Output Type
	SPDXYML = "spdx-yml"
	// GithubJSON Output Type
	GithubJSON = "github-json"
)

var (
	// OutputTypes - All Supported Output Types
	OutputTypes = map[string]string{
		JSON.ToOutput(): JSON.ToOutput(),
		Table:           Table,
		CycloneDXXML:    CycloneDXXML,
		CycloneDXJSON:   CycloneDXJSON,
		SPDXJSON:        SPDXJSON,
		SPDXTagValue:    SPDXTagValue,
		SPDXYML:         SPDXYML,
		GithubJSON:      GithubJSON,
	}

	// OutputList - List of supported output types
	OutputList = []string{
		JSON.ToOutput(),
		Table,
		CycloneDXXML,
		CycloneDXJSON,
		SPDXJSON,
		SPDXTagValue,
		SPDXYML,
		GithubJSON}
)

// ToOutput - returns the output type as string
func (c Output) ToOutput() string {
	return string(c)
}
