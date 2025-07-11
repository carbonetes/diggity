package swift

import (
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/scan/parsers/common"
)

// PodfileParser handles parsing of Podfile and Podfile.lock files
type PodfileParser struct {
	base *common.BaseParser
}

// NewPodfileParser creates a new Podfile parser
func NewPodfileParser(base *common.BaseParser) *PodfileParser {
	return &PodfileParser{base: base}
}

// Parse parses Podfile files
func (p *PodfileParser) Parse(filePath string) ([]model.Package, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var packages []model.Package

	// Simple regex to find pod dependencies
	podRegex := regexp.MustCompile(`pod\s+['"]([^'"]+)['"](?:\s*,\s*['"]([^'"]+)['"])?`)
	matches := podRegex.FindAllStringSubmatch(string(content), -1)

	for _, match := range matches {
		if len(match) > 1 {
			name := match[1]
			version := ""
			if len(match) > 2 {
				version = match[2]
			}

			metadata := map[string]interface{}{
				"manager": "cocoapods",
			}

			cleanVersion := p.cleanCocoaPodsVersion(version)
			pkg := p.base.CreatePackage(name, cleanVersion, filePath, false, metadata)
			packages = append(packages, pkg)
		}
	}

	return packages, nil
}

// ParseLock parses Podfile.lock files
func (p *PodfileParser) ParseLock(filePath string) ([]model.Package, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var packages []model.Package

	// Parse PODS section
	lines := strings.Split(string(content), "\n")
	inPodsSection := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "PODS:") {
			inPodsSection = true
			continue
		}

		if inPodsSection {
			if strings.HasPrefix(line, "DEPENDENCIES:") ||
				strings.HasPrefix(line, "SPEC REPOS:") ||
				strings.HasPrefix(line, "CHECKOUT OPTIONS:") {
				break
			}

			// Parse pod entry: "- PodName (version)"
			podRegex := regexp.MustCompile(`^-\s+([^(]+)\s*\(([^)]+)\)`)
			matches := podRegex.FindStringSubmatch(line)
			if len(matches) > 2 {
				name := strings.TrimSpace(matches[1])
				version := strings.TrimSpace(matches[2])

				metadata := map[string]interface{}{
					"manager": "cocoapods",
				}

				cleanVersion := p.cleanCocoaPodsVersion(version)
				pkg := p.base.CreatePackage(name, cleanVersion, filePath, true, metadata)
				packages = append(packages, pkg)
			}
		}
	}

	return packages, nil
}

// cleanCocoaPodsVersion cleans and normalizes CocoaPods version strings
func (p *PodfileParser) cleanCocoaPodsVersion(version string) string {
	version = strings.TrimSpace(version)
	prefixes := []string{"~>", ">=", "<=", ">", "<", "=", "~"}
	return p.base.GetVersionCleaner().CleanVersion(version, prefixes)
}
