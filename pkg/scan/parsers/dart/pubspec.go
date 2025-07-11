package dart

import (
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/scan/parsers/common"
	"gopkg.in/yaml.v3"
)

// PubspecParser handles parsing of pubspec.yaml and pubspec.lock files
type PubspecParser struct {
	base *common.BaseParser
}

// NewPubspecParser creates a new pubspec parser
func NewPubspecParser(base *common.BaseParser) *PubspecParser {
	return &PubspecParser{base: base}
}

// PubspecYaml represents the structure of pubspec.yaml
type PubspecYaml struct {
	Name            string                 `yaml:"name"`
	Version         string                 `yaml:"version"`
	Dependencies    map[string]interface{} `yaml:"dependencies"`
	DevDependencies map[string]interface{} `yaml:"dev_dependencies"`
}

// PubspecLock represents the structure of pubspec.lock
type PubspecLock struct {
	Packages map[string]PubPackage `yaml:"packages"`
}

// PubPackage represents a package in pubspec.lock
type PubPackage struct {
	Dependency string `yaml:"dependency"`
	Source     string `yaml:"source"`
	Version    string `yaml:"version"`
}

// Parse parses pubspec.yaml files
func (p *PubspecParser) Parse(filePath string) ([]model.Package, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var pubspec PubspecYaml
	if err := yaml.Unmarshal(content, &pubspec); err != nil {
		return nil, err
	}

	var packages []model.Package

	// Parse dependencies
	for name, versionInfo := range pubspec.Dependencies {
		if name == "flutter" {
			// Skip Flutter SDK dependency
			continue
		}

		version := p.extractVersion(versionInfo)
		if version == "" {
			continue
		}

		metadata := map[string]interface{}{
			"manager":    "pub",
			"dependency": "direct",
		}

		cleanVersion := p.cleanDartVersion(version)
		pkg := p.base.CreatePackage(name, cleanVersion, filePath, false, metadata)
		packages = append(packages, pkg)
	}

	// Parse dev dependencies
	for name, versionInfo := range pubspec.DevDependencies {
		if name == "flutter_test" || name == "flutter_lints" {
			// Skip Flutter SDK dev dependencies
			continue
		}

		version := p.extractVersion(versionInfo)
		if version == "" {
			continue
		}

		metadata := map[string]interface{}{
			"manager":    "pub",
			"dependency": "dev",
		}

		cleanVersion := p.cleanDartVersion(version)
		pkg := p.base.CreatePackage(name, cleanVersion, filePath, false, metadata)
		packages = append(packages, pkg)
	}

	return packages, nil
}

// ParseLock parses pubspec.lock files
func (p *PubspecParser) ParseLock(filePath string) ([]model.Package, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var lockFile PubspecLock
	if err := yaml.Unmarshal(content, &lockFile); err != nil {
		return nil, err
	}

	var packages []model.Package

	for name, pkg := range lockFile.Packages {
		if pkg.Source != "hosted" {
			// Only process packages from pub.dev
			continue
		}

		metadata := map[string]interface{}{
			"manager":    "pub",
			"dependency": pkg.Dependency,
			"source":     pkg.Source,
		}

		cleanVersion := p.cleanDartVersion(pkg.Version)
		modelPkg := p.base.CreatePackage(name, cleanVersion, filePath, true, metadata)
		packages = append(packages, modelPkg)
	}

	return packages, nil
}

// extractVersion extracts version string from various dependency formats
func (p *PubspecParser) extractVersion(versionInfo interface{}) string {
	switch v := versionInfo.(type) {
	case string:
		return v
	case map[string]interface{}:
		if version, ok := v["version"].(string); ok {
			return version
		}
		// Skip git, path, and SDK dependencies
		return ""
	default:
		return ""
	}
}

// cleanDartVersion cleans and normalizes Dart version strings
func (p *PubspecParser) cleanDartVersion(version string) string {
	version = strings.TrimSpace(version)

	// Remove common Dart version prefixes
	prefixes := []string{"^", "~", ">=", "<=", ">", "<", "="}
	cleanVersion := p.base.GetVersionCleaner().CleanVersion(version, prefixes)

	// Handle version ranges like ">=1.0.0 <2.0.0"
	if strings.Contains(cleanVersion, " ") {
		// Take the first version for ranges
		parts := strings.Fields(cleanVersion)
		if len(parts) > 0 {
			// Remove any remaining operators from the first part
			firstVersion := parts[0]
			re := regexp.MustCompile(`^[^\d]*(.+)`)
			if matches := re.FindStringSubmatch(firstVersion); len(matches) > 1 {
				cleanVersion = matches[1]
			}
		}
	}

	return cleanVersion
}
