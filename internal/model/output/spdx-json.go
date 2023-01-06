package output

import "time"

// SpdxJSONDocument Model
type SpdxJSONDocument struct {
	SPDXID            string       `json:"SPDXID"`
	Name              string       `json:"name,omitempty"`
	SpdxVersion       string       `json:"spdxVersion"`
	CreationInfo      CreationInfo `json:"creationInfo"`
	DataLicense       string       `json:"dataLicense"`
	DocumentNamespace string       `json:"documentNamespace"`
	// SpdxJsonPackages Actual Packages
	SpdxJSONPackages []SpdxJSONPackage `json:"packages,omitempty"`
}

// SpdxJSONPackage Model
type SpdxJSONPackage struct {
	SpdxID           string        `json:"SPDXID"`
	Name             string        `json:"name,omitempty"`
	LicenseConcluded string        `json:"licenseConcluded,omitempty"`
	Description      string        `json:"description,omitempty"`
	DownloadLocation string        `json:"downloadLocation,omitempty"`
	ExternalRefs     []ExternalRef `json:"externalRefs,omitempty"`
	FilesAnalyzed    bool          `json:"filesAnalyzed"`
	Homepage         string        `json:"homepage,omitempty"`
	LicenseDeclared  string        `json:"licenseDeclared,omitempty"`
	Originator       string        `json:"originator,omitempty"`
	SourceInfo       string        `json:"sourceInfo,omitempty"`
	VersionInfo      string        `json:"versionInfo,omitempty"`
	Copyright        string        `json:"copyright,omitempty"`
}

// ExternalRef Model
type ExternalRef struct {
	ReferenceCategory string `json:"referenceCategory,omitempty"`
	ReferenceLocator  string `json:"referenceLocator,omitempty"`
	ReferenceType     string `json:"referenceType,omitempty"`
}

// CreationInfo Model
type CreationInfo struct {
	Created            time.Time `json:"created"`
	Creators           []string  `json:"creators"`
	LicenseListVersion string    `json:"licenseListVersion"`
}
