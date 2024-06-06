package types

// Secret model
type Secret struct {
	Match       string `json:"match"`
	Description string `json:"description"`
	File        string `json:"file"`
	Layer       string `json:"layer"`
	Content     string `json:"content"`
	Line        int    `json:"line"`
}
