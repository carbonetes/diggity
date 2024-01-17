package types

type OSRelease struct {
	File    string                 `json:"file,omitempty"`
	Release map[string]interface{} `json:"release,omitempty"`
}

// Linux operating system information from /etc/os-release based on the https://www.freedesktop.org/software/systemd/man/latest/os-release.html
type Release struct {
	PrettyName              string   `json:"prettyName,omitempty"`
	Name                    string   `json:"name,omitempty"`
	ID                      string   `json:"id,omitempty"`
	IDLike                  []string `json:"idLike,omitempty"`
	Version                 string   `json:"version,omitempty"`
	VersionID               string   `json:"versionID,omitempty"`
	VersionCodename         string   `json:"versionCodename,omitempty"`
	BuildID                 string   `json:"buildID,omitempty"`
	ImageID                 string   `json:"imageID,omitempty"`
	ImageVersion            string   `json:"imageVersion,omitempty"`
	Variant                 string   `json:"variant,omitempty"`
	VariantID               string   `json:"variantID,omitempty"`
	DistributionID          string   `json:"distributionID,omitempty"`
	DistributionDescription string   `json:"distributionDescription,omitempty"`
	DistributionCodename    string   `json:"distributionCodename,omitempty"`
	HomeURL                 string   `json:"homeURL,omitempty"`
	DocumentationURL        string   `json:"documentationURL,omitempty"`
	SupportURL              string   `json:"supportURL,omitempty"`
	BugReportURL            string   `json:"bugReportURL,omitempty"`
	PrivacyPolicyURL        string   `json:"privacyPolicyURL,omitempty"`
	CPEName                 string   `json:"cpeName,omitempty"`
	SupportEndDate          string   `json:"supportEndDate,omitempty"`
}
