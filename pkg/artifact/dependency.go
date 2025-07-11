package artifact

// Dependency represents a package dependency relationship with full context
type Dependency struct {
	ID          string            `json:"id"`          // Unique dependency identifier
	Name        string            `json:"name"`        // Dependency name
	Version     string            `json:"version"`     // Dependency version
	Type        DependencyType    `json:"type"`        // Dependency type (direct, transitive, dev, etc.)
	Source      string            `json:"source"`      // Package manager source
	Constraints string            `json:"constraints"` // Version constraints
	Scope       string            `json:"scope"`       // Dependency scope (runtime, test, build, etc.)
	Optional    bool              `json:"optional"`    // Whether dependency is optional
	Metadata    map[string]string `json:"metadata"`    // Additional dependency metadata
}

// DependencyType represents the type of dependency relationship
type DependencyType string

const (
	DependencyTypeDirect     DependencyType = "direct"     // Direct dependency
	DependencyTypeTransitive DependencyType = "transitive" // Transitive dependency
	DependencyTypeDev        DependencyType = "dev"        // Development dependency
	DependencyTypePeer       DependencyType = "peer"       // Peer dependency
	DependencyTypeOptional   DependencyType = "optional"   // Optional dependency
	DependencyTypeBundled    DependencyType = "bundled"    // Bundled dependency
)

// NewDependency creates a new dependency with the given name and version
func NewDependency(name, version string, depType DependencyType) Dependency {
	return Dependency{
		ID:       name + "@" + version,
		Name:     name,
		Version:  version,
		Type:     depType,
		Optional: false,
		Metadata: make(map[string]string),
	}
}

// IsDirectDependency returns true if this is a direct dependency
func (d *Dependency) IsDirectDependency() bool {
	return d.Type == DependencyTypeDirect
}

// IsDevDependency returns true if this is a development dependency
func (d *Dependency) IsDevDependency() bool {
	return d.Type == DependencyTypeDev
}

// IsTransitiveDependency returns true if this is a transitive dependency
func (d *Dependency) IsTransitiveDependency() bool {
	return d.Type == DependencyTypeTransitive
}

// IsOptionalDependency returns true if this is an optional dependency
func (d *Dependency) IsOptionalDependency() bool {
	return d.Optional || d.Type == DependencyTypeOptional
}

// GetFullName returns the dependency name with version
func (d *Dependency) GetFullName() string {
	if d.Version != "" {
		return d.Name + "@" + d.Version
	}
	return d.Name
}
