package output

// Reference: https://docs.github.com/en/rest/dependency-graph/dependency-submission?apiVersion=2022-11-28
// Model Basis: https://gist.github.com/reiddraper/fdab2883db0f372c146d1a750fc1c43f

// DependencySnapshot Model
type DependencySnapshot struct {
	Version   int                    `json:"version"`
	Job       Job                    `json:"job,omitempty"`
	Sha       string                 `json:"sha,omitempty"`
	Ref       string                 `json:"ref,omitempty"`
	Detector  Detector               `json:"detector,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	Manifests PackageManifests       `json:"manifests,omitempty"`
	Scanned   string                 `json:"scanned,omitempty"`
}

// Job Model
type Job struct {
	Name    string `json:"correlator,omitempty"`
	ID      string `json:"id,omitempty"`
	HTMLURL string `json:"html_url,omitempty"`
}

// Detector Model
type Detector struct {
	Name    string `json:"name,omitempty"`
	URL     string `json:"url,omitempty"`
	Version string `json:"version,omitempty"`
}

// PackageManifests: manifests metadata from sbom packages
type PackageManifests map[string]PackageManifest

// PackageManifest: A collection of related dependencies, either declared in a file,
// or representing a logical group of dependencies.
type PackageManifest struct {
	Name     string                    `json:"name"`
	File     FileInfo                  `json:"file"`
	Metadata map[string]interface{}    `json:"metadata,omitempty"`
	Resolved map[string]DependencyNode `json:"resolved,omitempty"`
}

//FileInfo PackageManifest File Model
type FileInfo struct {
	SourceLocation string `json:"source_location,omitempty"`
}

// DependencyNode DependencyGraph Metadata
type DependencyNode struct {
	PURL         string                 `json:"package_url,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
	Relationship string                 `json:"relationship,omitempty"`
	Scope        string                 `json:"scope,omitempty"`
	Dependencies []string               `json:"dependencies,omitempty"`
}
