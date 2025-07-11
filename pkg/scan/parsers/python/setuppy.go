package python

import (
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/scan/parsers/common"
)

// SetupPyParser handles parsing of setup.py files
type SetupPyParser struct {
	base *common.BaseParser
}

// NewSetupPyParser creates a new setup.py parser
func NewSetupPyParser(base *common.BaseParser) *SetupPyParser {
	return &SetupPyParser{base: base}
}

// Parse parses setup.py files
func (s *SetupPyParser) Parse(filePath string) ([]model.Package, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var packages []model.Package
	contentStr := string(content)

	// Look for install_requires
	patterns := []string{
		`install_requires\s*=\s*\[(.*?)\]`,
		`requires\s*=\s*\[(.*?)\]`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(contentStr)
		if len(matches) > 1 {
			deps := matches[1]
			depPackages := s.parseDependencies(deps, filePath)
			packages = append(packages, depPackages...)
		}
	}

	return packages, nil
}

// parseDependencies parses dependencies from setup.py content
func (s *SetupPyParser) parseDependencies(deps, filePath string) []model.Package {
	var packages []model.Package

	// Split by comma and parse each dependency
	for _, dep := range strings.Split(deps, ",") {
		dep = strings.TrimSpace(dep)
		dep = strings.Trim(dep, "\"'")
		if dep == "" {
			continue
		}

		name, version := s.parseDepString(dep)
		if name != "" {
			metadata := map[string]interface{}{
				"manager": "setuptools",
			}

			pkg := s.base.CreatePackage(name, version, filePath, false, metadata)
			packages = append(packages, pkg)
		}
	}

	return packages
}

// parseDepString parses a dependency string to extract name and version
func (s *SetupPyParser) parseDepString(dep string) (string, string) {
	// Common patterns for setup.py dependencies
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
				version = s.cleanPythonVersion(version)
			}
			return name, version
		}
	}

	return "", ""
}

// cleanPythonVersion cleans and normalizes python version strings
func (s *SetupPyParser) cleanPythonVersion(version string) string {
	prefixes := []string{"^", "~", ">=", "<=", ">", "<", "==", "!="}
	return s.base.GetVersionCleaner().CleanVersion(version, prefixes)
}
