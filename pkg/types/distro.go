package types

type Distro struct {
	PrettyName              string   `json:"prettyName,omitempty"`
	Name                    string   `json:"name,omitempty"`
	ID                      string   `json:"id,omitempty"`
	IDLike                  []string `json:"idLike,omitempty"`
	Version                 string   `json:"version,omitempty"`
	VersionID               string   `json:"versionID,omitempty"`
	DistributionID          string   `json:"distributionID,omitempty"`
	DistributionDescription string   `json:"distributionDescription,omitempty"`
	DistributionCodename    string   `json:"distributionCodename,omitempty"`
	HomeURL                 string   `json:"homeURL,omitempty"`
	SupportURL              string   `json:"supportURL,omitempty"`
	BugReportURL            string   `json:"bugReportURL,omitempty"`
	PrivacyPolicyURL        string   `json:"privacyPolicyURL,omitempty"`
}
