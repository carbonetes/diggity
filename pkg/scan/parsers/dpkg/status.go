package dpkg

import (
	"bufio"
	"os"
	"strings"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/scan/parsers/common"
)

// StatusParser handles parsing of /var/lib/dpkg/status files
type StatusParser struct {
	base *common.BaseParser
}

// NewStatusParser creates a new DPKG status parser
func NewStatusParser(base *common.BaseParser) *StatusParser {
	return &StatusParser{base: base}
}

// Parse parses DPKG status files
func (p *StatusParser) Parse(filePath string) ([]model.Package, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var packages []model.Package
	scanner := bufio.NewScanner(file)

	var currentPackage map[string]string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Empty line indicates end of package entry
		if line == "" {
			if currentPackage != nil {
				if pkg := p.createPackageFromEntry(currentPackage, filePath); pkg != nil {
					packages = append(packages, *pkg)
				}
				currentPackage = nil
			}
			continue
		}

		// Start new package entry
		if currentPackage == nil {
			currentPackage = make(map[string]string)
		}

		// Parse field: value pairs
		if colonIndex := strings.Index(line, ":"); colonIndex > 0 {
			field := strings.TrimSpace(line[:colonIndex])
			value := strings.TrimSpace(line[colonIndex+1:])
			currentPackage[field] = value
		}
	}

	// Handle last package if file doesn't end with empty line
	if currentPackage != nil {
		if pkg := p.createPackageFromEntry(currentPackage, filePath); pkg != nil {
			packages = append(packages, *pkg)
		}
	}

	return packages, scanner.Err()
}

// createPackageFromEntry creates a package model from a DPKG status entry
func (p *StatusParser) createPackageFromEntry(entry map[string]string, filePath string) *model.Package {
	name, exists := entry["Package"]
	if !exists || name == "" {
		return nil
	}

	// Only include installed packages
	status, exists := entry["Status"]
	if !exists || !strings.Contains(status, "install ok installed") {
		return nil
	}

	version := entry["Version"]
	architecture := entry["Architecture"]
	description := entry["Description"]
	section := entry["Section"]
	priority := entry["Priority"]
	essential := entry["Essential"]
	depends := entry["Depends"]
	recommends := entry["Recommends"]
	suggests := entry["Suggests"]
	conflicts := entry["Conflicts"]
	replaces := entry["Replaces"]

	metadata := map[string]interface{}{
		"manager":      "dpkg",
		"status":       status,
		"architecture": architecture,
		"description":  description,
		"section":      section,
		"priority":     priority,
		"essential":    essential,
		"depends":      depends,
		"recommends":   recommends,
		"suggests":     suggests,
		"conflicts":    conflicts,
		"replaces":     replaces,
	}

	cleanVersion := p.cleanDpkgVersion(version)
	pkg := p.base.CreatePackage(name, cleanVersion, filePath, false, metadata)

	return &pkg
}

// cleanDpkgVersion cleans and normalizes DPKG version strings
func (p *StatusParser) cleanDpkgVersion(version string) string {
	version = strings.TrimSpace(version)

	// DPKG versions can have epoch (e.g., "1:2.3.4-1")
	if colonIndex := strings.Index(version, ":"); colonIndex > 0 {
		version = version[colonIndex+1:]
	}

	// Remove Debian revision (everything after last "-")
	if dashIndex := strings.LastIndex(version, "-"); dashIndex > 0 {
		version = version[:dashIndex]
	}

	return strings.TrimSpace(version)
}
