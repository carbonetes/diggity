package ci

type TokenCheckResponse struct {
	Expired     bool `json:"expired"`
	Permissions []struct {
		Label       string   `json:"label"`
		Permissions []string `json:"permissions"`
	} `json:"permissions"`
	PersonalAccessTokenId string `json:"personalAccessTokenId"`
	// Code string `json:"code"`
}

type ApplicationErrorResponse struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Details []string `json:"details"`
}

type PluginRepo struct {
	PersonalAccessTokenID string `gorm:"column:personal_access_token_id;" json:",omitempty"`
	RepoName              string `gorm:"column:repo_name;" json:",omitempty"`
	PluginType            string `gorm:"foreignkey:plugin_type" json:",omitempty"`
}
