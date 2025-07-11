package cargo

import (
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/scan/parsers/common"
)

// CargoTomlParser handles parsing of Cargo.toml files
type CargoTomlParser struct {
	base *common.BaseParser
}

// NewCargoTomlParser creates a new Cargo.toml parser
func NewCargoTomlParser(base *common.BaseParser) *CargoTomlParser {
	return &CargoTomlParser{base: base}
}

// Parse parses Cargo.toml files
func (c *CargoTomlParser) Parse(filePath string) ([]model.Package, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var packages []model.Package
	contentStr := string(content)

	// Parse [dependencies] section
	deps := c.parseTomlSection(contentStr, "dependencies", false)
	packages = append(packages, deps...)

	// Parse [dev-dependencies] section
	devDeps := c.parseTomlSection(contentStr, "dev-dependencies", true)
	packages = append(packages, devDeps...)

	return packages, nil
}

// parseTomlSection parses a specific TOML section for dependencies
func (c *CargoTomlParser) parseTomlSection(content, sectionName string, isDev bool) []model.Package {
	var packages []model.Package

	// Find the section
	sectionPattern := `(?s)\[` + regexp.QuoteMeta(sectionName) + `\](.*?)(?:\[|$)`
	re := regexp.MustCompile(sectionPattern)
	matches := re.FindStringSubmatch(content)

	if len(matches) < 2 {
		return packages
	}

	sectionContent := matches[1]
	lines := strings.Split(sectionContent, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.Contains(line, "=") {
			pkg := c.parseDependencyLine(line, isDev)
			if pkg != nil {
				packages = append(packages, *pkg)
			}
		}
	}

	return packages
}

// parseDependencyLine parses a single dependency line
func (c *CargoTomlParser) parseDependencyLine(line string, isDev bool) *model.Package {
	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return nil
	}

	name := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])

	// Remove quotes
	name = strings.Trim(name, "\"'")

	version := c.extractVersion(value)
	if version == "" {
		return nil
	}

	metadata := map[string]interface{}{
		"manager": "cargo",
	}

	// Clean version
	cleanVersion := c.cleanCargoVersion(version)

	pkg := c.base.CreatePackage(name, cleanVersion, "", isDev, metadata)
	return &pkg
}

// extractVersion extracts version from various Cargo.toml dependency formats
func (c *CargoTomlParser) extractVersion(value string) string {
	// Simple version: name = "1.0.0"
	if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
		return strings.Trim(value, "\"")
	}

	// Object format: name = { version = "1.0.0", features = ["..."] }
	if strings.Contains(value, "version") {
		versionPattern := `version\s*=\s*"([^"]+)"`
		re := regexp.MustCompile(versionPattern)
		matches := re.FindStringSubmatch(value)
		if len(matches) > 1 {
			return matches[1]
		}
	}

	return ""
}

// cleanCargoVersion cleans and normalizes cargo version strings
func (c *CargoTomlParser) cleanCargoVersion(version string) string {
	prefixes := []string{"^", "~", ">=", "<=", ">", "<", "="}
	return c.base.GetVersionCleaner().CleanVersion(version, prefixes)
}
