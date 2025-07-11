package rpm

import (
	"bufio"
	"os"
	"regexp"
	"strings"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/scan/parsers/common"
)

// RpmSpecParser handles parsing of RPM .spec files
type RpmSpecParser struct {
	base *common.BaseParser
}

// NewRpmSpecParser creates a new RPM spec parser
func NewRpmSpecParser(base *common.BaseParser) *RpmSpecParser {
	return &RpmSpecParser{base: base}
}

// Parse parses RPM .spec files
func (p *RpmSpecParser) Parse(filePath string) ([]model.Package, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var packages []model.Package
	scanner := bufio.NewScanner(file)

	spec := make(map[string]string)

	// Patterns for different spec file sections
	namePattern := regexp.MustCompile(`^Name:\s*(.+)$`)
	versionPattern := regexp.MustCompile(`^Version:\s*(.+)$`)
	releasePattern := regexp.MustCompile(`^Release:\s*(.+)$`)
	summaryPattern := regexp.MustCompile(`^Summary:\s*(.+)$`)
	requiresPattern := regexp.MustCompile(`^Requires:\s*(.+)$`)
	buildRequiresPattern := regexp.MustCompile(`^BuildRequires:\s*(.+)$`)
	sectionPattern := regexp.MustCompile(`^%(\w+)`)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip comments and empty lines
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Check for section headers (for future use)
		if matches := sectionPattern.FindStringSubmatch(line); len(matches) > 1 {
			// Section detected, could be used for future parsing improvements
			continue
		}

		// Parse main spec fields
		if matches := namePattern.FindStringSubmatch(line); len(matches) > 1 {
			spec["name"] = strings.TrimSpace(matches[1])
		} else if matches := versionPattern.FindStringSubmatch(line); len(matches) > 1 {
			spec["version"] = strings.TrimSpace(matches[1])
		} else if matches := releasePattern.FindStringSubmatch(line); len(matches) > 1 {
			spec["release"] = strings.TrimSpace(matches[1])
		} else if matches := summaryPattern.FindStringSubmatch(line); len(matches) > 1 {
			spec["summary"] = strings.TrimSpace(matches[1])
		} else if matches := requiresPattern.FindStringSubmatch(line); len(matches) > 1 {
			deps := spec["requires"]
			if deps == "" {
				deps = strings.TrimSpace(matches[1])
			} else {
				deps += ", " + strings.TrimSpace(matches[1])
			}
			spec["requires"] = deps
		} else if matches := buildRequiresPattern.FindStringSubmatch(line); len(matches) > 1 {
			buildDeps := spec["build_requires"]
			if buildDeps == "" {
				buildDeps = strings.TrimSpace(matches[1])
			} else {
				buildDeps += ", " + strings.TrimSpace(matches[1])
			}
			spec["build_requires"] = buildDeps
		}
	}

	// Create package from spec data
	if pkg := p.createSpecPackage(spec, filePath); pkg != nil {
		packages = append(packages, *pkg)
	}

	return packages, scanner.Err()
}

// createSpecPackage creates a package model from a RPM spec entry
func (p *RpmSpecParser) createSpecPackage(spec map[string]string, filePath string) *model.Package {
	name := spec["name"]
	if name == "" {
		return nil
	}

	version := spec["version"]
	release := spec["release"]
	summary := spec["summary"]
	requires := spec["requires"]
	buildRequires := spec["build_requires"]

	metadata := map[string]interface{}{
		"manager":        "rpm",
		"type":           "spec",
		"release":        release,
		"summary":        summary,
		"requires":       requires,
		"build_requires": buildRequires,
	}

	cleanVersion := p.cleanSpecVersion(version)
	pkg := p.base.CreatePackage(name, cleanVersion, filePath, false, metadata)

	return &pkg
}

// cleanSpecVersion cleans and normalizes RPM spec version strings
func (p *RpmSpecParser) cleanSpecVersion(version string) string {
	version = strings.TrimSpace(version)

	// Remove macro references like %{version}
	if strings.Contains(version, "%{") {
		return ""
	}

	return version
}
