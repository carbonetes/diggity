package npm

import (
	"bufio"
	"os"
	"regexp"
	"strings"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/scan/parsers/common"
)

// YarnLockParser handles parsing of yarn.lock files
type YarnLockParser struct {
	base *common.BaseParser
}

// NewYarnLockParser creates a new yarn.lock parser
func NewYarnLockParser(base *common.BaseParser) *YarnLockParser {
	return &YarnLockParser{base: base}
}

// Parse parses yarn.lock files
func (p *YarnLockParser) Parse(filePath string) ([]model.Package, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var packages []model.Package
	scanner := bufio.NewScanner(file)

	currentPackage := ""

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if p.shouldSkipLine(line) {
			continue
		}

		if p.isPackageDeclaration(line) {
			currentPackage = p.extractPackageName(line)
			continue
		}

		if p.isVersionLine(line) && currentPackage != "" {
			pkg := p.createPackageFromVersionLine(line, currentPackage, filePath)
			if pkg != nil {
				packages = append(packages, *pkg)
			}
			currentPackage = ""
		}
	}

	return packages, scanner.Err()
}

// shouldSkipLine checks if a line should be skipped
func (p *YarnLockParser) shouldSkipLine(line string) bool {
	return line == "" || strings.HasPrefix(line, "#")
}

// isVersionLine checks if a line contains version information
func (p *YarnLockParser) isVersionLine(line string) bool {
	return strings.HasPrefix(line, "version ")
}

// createPackageFromVersionLine creates a package from a version line
func (p *YarnLockParser) createPackageFromVersionLine(line, packageName, filePath string) *model.Package {
	versionMatch := regexp.MustCompile(`version\s+"([^"]+)"`).FindStringSubmatch(line)
	if len(versionMatch) <= 1 {
		return nil
	}

	version := versionMatch[1]
	metadata := map[string]interface{}{
		"manager": "yarn",
		"locked":  true,
	}

	cleanVersion := p.cleanNPMVersion(version)
	pkg := p.base.CreatePackage(packageName, cleanVersion, filePath, false, metadata)
	return &pkg
}

// isPackageDeclaration checks if a line is a package declaration
func (p *YarnLockParser) isPackageDeclaration(line string) bool {
	// Yarn lock package declarations typically look like:
	// "package@version", "package@^version", etc.
	return strings.Contains(line, "@") && (strings.HasSuffix(line, ":") || strings.HasSuffix(line, ", "))
}

// extractPackageName extracts package name from a package declaration line
func (p *YarnLockParser) extractPackageName(line string) string {
	// Remove trailing colon and comma
	line = strings.TrimSuffix(line, ":")
	line = strings.TrimSuffix(line, ",")
	line = strings.TrimSpace(line)
	line = strings.Trim(line, "\"'")

	// Extract package name (everything before @version)
	parts := strings.Split(line, "@")
	if len(parts) >= 2 {
		// Handle scoped packages like @scope/package@version
		if strings.HasPrefix(line, "@") && len(parts) >= 3 {
			return "@" + parts[1]
		}
		return parts[0]
	}

	return ""
}

// cleanNPMVersion cleans and normalizes npm version strings
func (p *YarnLockParser) cleanNPMVersion(version string) string {
	prefixes := []string{"^", "~", ">=", "<=", ">", "<", "="}
	return p.base.GetVersionCleaner().CleanVersion(version, prefixes)
}
