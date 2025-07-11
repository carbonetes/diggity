package python

import (
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/scan/parsers/common"
)

// PyprojectParser handles parsing of pyproject.toml files
type PyprojectParser struct {
	base *common.BaseParser
}

// NewPyprojectParser creates a new pyproject.toml parser
func NewPyprojectParser(base *common.BaseParser) *PyprojectParser {
	return &PyprojectParser{base: base}
}

// Parse parses pyproject.toml files
func (p *PyprojectParser) Parse(filePath string) ([]model.Package, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var packages []model.Package
	contentStr := string(content)

	// Parse Poetry dependencies
	poetryDeps := p.parsePoetryDependencies(contentStr, filePath)
	packages = append(packages, poetryDeps...)

	// Parse PEP 621 project dependencies
	projectDeps := p.parseProjectDependencies(contentStr, filePath)
	packages = append(packages, projectDeps...)

	return packages, nil
}

// parsePoetryDependencies parses dependencies from [tool.poetry.dependencies] section
func (p *PyprojectParser) parsePoetryDependencies(content, filePath string) []model.Package {
	var packages []model.Package

	// Look for [tool.poetry.dependencies] section
	sectionPattern := `\[tool\.poetry\.dependencies\](.*?)(?:\[|$)`
	re := regexp.MustCompile(sectionPattern)
	matches := re.FindStringSubmatch(content)

	if len(matches) > 1 {
		deps := matches[1]
		packages = append(packages, p.parseTomlDependencies(deps, filePath, "poetry")...)
	}

	// Also look for dev dependencies
	devSectionPattern := `\[tool\.poetry\.group\.dev\.dependencies\](.*?)(?:\[|$)`
	devRe := regexp.MustCompile(devSectionPattern)
	devMatches := devRe.FindStringSubmatch(content)

	if len(devMatches) > 1 {
		deps := devMatches[1]
		packages = append(packages, p.parseTomlDevDependencies(deps, filePath, "poetry")...)
	}

	return packages
}

// parseProjectDependencies parses dependencies from [project] section (PEP 621)
func (p *PyprojectParser) parseProjectDependencies(content, filePath string) []model.Package {
	var packages []model.Package

	// Look for dependencies array in [project] section
	pattern := `\[project\].*?dependencies\s*=\s*\[(.*?)\]`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(content)

	if len(matches) > 1 {
		deps := matches[1]
		// Split by comma and parse each dependency
		for _, dep := range strings.Split(deps, ",") {
			dep = strings.TrimSpace(dep)
			dep = strings.Trim(dep, "\"'")
			if dep == "" {
				continue
			}

			name, version := p.parseDepString(dep)
			if name != "" {
				metadata := map[string]interface{}{
					"manager": "pep621",
				}

				pkg := p.base.CreatePackage(name, version, filePath, false, metadata)
				packages = append(packages, pkg)
			}
		}
	}

	return packages
}

// parseTomlDependencies parses TOML format dependencies
func (p *PyprojectParser) parseTomlDependencies(deps, filePath, manager string) []model.Package {
	var packages []model.Package
	lines := strings.Split(deps, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.Contains(line, "=") {
			name, version := p.parseTomlLine(line)
			if name != "" {
				metadata := map[string]interface{}{
					"manager": manager,
				}

				pkg := p.base.CreatePackage(name, version, filePath, false, metadata)
				packages = append(packages, pkg)
			}
		}
	}

	return packages
}

// parseTomlDevDependencies parses TOML format dev dependencies
func (p *PyprojectParser) parseTomlDevDependencies(deps, filePath, manager string) []model.Package {
	var packages []model.Package
	lines := strings.Split(deps, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.Contains(line, "=") {
			name, version := p.parseTomlLine(line)
			if name != "" {
				metadata := map[string]interface{}{
					"manager": manager,
				}

				pkg := p.base.CreatePackage(name, version, filePath, true, metadata)
				packages = append(packages, pkg)
			}
		}
	}

	return packages
}

// parseTomlLine parses a TOML line to extract name and version
func (p *PyprojectParser) parseTomlLine(line string) (string, string) {
	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return "", ""
	}

	name := strings.TrimSpace(parts[0])
	version := strings.TrimSpace(parts[1])

	// Remove quotes
	name = strings.Trim(name, "\"'")
	version = strings.Trim(version, "\"'")

	return name, p.cleanPythonVersion(version)
}

// parseDepString parses a dependency string to extract name and version
func (p *PyprojectParser) parseDepString(dep string) (string, string) {
	// Common patterns for pyproject.toml dependencies
	patterns := []string{
		`^([a-zA-Z0-9\-_.]+)\s*([<>=!~]+)\s*([0-9][^,\s;]*)?`, // package>=1.0.0
		`^([a-zA-Z0-9\-_.]+)\s*==\s*([0-9][^\s,;]*)`,          // package==1.0.0
		`^([a-zA-Z0-9\-_.]+)\s*$`,                             // package (no version)
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(dep)

		if len(matches) >= 2 {
			name := matches[1]
			version := ""
			if len(matches) >= 3 && matches[2] != "" {
				if len(matches) >= 4 && matches[3] != "" {
					version = matches[3]
				} else {
					version = matches[2]
				}
				version = p.cleanPythonVersion(version)
			}
			return name, version
		}
	}

	return "", ""
}

// cleanPythonVersion cleans and normalizes python version strings
func (p *PyprojectParser) cleanPythonVersion(version string) string {
	prefixes := []string{"^", "~", ">=", "<=", ">", "<", "==", "!="}
	return p.base.GetVersionCleaner().CleanVersion(version, prefixes)
}
