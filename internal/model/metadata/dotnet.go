package metadata

// DotnetDeps - .NET Dependencies
type DotnetDeps struct {
	Libraries map[string]DotnetLibrary `json:"libraries"`
}

// DotnetLibrary - .NET libraries
type DotnetLibrary struct {
	Type     string `json:"type"`
	Path     string `json:"path"`
	Sha512   string `json:"sha512"`
	HashPath string `json:"hashPath"`
}
