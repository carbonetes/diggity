package npm

import (
	"encoding/json"
	"io/ioutil"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/scan/parsers/common"
)

// PackageLockParser handles parsing of package-lock.json files
type PackageLockParser struct {
	base *common.BaseParser
}

// NewPackageLockParser creates a new package-lock.json parser
func NewPackageLockParser(base *common.BaseParser) *PackageLockParser {
	return &PackageLockParser{base: base}
}

// PackageLock represents the structure of a package-lock.json file
type PackageLock struct {
	Name         string                           `json:"name"`
	Version      string                           `json:"version"`
	Dependencies map[string]PackageLockDependency `json:"dependencies"`
}

// PackageLockDependency represents a dependency in package-lock.json
type PackageLockDependency struct {
	Version      string                           `json:"version"`
	Dev          bool                             `json:"dev,omitempty"`
	Dependencies map[string]PackageLockDependency `json:"dependencies,omitempty"`
}

// Parse parses package-lock.json files
func (p *PackageLockParser) Parse(filePath string) ([]model.Package, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var packageLock PackageLock
	if err := json.Unmarshal(content, &packageLock); err != nil {
		return nil, err
	}

	var packages []model.Package

	// Parse all dependencies recursively
	packages = append(packages, p.parseDependencies(packageLock.Dependencies, filePath)...)

	return packages, nil
}

// parseDependencies recursively parses dependencies from package-lock.json
func (p *PackageLockParser) parseDependencies(dependencies map[string]PackageLockDependency, filePath string) []model.Package {
	var packages []model.Package

	for name, dep := range dependencies {
		metadata := map[string]interface{}{
			"manager": "npm",
			"locked":  true,
		}

		cleanVersion := p.cleanNPMVersion(dep.Version)
		pkg := p.base.CreatePackage(name, cleanVersion, filePath, dep.Dev, metadata)
		packages = append(packages, pkg)

		// Recursively parse nested dependencies
		if dep.Dependencies != nil {
			nestedPackages := p.parseDependencies(dep.Dependencies, filePath)
			packages = append(packages, nestedPackages...)
		}
	}

	return packages
}

// cleanNPMVersion cleans and normalizes npm version strings
func (p *PackageLockParser) cleanNPMVersion(version string) string {
	prefixes := []string{"^", "~", ">=", "<=", ">", "<", "="}
	return p.base.GetVersionCleaner().CleanVersion(version, prefixes)
}
