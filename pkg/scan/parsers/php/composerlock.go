package php

import (
	"encoding/json"
	"io/ioutil"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/scan/parsers/common"
)

// ComposerLockParser handles parsing of composer.lock files
type ComposerLockParser struct {
	base *common.BaseParser
}

// NewComposerLockParser creates a new composer.lock parser
func NewComposerLockParser(base *common.BaseParser) *ComposerLockParser {
	return &ComposerLockParser{base: base}
}

// ComposerLock represents the structure of a composer.lock file
type ComposerLock struct {
	Packages    []ComposerPackage `json:"packages"`
	PackagesDev []ComposerPackage `json:"packages-dev"`
}

// ComposerPackage represents a package in composer.lock
type ComposerPackage struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Type    string `json:"type"`
	Source  Source `json:"source"`
	Dist    Dist   `json:"dist"`
}

// Source represents source information
type Source struct {
	Type      string `json:"type"`
	URL       string `json:"url"`
	Reference string `json:"reference"`
}

// Dist represents distribution information
type Dist struct {
	Type      string `json:"type"`
	URL       string `json:"url"`
	Reference string `json:"reference"`
	Shasum    string `json:"shasum"`
}

// Parse parses composer.lock files
func (c *ComposerLockParser) Parse(filePath string) ([]model.Package, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var composerLock ComposerLock
	if err := json.Unmarshal(content, &composerLock); err != nil {
		return nil, err
	}

	var packages []model.Package

	// Parse production packages
	for _, pkg := range composerLock.Packages {
		metadata := map[string]interface{}{
			"manager": "composer",
			"locked":  true,
			"type":    pkg.Type,
		}

		if pkg.Source.Reference != "" {
			metadata["source_reference"] = pkg.Source.Reference
		}

		if pkg.Dist.Shasum != "" {
			metadata["shasum"] = pkg.Dist.Shasum
		}

		cleanVersion := c.cleanPHPVersion(pkg.Version)
		modelPkg := c.base.CreatePackage(pkg.Name, cleanVersion, filePath, false, metadata)
		packages = append(packages, modelPkg)
	}

	// Parse development packages
	for _, pkg := range composerLock.PackagesDev {
		metadata := map[string]interface{}{
			"manager": "composer",
			"locked":  true,
			"type":    pkg.Type,
		}

		if pkg.Source.Reference != "" {
			metadata["source_reference"] = pkg.Source.Reference
		}

		if pkg.Dist.Shasum != "" {
			metadata["shasum"] = pkg.Dist.Shasum
		}

		cleanVersion := c.cleanPHPVersion(pkg.Version)
		modelPkg := c.base.CreatePackage(pkg.Name, cleanVersion, filePath, true, metadata)
		packages = append(packages, modelPkg)
	}

	return packages, nil
}

// cleanPHPVersion cleans and normalizes PHP version strings
func (c *ComposerLockParser) cleanPHPVersion(version string) string {
	prefixes := []string{"v", "^", "~", ">=", "<=", ">", "<", "="}
	return c.base.GetVersionCleaner().CleanVersion(version, prefixes)
}
