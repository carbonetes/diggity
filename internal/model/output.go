package model

// Output type
type Output string

const (
	// JSON Output Type
	JSON Output = "json"
	// Table Output Type (Default)
	Table = "table"
	// CycloneDX Output Type
	CycloneDX = "cyclonedx"
	// CycloneDXJSON Output Type
	CycloneDXJSON = "cyclonedx-json"
	// SPDXJSON Output Type
	SPDXJSON = "spdx-json"
	// SPDXTagValue Output Type
	SPDXTagValue = "spdx-tag-value"
)

var (
	// OutputTypes - All Supported Output Types
	OutputTypes = map[string]string{
		JSON.ToOutput(): JSON.ToOutput(),
		Table:           Table,
		CycloneDX:       CycloneDX,
		CycloneDXJSON:   CycloneDXJSON,
		SPDXJSON:        SPDXJSON,
		SPDXTagValue:    SPDXTagValue,
	}
)

// ToOutput - returns the output type as string
func (c Output) ToOutput() string {
	return string(c)
}
