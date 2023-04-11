package model

type Relationship struct {
	Parent string `json:"parent"`
	Child  string `json:"child"`
	Type   string `json:"type"`
}

type Ownership struct {
	Main string
	Sub  string
	Path string
}
