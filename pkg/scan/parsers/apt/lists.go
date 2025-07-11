package apt

import (
	"bufio"
	"os"
	"strings"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/scan/parsers/common"
)

// ListsParser handles parsing of /var/lib/apt/lists/**/Packages files
type ListsParser struct {
	base *common.BaseParser
}

// NewListsParser creates a new APT lists parser
func NewListsParser(base *common.BaseParser) *ListsParser {
	return &ListsParser{base: base}
}

// Parse parses APT Packages files
func (p *ListsParser) Parse(filePath string) ([]model.Package, error) {
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

// createPackageFromEntry creates a package model from an APT Packages entry
func (p *ListsParser) createPackageFromEntry(entry map[string]string, filePath string) *model.Package {
	name, exists := entry["Package"]
	if !exists || name == "" {
		return nil
	}

	version := entry["Version"]
	architecture := entry["Architecture"]
	description := entry["Description"]
	section := entry["Section"]
	priority := entry["Priority"]
	filename := entry["Filename"]
	size := entry["Size"]
	sha256 := entry["SHA256"]

	metadata := map[string]interface{}{
		"manager":      "apt",
		"architecture": architecture,
		"description":  description,
		"section":      section,
		"priority":     priority,
		"filename":     filename,
		"size":         size,
		"sha256":       sha256,
	}

	cleanVersion := p.cleanAptVersion(version)
	pkg := p.base.CreatePackage(name, cleanVersion, filePath, false, metadata)

	return &pkg
}

// cleanAptVersion cleans and normalizes APT version strings
func (p *ListsParser) cleanAptVersion(version string) string {
	version = strings.TrimSpace(version)

	// APT versions can have epoch (e.g., "1:2.3.4-1")
	if colonIndex := strings.Index(version, ":"); colonIndex > 0 {
		version = version[colonIndex+1:]
	}

	// Remove Debian revision (everything after last "-")
	if dashIndex := strings.LastIndex(version, "-"); dashIndex > 0 {
		version = version[:dashIndex]
	}

	return strings.TrimSpace(version)
}
