package model

import "github.com/docker/docker/api/types"

// NewRegistryAuth returns a new types.AuthConfig struct with the values set from the given model.Arguments struct.
func NewRegistryAuth(arguments *Arguments) *types.AuthConfig {
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