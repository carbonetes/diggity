package metadata

// PackageJSON - packages.json model
type PackageJSON struct {
	Version      string                 `json:"version"`
	Latest       []string               `json:"latest"`
	Contributors interface{}            `json:"contributors"`
	License      interface{}            `json:"license"`
	Name         string                 `json:"name"`
	Homepage     string                 `json:"homepage"`
	Description  string                 `json:"description"`
	Dependencies map[string]interface{} `json:"dependencies"`
	Repository   interface{}            `json:"repository"`
	Author       interface{}            `json:"author"`
}

// Contributors - PackageJSON contributors
type Contributors struct {
	Name     string `json:"name" mapstruct:"name"`
	Username string `json:"email" mapstruct:"username"`
	URL      string `json:"url" mapstruct:"url"`
}

// Repository - PackageJSON repository
type Repository struct {
	Type string `json:"type" mapstructure:"type"`
	URL  string `json:"url" mapstructure:"url"`
}

//PackageLock - PackageLock model
type PackageLock struct {
	Requires        bool `json:"requires"`
	LockfileVersion int  `json:"lockfileVersion"`
	Dependencies    map[string]LockDependency
}

// LockDependency - PackageLock dependencies
type LockDependency struct {
	Version   string `json:"version"`
	Resolved  string `json:"resolved"`
	Integrity string `json:"integrity"`
	Requires  map[string]string
}
