package npm

import (
	"encoding/json"
	"io/ioutil"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/scan/parsers/common"
)

// PackageJsonParser handles parsing of package.json files
type PackageJsonParser struct {
	base *common.BaseParser
}

// NewPackageJsonParser creates a new package.json parser
func NewPackageJsonParser(base *common.BaseParser) *PackageJsonParser {
	return &PackageJsonParser{base: base}
}

// PackageJson represents the structure of a package.json file
type PackageJson struct {
	Name            string            `json:"name"`
	Version         string            `json:"version"`
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

// Parse parses package.json files
func (p *PackageJsonParser) Parse(filePath string) ([]model.Package, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var packageJson PackageJson
	if err := json.Unmarshal(content, &packageJson); err != nil {
		return nil, err
	}

	var packages []model.Package

	// Parse production dependencies
	for name, version := range packageJson.Dependencies {
		metadata := map[string]interface{}{
			"manager": "npm",
		}

		cleanVersion := p.cleanNPMVersion(version)
		pkg := p.base.CreatePackage(name, cleanVersion, filePath, false, metadata)
		packages = append(packages, pkg)
	}

	// Parse development dependencies
	for name, version := range packageJson.DevDependencies {
		metadata := map[string]interface{}{
			"manager": "npm",
		}

		cleanVersion := p.cleanNPMVersion(version)
		pkg := p.base.CreatePackage(name, cleanVersion, filePath, true, metadata)
		packages = append(packages, pkg)
	}

	return packages, nil
}

// cleanNPMVersion cleans and normalizes npm version strings
func (p *PackageJsonParser) cleanNPMVersion(version string) string {
	prefixes := []string{"^", "~", ">=", "<=", ">", "<", "="}
	return p.base.GetVersionCleaner().CleanVersion(version, prefixes)
}
