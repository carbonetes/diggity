package swift

import (
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/scan/parsers/common"
)

// PackageSwiftParser handles parsing of Package.swift files
type PackageSwiftParser struct {
	base *common.BaseParser
}

// NewPackageSwiftParser creates a new Package.swift parser
func NewPackageSwiftParser(base *common.BaseParser) *PackageSwiftParser {
	return &PackageSwiftParser{base: base}
}

// Parse parses Package.swift files (basic dependency extraction)
func (p *PackageSwiftParser) Parse(filePath string) ([]model.Package, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var packages []model.Package

	// Simple regex to find .package dependencies in Package.swift
	// This is a simplified approach - full Swift parsing would be more complex
	packageRegex := regexp.MustCompile(`\.package\s*\(\s*url:\s*"([^"]+)"[^)]*\)`)
	matches := packageRegex.FindAllStringSubmatch(string(content), -1)

	for _, match := range matches {
		if len(match) > 1 {
			url := match[1]

			// Extract package name from URL
			name := p.extractPackageNameFromURL(url)
			if name == "" {
				continue
			}

			metadata := map[string]interface{}{
				"manager": "swift-package-manager",
				"url":     url,
			}

			pkg := p.base.CreatePackage(name, "", filePath, false, metadata)
			packages = append(packages, pkg)
		}
	}

	return packages, nil
}

// extractPackageNameFromURL extracts package name from a Git URL
func (p *PackageSwiftParser) extractPackageNameFromURL(url string) string {
	// Remove .git suffix
	url = strings.TrimSuffix(url, ".git")

	// Extract the last part of the path
	parts := strings.Split(url, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}

	return ""
}
