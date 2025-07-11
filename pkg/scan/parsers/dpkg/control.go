package dpkg

import (
	"bufio"
	"os"
	"strings"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/scan/parsers/common"
)

// ControlParser handles parsing of DEBIAN/control files
type ControlParser struct {
	base *common.BaseParser
}

// NewControlParser creates a new Debian control parser
func NewControlParser(base *common.BaseParser) *ControlParser {
	return &ControlParser{base: base}
}

// Parse parses Debian control files
func (p *ControlParser) Parse(filePath string) ([]model.Package, error) {
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

// createPackageFromEntry creates a package model from a Debian control entry
func (p *ControlParser) createPackageFromEntry(entry map[string]string, filePath string) *model.Package {
	name, exists := entry["Package"]
	if !exists || name == "" {
		return nil
	}

	version := entry["Version"]
	architecture := entry["Architecture"]
	description := entry["Description"]
	section := entry["Section"]
	priority := entry["Priority"]
	depends := entry["Depends"]
	recommends := entry["Recommends"]
	suggests := entry["Suggests"]
	conflicts := entry["Conflicts"]
	replaces := entry["Replaces"]
	provides := entry["Provides"]
	essential := entry["Essential"]

	metadata := map[string]interface{}{
		"manager":      "dpkg",
		"type":         "control",
		"architecture": architecture,
		"description":  description,
		"section":      section,
		"priority":     priority,
		"depends":      depends,
		"recommends":   recommends,
		"suggests":     suggests,
		"conflicts":    conflicts,
		"replaces":     replaces,
		"provides":     provides,
		"essential":    essential,
	}

	cleanVersion := p.cleanDebVersion(version)
	pkg := p.base.CreatePackage(name, cleanVersion, filePath, false, metadata)

	return &pkg
}

// cleanDebVersion cleans and normalizes Debian version strings
func (p *ControlParser) cleanDebVersion(version string) string {
	version = strings.TrimSpace(version)

	// Debian versions can have epoch (e.g., "1:2.3.4-1")
	if colonIndex := strings.Index(version, ":"); colonIndex > 0 {
		version = version[colonIndex+1:]
	}

	// Remove Debian revision (everything after last "-")
	if dashIndex := strings.LastIndex(version, "-"); dashIndex > 0 {
		version = version[:dashIndex]
	}

	return strings.TrimSpace(version)
}
