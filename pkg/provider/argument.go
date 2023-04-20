package provider

import (
	"github.com/carbonetes/diggity/pkg/model"
)

// NewArguments returns a new model.Arguments struct with all fields initialized to their zero values.
func NewArguments() *model.Arguments {
	return &model.Arguments{
		Image:               new(string),
		Output:              new(model.Output),
		DisableFileListing:  new(bool),
		SecretContentRegex:  new(string),
		DisableSecretSearch: new(bool),
		DisablePullTimeout:  new(bool),
		Dir:                 new(string),
		Tar:                 new(string),
		Quiet:               new(bool),
		OutputFile:          new(string),
		ExcludedFilenames:   &[]string{},
		SecretExtensions:    &[]string{},
		EnabledParsers:      &[]string{},
		RegistryURI:         new(string),
		RegistryUsername:    new(string),
		RegistryPassword:    new(string),
		RegistryToken:       new(string),
		Provenance:          new(string),
	}
}
