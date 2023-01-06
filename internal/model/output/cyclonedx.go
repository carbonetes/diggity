package output

import "encoding/xml"

//CycloneFormat - CycloneDX Output Model
type CycloneFormat struct {
	// XML specific fields
	XMLName      xml.Name     `json:"-" xml:"bom"`
	XMLNS        string       `json:"-" xml:"xmlns,attr"`
	SerialNumber string       `json:"serialNumber,omitempty" xml:"serialNumber,attr,omitempty"`
	Metadata     *Metadata    `json:"metadata,omitempty" xml:"metadata,omitempty"`
	Components   *[]Component `json:"components,omitempty" xml:"components>component,omitempty"`
}

// Metadata - cyclone format metadata
type Metadata struct {
	Timestamp  string      `json:"timestamp,omitempty" xml:"timestamp,omitempty"`
	Tools      *[]Tool     `json:"tools,omitempty" xml:"tools>tool,omitempty"`
	Component  *Component  `json:"component,omitempty" xml:"component,omitempty"`
	Licenses   *[]License  `json:"licenses,omitempty" xml:"licenses>license,omitempty"`
	Properties *[]Property `json:"properties,omitempty" xml:"properties>property,omitempty"`
}

// Tool - metadata tool
type Tool struct {
	Vendor  string `json:"vendor,omitempty" xml:"vendor,omitempty"`
	Name    string `json:"name" xml:"name"`
	Version string `json:"version,omitempty" xml:"version,omitempty"`
}

//ComponentLibrary - component library type
type ComponentLibrary string

// OperatingSystem - operating system type
type OperatingSystem string

// Component - CycloneFormat component
type Component struct {
	BOMRef             string               `json:"bom-ref,omitempty" xml:"bom-ref,attr,omitempty"`
	MIMEType           string               `json:"mime-type,omitempty" xml:"mime-type,attr,omitempty"`
	Type               ComponentLibrary     `json:"type" xml:"type,attr"`
	Author             string               `json:"author,omitempty" xml:"author,omitempty"`
	Publisher          string               `json:"publisher,omitempty" xml:"publisher,omitempty"`
	Group              string               `json:"group,omitempty" xml:"group,omitempty"`
	Name               string               `json:"name" xml:"name"`
	Version            string               `json:"version,omitempty" xml:"version,omitempty"`
	Description        string               `json:"description,omitempty" xml:"description,omitempty"`
	Licenses           *[]License           `json:"licenses,omitempty" xml:"licenses>license,omitempty"`
	Copyright          string               `json:"copyright,omitempty" xml:"copyright,omitempty"`
	CPE                string               `json:"cpe,omitempty" xml:"cpe,omitempty"`
	PackageURL         string               `json:"purl,omitempty" xml:"purl,omitempty"`
	ExternalReferences *[]ExternalReference `json:"externalReferences,omitempty" xml:"externalReferences>reference,omitempty"`
	Modified           *bool                `json:"modified,omitempty" xml:"modified,omitempty"`
	Properties         *[]Property          `json:"properties,omitempty" xml:"properties>property,omitempty"`
	Components         *[]Component         `json:"components,omitempty" xml:"components>component,omitempty"`
}

// License - Component Licenses
type License struct {
	ID   string `json:"id,omitempty" xml:"id,omitempty"`
	Name string `json:"name,omitempty" xml:"name,omitempty"`
	URL  string `json:"url,omitempty" xml:"url,omitempty"`
}

//Property - Component Properties
type Property struct {
	Name  string `json:"name" xml:"name,attr"`
	Value string `json:"value" xml:",chardata"`
}

//ExternalReference - Component External References
type ExternalReference struct {
	URL     string                `json:"url" xml:"url"`
	Comment string                `json:"comment,omitempty" xml:"comment,omitempty"`
	Type    ExternalReferenceType `json:"type" xml:"type,attr"`
}

//ExternalReferenceType - External Reference Type
type ExternalReferenceType string
