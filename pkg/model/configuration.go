package model

// Configuration YAML file config
type Configuration struct {
	SecretConfig       SecretConfig      `yaml:"secret-config"`
	EnabledParsers     []string          `yaml:"enabled-parsers"`
	DisableFileListing bool              `yaml:"disable-file-listing"`
	DisablePullTimeout bool              `yaml:"disable-pull-timeout"`
	Quiet              bool              `yaml:"quiet"`
	OutputFile         string            `yaml:"output-file"`
	Output             *[]string         `yaml:"output"`
	Registry           Registry          `yaml:"registry"`
	AttestationConfig  AttestationConfig `yaml:"attestation"`
}

// Registry config
type Registry struct {
	URI      string `yaml:"uri"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Token    string `yaml:"token"`
}
