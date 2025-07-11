package python

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/scan/parsers/common"
)

// PipfileParser handles parsing of Pipfile files
type PipfileParser struct {
	base *common.BaseParser
}

// NewPipfileParser creates a new Pipfile parser
func NewPipfileParser(base *common.BaseParser) *PipfileParser {
	return &PipfileParser{base: base}
}

// Parse parses Pipfile files
func (p *PipfileParser) Parse(filePath string) ([]model.Package, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var packages []model.Package
	contentStr := string(content)

	// Parse both production and development packages
	prodDeps := p.parsePipfileSection(contentStr, "packages", false)
	devDeps := p.parsePipfileSection(contentStr, "dev-packages", true)

	// Combine dependencies
	allDeps := make(map[string]pipfileDependency)
	for name, dep := range prodDeps {
		allDeps[name] = dep
	}
	for name, dep := range devDeps {
		allDeps[name] = dep
	}

	// Convert to packages
	for name, dep := range allDeps {
		metadata := map[string]interface{}{
			"manager": "pipenv",
		}

		version := p.cleanPythonVersion(dep.version)
		pkg := p.base.CreatePackage(name, version, filePath, dep.isDev, metadata)
		packages = append(packages, pkg)
	}

	return packages, nil
}

// pipfileDependency represents a dependency in a Pipfile
type pipfileDependency struct {
	version string
	isDev   bool
}

// parsePipfileSection parses a specific section from Pipfile
func (p *PipfileParser) parsePipfileSection(content, sectionName string, isDev bool) map[string]pipfileDependency {
	pattern := fmt.Sprintf(`(?s)\[%s\](.*?)(?:\[|$)`, regexp.QuoteMeta(sectionName))
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(content)

	if len(matches) < 2 {
		return nil
	}

	section := matches[1]
	result := make(map[string]pipfileDependency)

	// Parse key = value pairs
	lines := strings.Split(section, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.Contains(line, "=") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				key := p.cleanQuotes(strings.TrimSpace(parts[0]))
				value := p.cleanQuotes(strings.TrimSpace(parts[1]))

				result[key] = pipfileDependency{
					version: value,
					isDev:   isDev,
				}
			}
		}
	}

	return result
}

// cleanQuotes removes quotes from strings
func (p *PipfileParser) cleanQuotes(str string) string {
	return strings.Trim(str, "\"'")
}

// cleanPythonVersion cleans and normalizes python version strings
func (p *PipfileParser) cleanPythonVersion(version string) string {
	prefixes := []string{"^", "~", ">=", "<=", ">", "<", "==", "!="}
	return p.base.GetVersionCleaner().CleanVersion(version, prefixes)
}
