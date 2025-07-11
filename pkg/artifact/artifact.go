package artifact

import (
	"time"

	"github.com/carbonetes/diggity/pkg/model"
)

// Artifact represents a native structure for scanning packages with their licenses,
// potential secrets, and license files, plus dependency mapping
type Artifact struct {
	ID           string            `json:"id"`            // Unique identifier for artifact
	Package      *model.Package    `json:"package"`       // Core package information
	Licenses     []model.License   `json:"licenses"`      // Associated licenses
	Secrets      []model.Secret    `json:"secrets"`       // Detected secrets
	LicenseFiles []string          `json:"license_files"` // Paths to license files
	Dependencies []Dependency      `json:"dependencies"`  // Enhanced dependency structure
	Path         string            `json:"path"`          // File system path
	Type         string            `json:"type"`          // Package manager type (npm, cargo, etc.)
	ScanTime     time.Time         `json:"scan_time"`     // When this artifact was scanned
	Metadata     map[string]string `json:"metadata"`      // Additional metadata
}

// NewArtifact creates a new artifact with the given ID and package information
func NewArtifact(id string, pkg *model.Package) *Artifact {
	return &Artifact{
		ID:           id,
		Package:      pkg,
		Licenses:     make([]model.License, 0),
		Secrets:      make([]model.Secret, 0),
		LicenseFiles: make([]string, 0),
		Dependencies: make([]Dependency, 0),
		ScanTime:     time.Now(),
		Metadata:     make(map[string]string),
	}
}

// AddLicense adds a license to the artifact
func (a *Artifact) AddLicense(license model.License) {
	a.Licenses = append(a.Licenses, license)
}

// AddSecret adds a secret to the artifact
func (a *Artifact) AddSecret(secret model.Secret) {
	a.Secrets = append(a.Secrets, secret)
}

// AddLicenseFile adds a license file path to the artifact
func (a *Artifact) AddLicenseFile(filePath string) {
	a.LicenseFiles = append(a.LicenseFiles, filePath)
}

// AddDependency adds a dependency to the artifact
func (a *Artifact) AddDependency(dep Dependency) {
	a.Dependencies = append(a.Dependencies, dep)
}

// HasSecrets returns true if the artifact has any detected secrets
func (a *Artifact) HasSecrets() bool {
	return len(a.Secrets) > 0
}

// HasLicenses returns true if the artifact has any associated licenses
func (a *Artifact) HasLicenses() bool {
	return len(a.Licenses) > 0
}

// HasLicenseFiles returns true if the artifact has any license files
func (a *Artifact) HasLicenseFiles() bool {
	return len(a.LicenseFiles) > 0
}

// HasDependencies returns true if the artifact has any dependencies
func (a *Artifact) HasDependencies() bool {
	return len(a.Dependencies) > 0
}

// GetDirectDependencies returns only direct dependencies
func (a *Artifact) GetDirectDependencies() []Dependency {
	var direct []Dependency
	for _, dep := range a.Dependencies {
		if dep.Type == DependencyTypeDirect {
			direct = append(direct, dep)
		}
	}
	return direct
}

// GetDevDependencies returns only development dependencies
func (a *Artifact) GetDevDependencies() []Dependency {
	var dev []Dependency
	for _, dep := range a.Dependencies {
		if dep.Type == DependencyTypeDev {
			dev = append(dev, dep)
		}
	}
	return dev
}

// GetCriticalSecrets returns secrets with critical or high severity
func (a *Artifact) GetCriticalSecrets() []model.Secret {
	var critical []model.Secret
	for _, secret := range a.Secrets {
		if secret.Severity == model.SeverityCritical || secret.Severity == model.SeverityHigh {
			critical = append(critical, secret)
		}
	}
	return critical
}
