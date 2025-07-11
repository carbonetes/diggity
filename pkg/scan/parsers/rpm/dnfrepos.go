package rpm

import (
	"bufio"
	"os"
	"strings"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/scan/parsers/common"
)

// DnfReposParser handles parsing of DNF repository files
type DnfReposParser struct {
	base *common.BaseParser
}

// NewDnfReposParser creates a new DNF repos parser
func NewDnfReposParser(base *common.BaseParser) *DnfReposParser {
	return &DnfReposParser{base: base}
}

// Parse parses DNF .repo files (similar to YUM but newer format)
func (p *DnfReposParser) Parse(filePath string) ([]model.Package, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var packages []model.Package
	scanner := bufio.NewScanner(file)

	var currentRepo map[string]string
	var repoName string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip comments and empty lines
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Repository section header [reponame]
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			// Finalize previous repo
			if currentRepo != nil && repoName != "" {
				if pkg := p.createRepoPackage(repoName, currentRepo, filePath); pkg != nil {
					packages = append(packages, *pkg)
				}
			}

			// Start new repo
			repoName = strings.Trim(line, "[]")
			currentRepo = make(map[string]string)
			continue
		}

		// Parse key=value pairs
		if equalIndex := strings.Index(line, "="); equalIndex > 0 && currentRepo != nil {
			key := strings.TrimSpace(line[:equalIndex])
			value := strings.TrimSpace(line[equalIndex+1:])
			currentRepo[key] = value
		}
	}

	// Handle last repository
	if currentRepo != nil && repoName != "" {
		if pkg := p.createRepoPackage(repoName, currentRepo, filePath); pkg != nil {
			packages = append(packages, *pkg)
		}
	}

	return packages, scanner.Err()
}

// createRepoPackage creates a package model from a DNF repository entry
func (p *DnfReposParser) createRepoPackage(name string, repo map[string]string, filePath string) *model.Package {
	if name == "" {
		return nil
	}

	// Only include enabled repositories
	enabled := repo["enabled"]
	if enabled != "1" && enabled != "true" {
		return nil
	}

	baseurl := repo["baseurl"]
	mirrorlist := repo["mirrorlist"]
	metalink := repo["metalink"]
	gpgcheck := repo["gpgcheck"]
	gpgkey := repo["gpgkey"]
	modulesFilter := repo["module_hotfixes"]

	metadata := map[string]interface{}{
		"manager":         "dnf",
		"type":            "repository",
		"baseurl":         baseurl,
		"mirrorlist":      mirrorlist,
		"metalink":        metalink,
		"gpgcheck":        gpgcheck,
		"gpgkey":          gpgkey,
		"enabled":         enabled,
		"module_hotfixes": modulesFilter,
	}

	pkg := p.base.CreatePackage(name, "", filePath, false, metadata)
	return &pkg
}
