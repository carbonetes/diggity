package ruby

import (
	"bufio"
	"os"
	"regexp"
	"strings"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/scan/parsers/common"
)

// GemspecParser handles parsing of .gemspec files
type GemspecParser struct {
	base *common.BaseParser
}

// NewGemspecParser creates a new .gemspec parser
func NewGemspecParser(base *common.BaseParser) *GemspecParser {
	return &GemspecParser{base: base}
}

// Parse parses .gemspec files
func (g *GemspecParser) Parse(filePath string) ([]model.Package, error) {
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

		// Parse dependency lines
		pkg := g.parseDependencyLine(line, filePath)
		if pkg != nil {
			packages = append(packages, *pkg)
		}
	}

	return packages, scanner.Err()
}

// parseDependencyLine parses dependency declarations in .gemspec files
func (g *GemspecParser) parseDependencyLine(line, filePath string) *model.Package {
	// Patterns for gemspec dependencies:
	// spec.add_dependency 'gem_name', 'version'
	// spec.add_development_dependency 'gem_name', 'version'
	// s.add_dependency 'gem_name', 'version'
	patterns := map[string]bool{
		`(?:spec|s)\.add_dependency\s+['"]([\w\-_]+)['"](?:,\s*['"]([\w\-\.\~>\s<>=!]+)['"])?`:             false, // production
		`(?:spec|s)\.add_development_dependency\s+['"]([\w\-_]+)['"](?:,\s*['"]([\w\-\.\~>\s<>=!]+)['"])?`: true,  // development
		`(?:spec|s)\.add_runtime_dependency\s+['"]([\w\-_]+)['"](?:,\s*['"]([\w\-\.\~>\s<>=!]+)['"])?`:     false, // production
	}

	for pattern, isDev := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(line)

		if len(matches) >= 2 {
			name := matches[1]
			version := ""
			if len(matches) > 2 && matches[2] != "" {
				version = g.cleanRubyVersion(matches[2])
			}

			metadata := map[string]interface{}{
				"manager": "rubygems",
			}

			pkg := g.base.CreatePackage(name, version, filePath, isDev, metadata)
			return &pkg
		}
	}

	return nil
}

// cleanRubyVersion cleans and normalizes ruby version strings
func (g *GemspecParser) cleanRubyVersion(version string) string {
	prefixes := []string{"~>", ">=", "<=", ">", "<", "="}
	return g.base.GetVersionCleaner().CleanVersion(version, prefixes)
}
