package model

// Distro docker image distro
type Distro struct {
	PrettyName         string   `json:"prettyName,omitempty"`
	Name               string   `json:"name,omitempty"`
	ID                 string   `json:"id,omitempty"`
	IDLike             []string `json:"idLike,omitempty"`
	Version            string   `json:"version,omitempty"`
	VersionID          string   `json:"versionID,omitempty"`
	DistribID          string   `json:"distribID,omitempty"`
	DistribDescription string   `json:"distribDescription,omitempty"`
	DistribCodename    string   `json:"versionCodename,omitempty"`
	HomeURL            string   `json:"homeURL,omitempty"`
	SupportURL         string   `json:"supportURL,omitempty"`
	BugReportURL       string   `json:"bugReportURL,omitempty"`
	PrivacyPolicyURL   string   `json:"privacyPolicyURL,omitempty"`
}
