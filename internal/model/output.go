package model

// Output type
type Output string

const (
	// JSON Output Type
	JSON Output = "json"
	// Table Output Type (Default)
	Table = "table"
	// CycloneDX Output Type
	CycloneDXXML = "cyclonedx-xml"
	// CycloneDXJSON Output Type
	CycloneDXJSON = "cyclonedx-json"
	// SPDXJSON Output Type
	SPDXJSON = "spdx-json"
	// SPDXTagValue Output Type
	SPDXTagValue = "spdx-tag-value"
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
		GithubJSON:      GithubJSON,
	}

	// OutputList - List of supported output types
	OutputList = []string{JSON.ToOutput(), Table, CycloneDXXML, CycloneDXJSON, SPDXJSON, SPDXTagValue, GithubJSON}

	// OutputAliases - valid aliases of the output types
	OutputAliases = map[string]string{
		// CycloneDX-XML
		"cyclonedxxml": "cyclonedxxml",
		"cyclonedx":    "cyclonedx",
		"cyclone":      "cyclone",
		// CycloneDX-JSON
		"cyclonedxjson": "cyclonedxjson",
		// SPDX-JSON
		"spdxjson": "spdxjson",
		// SPDX-Tag-Value
		"spdxtagvalue": "spdxtagvalue",
		"spdx":         "spdx",
		"spdxtv":       "spdxtv",
		// Github JSON
		"githubjson": "githubjson",
		"github":     "github",
	}
)

// ToOutput - returns the output type as string
func (c Output) ToOutput() string {
	return string(c)
}
