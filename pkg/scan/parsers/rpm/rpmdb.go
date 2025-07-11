package rpm

import (
	"bufio"
	"os"
	"regexp"
	"strings"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/scan/parsers/common"
)

// RpmDbParser handles parsing of RPM database files
type RpmDbParser struct {
	base *common.BaseParser
}

// NewRpmDbParser creates a new RPM database parser
func NewRpmDbParser(base *common.BaseParser) *RpmDbParser {
	return &RpmDbParser{base: base}
}

// Parse parses RPM database files (usually requires rpm -qa output)
func (p *RpmDbParser) Parse(filePath string) ([]model.Package, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var packages []model.Package
	scanner := bufio.NewScanner(file)

	// Pattern to match RPM package names: name-version-release.arch
	rpmPattern := regexp.MustCompile(`^(.+)-([^-]+)-([^-]+)\.(.+)$`)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		// Try to parse as RPM package format
		matches := rpmPattern.FindStringSubmatch(line)
		if len(matches) == 5 {
			name := matches[1]
			version := matches[2]
			release := matches[3]
			arch := matches[4]

			metadata := map[string]interface{}{
				"manager":      "rpm",
				"release":      release,
				"architecture": arch,
			}

			cleanVersion := p.cleanRpmVersion(version)
			pkg := p.base.CreatePackage(name, cleanVersion, filePath, false, metadata)
			packages = append(packages, pkg)
		}
	}

	return packages, scanner.Err()
}

// cleanRpmVersion cleans and normalizes RPM version strings
func (p *RpmDbParser) cleanRpmVersion(version string) string {
	version = strings.TrimSpace(version)

	// RPM versions can have epoch (e.g., "1:2.3.4")
	if colonIndex := strings.Index(version, ":"); colonIndex > 0 {
		version = version[colonIndex+1:]
	}

	return strings.TrimSpace(version)
}
