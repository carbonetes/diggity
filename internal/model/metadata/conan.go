package metadata

// ConanMetadata conan metadata
type ConanMetadata struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// ConanLockMetadata conan.lock metadata
type ConanLockMetadata struct {
	GraphLock GraphLock `json:"graph_lock"`
	Version   string    `json:"version"`
}

// GraphLock conan.lock metadata containing nodes
type GraphLock struct {
	Nodes            map[string]ConanLockNode `json:"nodes"`
	RevisionsEnabled bool                     `json:"revisions_enabled"`
}

// ConanLockNode conan.lock packages metadata
type ConanLockNode struct {
	Ref            string      `json:"ref"`
	Path           string      `json:"path,omitempty"`
	Context        string      `json:"context,omitempty"`
	Requires       []string    `json:"requires,omitempty"`
	PackageID      string      `json:"package_id,omitempty"`
	Prev           string      `json:"prev,omitempty"`
	BuildRequires  string      `json:"build_requires,omitempty"`
	PythonRequires string      `json:"py_requires,omitempty"`
	Options        interface{} `json:"options,omitempty"`
}
