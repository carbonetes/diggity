package provider

import (
    "github.com/carbonetes/diggity/pkg/model"
    "github.com/docker/docker/api/types"
)

// NewArguments returns a new model.Arguments struct with all fields initialized to their zero values.
func NewArguments() *model.Arguments {
    return &model.Arguments{
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

// NewRegistryAuth returns a new types.AuthConfig struct with the values set from the given model.Arguments struct.
func NewRegistryAuth(arguments *model.Arguments) *types.AuthConfig {
    return &types.AuthConfig{
        Username:      *setValue(arguments.RegistryUsername),
        Password:      *setValue(arguments.RegistryPassword),
        RegistryToken: *setValue(arguments.RegistryToken),
        ServerAddress: *setValue(arguments.RegistryURI),
    }
}

// setValue returns the value of the given string pointer if it is not nil, otherwise it returns nil.
func setValue(value *string) *string {
    if value != nil {
        return value
    }
    return nil
}
