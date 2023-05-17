package model

// File - OS Files
type File struct {
	Path        string      `json:"path"`
	OwnerUID    string      `json:"ownerUid,omitempty"`
	OwnerGID    string      `json:"ownerGid,omitempty"`
	Permissions string      `json:"permissions,omitempty"`
	Digest      interface{} `json:"digest,omitempty"`
}

type Conffile struct {
	Path         string `json:"path"`
	Digest       Digest `json:"digest"`
	IsConfigFile bool   `json:"isConfigFile"`
}

type Digest struct {
	Algorithm string `json:"algorithm"`
	Value     string `json:"value"`
}
