package model

// File - OS Files
type File struct {
	Path        string      `json:"path"`
	OwnerUID    string      `json:"ownerUid,omitempty"`
	OwnerGID    string      `json:"ownerGid,omitempty"`
	Permissions string      `json:"permissions,omitempty"`
	Digest      interface{} `json:"digest,omitempty"`
}
