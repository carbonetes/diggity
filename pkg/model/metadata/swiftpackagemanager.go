package metadata

type SwiftPackageManagerMetadata struct {
	Object  Object `json:"object,omitempty"`
	Pins    []Pin  `json:"pins,omitempty"`
	Version int    `json:"version,omitempty"`
}

type Object struct {
	Pins []Pin `json:"pins,omitempty"`
}

type Pin struct {
	Identity      string `json:"identity,omitempty"`
	Name          string `json:"package,omitempty"`
	Kind          string `json:"kind,omitempty"`
	RepositoryURL string `json:"repositoryURL,omitempty"`
	Location      string `json:"location,omitempty"`
	State         State  `json:"state,omitempty"`
}

type State struct {
	Revision string `json:"revision,omitempty"`
	Version  string `json:"version,omitempty"`
}
