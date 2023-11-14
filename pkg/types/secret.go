package types

// Secret model
type Secret struct {
	Match   string `json:"match"`
	File    string `json:"file"`
	Content string `json:"content"`
	Line    int    `json:"line"`
}

// SecretResults the final result that will be displayed
type SecretResult struct {
	Parameters SecretParameters `json:"parameters"`
	Secrets    []Secret         `json:"secrets"`
}
