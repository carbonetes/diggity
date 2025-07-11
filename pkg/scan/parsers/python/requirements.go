package python

import (
	"bufio"
	"os"
	"regexp"
	"strings"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/scan/parsers/common"
)

// RequirementsParser handles parsing of requirements.txt files
type RequirementsParser struct {
	base *common.BaseParser
}

// NewRequirementsParser creates a new requirements parser
func NewRequirementsParser(base *common.BaseParser) *RequirementsParser {
	return &RequirementsParser{base: base}
}

// Parse parses requirements.txt files
func (r *RequirementsParser) Parse(filePath string) ([]model.Package, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var packages []model.Package
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Skip -r includes and other options
		if strings.HasPrefix(line, "-") {
			continue
		}

		pkg := r.parseRequirementLine(line, filePath)
		if pkg != nil {
			packages = append(packages, *pkg)
		}
	}

	return packages, scanner.Err()
}

// parseRequirementLine parses a single requirement line
func (r *RequirementsParser) parseRequirementLine(line, filePath string) *model.Package {
	line = r.cleanLine(line)
	if line == "" {
		return nil
	}

	name, version := r.extractNameAndVersion(line)
	if name == "" {
		return nil
	}

	isDev := r.isDevDependency(filePath)
	metadata := map[string]interface{}{
		"manager": "pip",
	}

	pkg := r.base.CreatePackage(name, version, filePath, isDev, metadata)
	return &pkg
}

// cleanLine removes comments and trims whitespace
func (r *RequirementsParser) cleanLine(line string) string {
	// Remove inline comments
	if idx := strings.Index(line, "#"); idx != -1 {
		line = line[:idx]
	}
	return strings.TrimSpace(line)
}

// extractNameAndVersion extracts package name and version from requirement line
func (r *RequirementsParser) extractNameAndVersion(line string) (string, string) {
	patterns := []string{
		`^([a-zA-Z0-9\-_.]+)\s*([<>=!~]+)\s*([0-9][^,\s;]*)?`, // package>=1.0.0
		`^([a-zA-Z0-9\-_.]+)\s*==\s*([0-9][^\s,;]*)`,          // package==1.0.0
		`^([a-zA-Z0-9\-_.]+)\s*$`,                             // package (no version)
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(line)

		if len(matches) >= 2 {
			name := matches[1]
			version := r.extractVersionFromMatches(matches)
			return name, r.cleanPythonVersion(version)
		}
	}

	return "", ""
}

// extractVersionFromMatches extracts version from regex matches
func (r *RequirementsParser) extractVersionFromMatches(matches []string) string {
	if len(matches) >= 3 && matches[2] != "" {
		if len(matches) >= 4 && matches[3] != "" {
			return matches[3]
		}
		return matches[2]
	}
	return ""
}

// isDevDependency checks if the file indicates development dependencies
func (r *RequirementsParser) isDevDependency(filePath string) bool {
	return strings.Contains(filePath, "dev") || strings.Contains(filePath, "test")
}

// cleanPythonVersion cleans and normalizes python version strings
func (r *RequirementsParser) cleanPythonVersion(version string) string {
	prefixes := []string{"^", "~", ">=", "<=", ">", "<", "==", "!="}
	return r.base.GetVersionCleaner().CleanVersion(version, prefixes)
}
