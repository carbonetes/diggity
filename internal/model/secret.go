package model

// Secret model
type Secret struct {
	ContentRegexName string `json:"contentRegexName"`
	FileName         string `json:"fileName"`
	FilePath         string `json:"filePath"`
	LineNumber       string `json:"lineNumber"`
}

// SecretConfig model
type SecretConfig struct {
	Disabled    bool      `yaml:"disabled" json:"disabled"`
	SecretRegex string    `yaml:"secret-regex" json:"secretRegex"`
	Excludes    *[]string `yaml:"excludes-filenames" json:"excludesFilenames"`
	MaxFileSize int64     `yaml:"max-file-size" json:"maxFileSize"`
	Extensions  *[]string `yaml:"extensions" json:"extensions,omitempty"`
}

// SecretResults the final result that will be displayed
type SecretResults struct {
	Configuration SecretConfig `json:"applied-configuration"`
	Secrets       []Secret     `json:"secrets"`
}

// DefaultSecretExtensions contains a list of common file extensions containing secrets.
// Additional Reference: https://blog.gitguardian.com/top-10-file-extensions/
var DefaultSecretExtensions = []string{"env", "h", "so", "pem", "properties", "xml", "yml", "yaml", "json", "py", "js", "ts", "PHP"}
