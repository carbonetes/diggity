package dotnet

import (
	"encoding/xml"
	"io/ioutil"
	"strings"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/scan/parsers/common"
)

// ProjectFileParser handles parsing of .csproj, .vbproj, .fsproj files
type ProjectFileParser struct {
	base *common.BaseParser
}

// NewProjectFileParser creates a new project file parser
func NewProjectFileParser(base *common.BaseParser) *ProjectFileParser {
	return &ProjectFileParser{base: base}
}

// PackageConfigParser handles parsing of packages.config files
type PackageConfigParser struct {
	base *common.BaseParser
}

// NewPackageConfigParser creates a new packages.config parser
func NewPackageConfigParser(base *common.BaseParser) *PackageConfigParser {
	return &PackageConfigParser{base: base}
}

// PackagesLockParser handles parsing of packages.lock.json files
type PackagesLockParser struct {
	base *common.BaseParser
}

// NewPackagesLockParser creates a new packages.lock.json parser
func NewPackagesLockParser(base *common.BaseParser) *PackagesLockParser {
	return &PackagesLockParser{base: base}
}

// Project represents the structure of a .NET project file
type Project struct {
	XMLName    xml.Name    `xml:"Project"`
	ItemGroups []ItemGroup `xml:"ItemGroup"`
}

// ItemGroup represents an ItemGroup in a project file
type ItemGroup struct {
	PackageReferences []PackageReference `xml:"PackageReference"`
	Packages          []Package          `xml:"package"`
}

// PackageReference represents a PackageReference in modern .NET projects
type PackageReference struct {
	Include string `xml:"Include,attr"`
	Version string `xml:"Version,attr"`
}

// Package represents a package in packages.config
type Package struct {
	ID      string `xml:"id,attr"`
	Version string `xml:"version,attr"`
}

// Parse parses .NET project files (.csproj, .vbproj, .fsproj)
func (p *ProjectFileParser) Parse(filePath string) ([]model.Package, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var project Project
	if err := xml.Unmarshal(content, &project); err != nil {
		return nil, err
	}

	var packages []model.Package

	for _, itemGroup := range project.ItemGroups {
		for _, packageRef := range itemGroup.PackageReferences {
			if packageRef.Include == "" {
				continue
			}

			metadata := map[string]interface{}{
				"manager": "nuget",
			}

			cleanVersion := p.cleanDotNetVersion(packageRef.Version)
			pkg := p.base.CreatePackage(packageRef.Include, cleanVersion, filePath, false, metadata)
			packages = append(packages, pkg)
		}
	}

	return packages, nil
}

// Parse parses packages.config files
func (p *PackageConfigParser) Parse(filePath string) ([]model.Package, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var packages struct {
		XMLName  xml.Name  `xml:"packages"`
		Packages []Package `xml:"package"`
	}

	if err := xml.Unmarshal(content, &packages); err != nil {
		return nil, err
	}

	var result []model.Package

	for _, pkg := range packages.Packages {
		if pkg.ID == "" {
			continue
		}

		metadata := map[string]interface{}{
			"manager": "nuget",
		}

		cleanVersion := p.cleanDotNetVersion(pkg.Version)
		modelPkg := p.base.CreatePackage(pkg.ID, cleanVersion, filePath, false, metadata)
		result = append(result, modelPkg)
	}

	return result, nil
}

// Parse parses packages.lock.json files (simplified implementation)
func (p *PackagesLockParser) Parse(filePath string) ([]model.Package, error) {
	// For now, return empty as packages.lock.json is complex and not always present
	return []model.Package{}, nil
}

// cleanDotNetVersion cleans and normalizes .NET version strings
func (p *ProjectFileParser) cleanDotNetVersion(version string) string {
	// Remove common .NET version prefixes
	version = strings.TrimSpace(version)
	prefixes := []string{"[", "]", "(", ")", ">=", "<=", ">", "<", "="}
	return p.base.GetVersionCleaner().CleanVersion(version, prefixes)
}

// cleanDotNetVersion cleans and normalizes .NET version strings
func (p *PackageConfigParser) cleanDotNetVersion(version string) string {
	version = strings.TrimSpace(version)
	prefixes := []string{"[", "]", "(", ")", ">=", "<=", ">", "<", "="}
	return p.base.GetVersionCleaner().CleanVersion(version, prefixes)
}
