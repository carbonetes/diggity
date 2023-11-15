package types

// Secret model
type Secret struct {
	Match   string `json:"match"`
	File    string `json:"file"`
	Content string `json:"content"`
	Line    int    `json:"line"`
}

