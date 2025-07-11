package alpine

import (
	"bufio"
	"os"
	"strings"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/scan/parsers/common"
)

// ApkDbParser handles parsing of APK database files (/lib/apk/db/installed)
type ApkDbParser struct {
	base *common.BaseParser
}

// NewApkDbParser creates a new APK database parser
func NewApkDbParser(base *common.BaseParser) *ApkDbParser {
	return &ApkDbParser{base: base}
}

// Parse parses APK installed database files
func (p *ApkDbParser) Parse(filePath string) ([]model.Package, error) {
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

		// Start new package entry if this is a package line
		if strings.HasPrefix(line, "P:") {
			currentPackage = make(map[string]string)
		}

		if currentPackage == nil {
			continue
		}

		// Parse APK database fields
		if len(line) > 2 && line[1] == ':' {
			field := line[0:1]
			value := strings.TrimSpace(line[2:])

			switch field {
			case "P": // Package name
				currentPackage["name"] = value
			case "V": // Version
				currentPackage["version"] = value
			case "A": // Architecture
				currentPackage["architecture"] = value
			case "S": // Size
				currentPackage["size"] = value
			case "I": // Installed size
				currentPackage["installed_size"] = value
			case "T": // Description
				currentPackage["description"] = value
			case "U": // URL
				currentPackage["url"] = value
			case "L": // License
				currentPackage["license"] = value
			case "o": // Origin
				currentPackage["origin"] = value
			case "m": // Maintainer
				currentPackage["maintainer"] = value
			case "t": // Install timestamp
				currentPackage["install_time"] = value
			case "c": // Checksum
				currentPackage["checksum"] = value
			case "D": // Dependencies
				deps := currentPackage["dependencies"]
				if deps == "" {
					deps = value
				} else {
					deps += " " + value
				}
				currentPackage["dependencies"] = deps
			case "p": // Provides
				provides := currentPackage["provides"]
				if provides == "" {
					provides = value
				} else {
					provides += " " + value
				}
				currentPackage["provides"] = provides
			}
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

// createPackageFromEntry creates a package model from an APK database entry
func (p *ApkDbParser) createPackageFromEntry(entry map[string]string, filePath string) *model.Package {
	name := entry["name"]
	if name == "" {
		return nil
	}

	version := entry["version"]
	architecture := entry["architecture"]
	description := entry["description"]
	size := entry["size"]
	installedSize := entry["installed_size"]
	url := entry["url"]
	license := entry["license"]
	origin := entry["origin"]
	maintainer := entry["maintainer"]
	dependencies := entry["dependencies"]
	provides := entry["provides"]

	metadata := map[string]interface{}{
		"manager":        "apk",
		"architecture":   architecture,
		"description":    description,
		"size":           size,
		"installed_size": installedSize,
		"url":            url,
		"license":        license,
		"origin":         origin,
		"maintainer":     maintainer,
		"dependencies":   dependencies,
		"provides":       provides,
	}

	cleanVersion := p.cleanApkVersion(version)
	pkg := p.base.CreatePackage(name, cleanVersion, filePath, false, metadata)

	return &pkg
}

// cleanApkVersion cleans and normalizes APK version strings
func (p *ApkDbParser) cleanApkVersion(version string) string {
	version = strings.TrimSpace(version)

	// APK versions can have revision (e.g., "2.3.4-r1")
	if dashIndex := strings.LastIndex(version, "-r"); dashIndex > 0 {
		version = version[:dashIndex]
	}

	return strings.TrimSpace(version)
}
