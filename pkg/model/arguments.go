package model

// Arguments - CLI Arguments
type Arguments struct {
	Image                *string
	Output               *Output
	Quiet                *bool
	OutputFile           *string
	EnabledParsers       *[]string
	DisableFileListing   *bool
	DisableRelationships *bool
	SecretContentRegex   *string
	DisableSecretSearch  *bool
	SecretMaxFileSize    int64
	RegistryURI          *string
	RegistryUsername     *string
	RegistryPassword     *string
	RegistryToken        *string
	Dir                  *string
	Tar                  *string
	ExcludedFilenames    *[]string
	SecretExtensions     *[]string
	Provenance           *string
}
