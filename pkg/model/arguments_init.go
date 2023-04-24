package model

func NewArguments() *Arguments {
	return &Arguments{
		Image:               new(string),
		Output:              new(Output),
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