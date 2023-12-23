package types

const ConfigVersion string = "1.0"

type Config struct {
	Version      string             `json:"version" yaml:"version"`
	MaxFileSize  int64              `json:"max_file_size" yaml:"max_file_size"`
	Registry     RegistryParameters `json:"registry" yaml:"registry"`
	SecretConfig SecretConfig       `json:"secret_config" yaml:"secret_config"`
}

type SecretConfig struct {
	Whitelist Whitelist `json:"whitelist" yaml:"whitelist"`
	Rules     []Rule    `json:"rules" yaml:"rules"`
}

type Whitelist struct {
	Patterns []string `json:"patterns" yaml:"patterns"`
	Keywords []string `json:"keywords" yaml:"keywords"`
}

type Rule struct {
	ID          string   `json:"id" yaml:"id"`
	Description string   `json:"description" yaml:"description"`
	Pattern     string   `json:"pattern" yaml:"pattern"`
	Keywords    []string `json:"keywords" yaml:"keywords"`
}
