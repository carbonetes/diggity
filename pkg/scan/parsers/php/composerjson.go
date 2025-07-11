package php

import (
	"encoding/json"
	"io/ioutil"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/scan/parsers/common"
)

// ComposerJsonParser handles parsing of composer.json files
type ComposerJsonParser struct {
	base *common.BaseParser
}

// NewComposerJsonParser creates a new composer.json parser
func NewComposerJsonParser(base *common.BaseParser) *ComposerJsonParser {
	return &ComposerJsonParser{base: base}
}

// ComposerJson represents the structure of a composer.json file
type ComposerJson struct {
	Name       string            `json:"name"`
	Version    string            `json:"version"`
	Require    map[string]string `json:"require"`
	RequireDev map[string]string `json:"require-dev"`
}

// Parse parses composer.json files
func (c *ComposerJsonParser) Parse(filePath string) ([]model.Package, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var composerJson ComposerJson
	if err := json.Unmarshal(content, &composerJson); err != nil {
		return nil, err
	}

	var packages []model.Package

	// Parse production dependencies
	for name, version := range composerJson.Require {
		// Skip PHP itself and other non-package constraints
		if c.isSystemConstraint(name) {
			continue
		}

		metadata := map[string]interface{}{
			"manager": "composer",
		}

		cleanVersion := c.cleanPHPVersion(version)
		pkg := c.base.CreatePackage(name, cleanVersion, filePath, false, metadata)
		packages = append(packages, pkg)
	}

	// Parse development dependencies
	for name, version := range composerJson.RequireDev {
		// Skip PHP itself and other non-package constraints
		if c.isSystemConstraint(name) {
			continue
		}

		metadata := map[string]interface{}{
			"manager": "composer",
		}

		cleanVersion := c.cleanPHPVersion(version)
		pkg := c.base.CreatePackage(name, cleanVersion, filePath, true, metadata)
		packages = append(packages, pkg)
	}

	return packages, nil
}

// isSystemConstraint checks if the dependency is a system constraint (like PHP version)
func (c *ComposerJsonParser) isSystemConstraint(name string) bool {
	systemConstraints := []string{
		"php",
		"ext-",     // PHP extensions
		"lib-",     // System libraries
		"composer", // Composer itself
	}

	for _, constraint := range systemConstraints {
		if name == constraint || (len(name) > len(constraint) && name[:len(constraint)] == constraint) {
			return true
		}
	}

	return false
}

// cleanPHPVersion cleans and normalizes PHP version strings
func (c *ComposerJsonParser) cleanPHPVersion(version string) string {
	prefixes := []string{"^", "~", ">=", "<=", ">", "<", "="}
	return c.base.GetVersionCleaner().CleanVersion(version, prefixes)
}
