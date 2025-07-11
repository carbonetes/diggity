package r

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/scan/parsers/common"
)

// RenvParser handles parsing of renv.lock files
type RenvParser struct {
	base *common.BaseParser
}

// NewRenvParser creates a new renv.lock parser
func NewRenvParser(base *common.BaseParser) *RenvParser {
	return &RenvParser{base: base}
}

// RenvLock represents the structure of renv.lock
type RenvLock struct {
	R        RVersion               `json:"R"`
	Packages map[string]RenvPackage `json:"Packages"`
}

// RVersion represents R version info
type RVersion struct {
	Version string `json:"Version"`
}

// RenvPackage represents a package in renv.lock
type RenvPackage struct {
	Package    string `json:"Package"`
	Version    string `json:"Version"`
	Source     string `json:"Source"`
	Repository string `json:"Repository,omitempty"`
}

// Parse parses renv.lock files
func (p *RenvParser) Parse(filePath string) ([]model.Package, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var lockFile RenvLock
	if err := json.Unmarshal(content, &lockFile); err != nil {
		return nil, err
	}

	var packages []model.Package

	for name, pkg := range lockFile.Packages {
		if name == "renv" {
			// Skip renv itself
			continue
		}

		metadata := map[string]interface{}{
			"manager":    "cran",
			"source":     pkg.Source,
			"repository": pkg.Repository,
		}

		cleanVersion := p.cleanRVersion(pkg.Version)
		modelPkg := p.base.CreatePackage(pkg.Package, cleanVersion, filePath, true, metadata)
		packages = append(packages, modelPkg)
	}

	return packages, nil
}

// cleanRVersion cleans and normalizes R version strings
func (p *RenvParser) cleanRVersion(version string) string {
	version = strings.TrimSpace(version)
	prefixes := []string{">=", "<=", ">", "<", "=", "=="}
	return p.base.GetVersionCleaner().CleanVersion(version, prefixes)
}
