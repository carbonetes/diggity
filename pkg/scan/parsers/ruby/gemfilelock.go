package ruby

import (
	"bufio"
	"os"
	"regexp"
	"strings"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/scan/parsers/common"
)

// GemfileLockParser handles parsing of Gemfile.lock files
type GemfileLockParser struct {
	base *common.BaseParser
}

// NewGemfileLockParser creates a new Gemfile.lock parser
func NewGemfileLockParser(base *common.BaseParser) *GemfileLockParser {
	return &GemfileLockParser{base: base}
}

// Parse parses Gemfile.lock files
func (g *GemfileLockParser) Parse(filePath string) ([]model.Package, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var packages []model.Package
	scanner := bufio.NewScanner(file)

	inGemsSection := false

	for scanner.Scan() {
		line := scanner.Text()

		// Check for GEM section
		if strings.Contains(line, "GEM") {
			inGemsSection = true
			continue
		}

		// Check for end of section
		if inGemsSection && g.isNewSection(line) {
			inGemsSection = false
			continue
		}

		if inGemsSection {
			pkg := g.parseGemLockLine(line, filePath)
			if pkg != nil {
				packages = append(packages, *pkg)
			}
		}
	}

	return packages, scanner.Err()
}

// isNewSection checks if we've reached a new section (like PLATFORMS, DEPENDENCIES, etc.)
func (g *GemfileLockParser) isNewSection(line string) bool {
	sections := []string{"PLATFORMS", "DEPENDENCIES", "BUNDLED WITH"}
	trimmed := strings.TrimSpace(line)

	for _, section := range sections {
		if trimmed == section {
			return true
		}
	}

	return false
}

// parseGemLockLine parses a line from the GEM section of Gemfile.lock
func (g *GemfileLockParser) parseGemLockLine(line, filePath string) *model.Package {
	// Skip lines that don't contain gem specifications
	if !strings.Contains(line, "(") || !strings.Contains(line, ")") {
		return nil
	}

	// Pattern for gem lines: "    gemname (version)"
	pattern := `^\s+([\w\-_]+)\s+\(([\d\.\w\-\+]+)\)`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(line)

	if len(matches) >= 3 {
		name := matches[1]
		version := matches[2]

		metadata := map[string]interface{}{
			"manager": "bundler",
			"locked":  true,
		}

		cleanVersion := g.cleanRubyVersion(version)
		pkg := g.base.CreatePackage(name, cleanVersion, filePath, false, metadata)
		return &pkg
	}

	return nil
}

// cleanRubyVersion cleans and normalizes ruby version strings
func (g *GemfileLockParser) cleanRubyVersion(version string) string {
	prefixes := []string{"~>", ">=", "<=", ">", "<", "="}
	return g.base.GetVersionCleaner().CleanVersion(version, prefixes)
}
