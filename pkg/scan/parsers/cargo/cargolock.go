package cargo

import (
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/scan/parsers/common"
)

// CargoLockParser handles parsing of Cargo.lock files
type CargoLockParser struct {
	base *common.BaseParser
}

// NewCargoLockParser creates a new Cargo.lock parser
func NewCargoLockParser(base *common.BaseParser) *CargoLockParser {
	return &CargoLockParser{base: base}
}

// Parse parses Cargo.lock files
func (c *CargoLockParser) Parse(filePath string) ([]model.Package, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var packages []model.Package
	contentStr := string(content)

	// Parse [[package]] sections
	packageSections := c.extractPackageSections(contentStr)

	for _, section := range packageSections {
		pkg := c.parsePackageSection(section, filePath)
		if pkg != nil {
			packages = append(packages, *pkg)
		}
	}

	return packages, nil
}

// extractPackageSections extracts all [[package]] sections from Cargo.lock
func (c *CargoLockParser) extractPackageSections(content string) []string {
	// Find all [[package]] sections
	pattern := `(?s)\[\[package\]\](.*?)(?:\[\[|$)`
	re := regexp.MustCompile(pattern)
	matches := re.FindAllStringSubmatch(content, -1)

	var sections []string
	for _, match := range matches {
		if len(match) > 1 {
			sections = append(sections, match[1])
		}
	}

	return sections
}

// parsePackageSection parses a single [[package]] section
func (c *CargoLockParser) parsePackageSection(section, filePath string) *model.Package {
	lines := strings.Split(section, "\n")

	var name, version, source string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "name = ") {
			name = c.extractValue(line)
		} else if strings.HasPrefix(line, "version = ") {
			version = c.extractValue(line)
		} else if strings.HasPrefix(line, "source = ") {
			source = c.extractValue(line)
		}
	}

	if name == "" || version == "" {
		return nil
	}

	metadata := map[string]interface{}{
		"manager": "cargo",
		"locked":  true,
	}

	if source != "" {
		metadata["source"] = source
	}

	pkg := c.base.CreatePackage(name, version, filePath, false, metadata)
	return &pkg
}

// extractValue extracts the value from a TOML key-value line
func (c *CargoLockParser) extractValue(line string) string {
	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return ""
	}

	value := strings.TrimSpace(parts[1])
	return strings.Trim(value, "\"'")
}
