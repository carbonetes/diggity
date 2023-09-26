package model

// Output type
type Output string

const (
	// JSON Output Type
	JSON Output = "json"
	// Table Output Type (Default)
	Table = "table"
	// CycloneDXXML Output Type
	CycloneDXXML = "cdx-xml"
	// CycloneDXJSON Output Type
	CycloneDXJSON = "cdx-json"
	// SPDXJSON Output Type
	SPDXJSON = "spdx-json"
	// SPDXTagValue Output Type
	SPDXTagValue = "spdx-tag"
	// SPDXYML Output Type
	SPDXYML = "spdx-yml"
	// GithubJSON Output Type
	GithubJSON = "snapshot-json"
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
