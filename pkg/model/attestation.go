package model

// AttestationConfig model
type AttestationConfig struct {
	Key      string `yaml:"key"`
	Pub      string `yaml:"pub"`
	Password string `yaml:"password"`
}

// AttestationOptions model
type AttestationOptions struct {
	Key        *string
	Pub        *string
	AttestType *string
	Predicate  *string
	Password   *string
	OutputFile *string
	OutputType *string
	BomArgs    *Arguments
	Provenance *string
}
