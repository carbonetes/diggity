package ruby

import (
	"bufio"
	"os"
	"regexp"
	"strings"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/scan/parsers/common"
)

// GemfileParser handles parsing of Gemfile files
type GemfileParser struct {
	base *common.BaseParser
}

// NewGemfileParser creates a new Gemfile parser
func NewGemfileParser(base *common.BaseParser) *GemfileParser {
	return &GemfileParser{base: base}
}

// Parse parses Gemfile files
func (g *GemfileParser) Parse(filePath string) ([]model.Package, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var packages []model.Package
	scanner := bufio.NewScanner(file)

	currentGroup := ""

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Track group context
		if group := g.extractGroup(line); group != "" {
			currentGroup = group
			continue
		}

		// Parse gem line
		pkg := g.parseGemLine(line, filePath, currentGroup)
		if pkg != nil {
			packages = append(packages, *pkg)
		}
	}

	return packages, scanner.Err()
}

// extractGroup extracts group information from group lines
func (g *GemfileParser) extractGroup(line string) string {
	// Look for group declarations like: group :development do
	groupPattern := `group\s+:(\w+)`
	re := regexp.MustCompile(groupPattern)
	matches := re.FindStringSubmatch(line)

	if len(matches) > 1 {
		return matches[1]
	}

	// Reset group on "end"
	if strings.Contains(line, "end") {
		return "reset"
	}

	return ""
}

// parseGemLine parses a gem declaration line
func (g *GemfileParser) parseGemLine(line, filePath, group string) *model.Package {
	// Patterns for gem declarations:
	// gem 'name', 'version'
	// gem "name", "version"
	// gem 'name', '~> version'
	patterns := []string{
		`gem\s+['"]([\w\-_]+)['"],\s*['"]([\w\-\.\~>\s<>=!]+)['"]`,
		`gem\s+['"]([\w\-_]+)['"]`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(line)

		if len(matches) >= 2 {
			name := matches[1]
			version := ""
			if len(matches) > 2 {
				version = g.cleanRubyVersion(matches[2])
			}

			isDev := g.isDevGroup(group, line)

			metadata := map[string]interface{}{
				"manager": "bundler",
			}

			if group != "" && group != "reset" {
				metadata["group"] = group
			}

			pkg := g.base.CreatePackage(name, version, filePath, isDev, metadata)
			return &pkg
		}
	}

	return nil
}

// isDevGroup determines if a gem is in a development group
func (g *GemfileParser) isDevGroup(group, line string) bool {
	devGroups := []string{"development", "test", "debug"}

	// Check current group
	for _, devGroup := range devGroups {
		if group == devGroup {
			return true
		}
	}

	// Check for inline group specification
	for _, devGroup := range devGroups {
		if strings.Contains(line, ":"+devGroup) {
			return true
		}
	}

	return false
}

// cleanRubyVersion cleans and normalizes ruby version strings
func (g *GemfileParser) cleanRubyVersion(version string) string {
	prefixes := []string{"~>", ">=", "<=", ">", "<", "="}
	return g.base.GetVersionCleaner().CleanVersion(version, prefixes)
}
