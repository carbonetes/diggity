package terraform

import (
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/scan/parsers/common"
)

// LockfileParser handles parsing of .terraform.lock.hcl files
type LockfileParser struct {
	base *common.BaseParser
}

// NewLockfileParser creates a new lockfile parser
func NewLockfileParser(base *common.BaseParser) *LockfileParser {
	return &LockfileParser{base: base}
}

// TerraformProvider represents a provider in the lock file
type TerraformProvider struct {
	Name    string
	Version string
	Source  string
	Hashes  []string
}

// Parse parses .terraform.lock.hcl files
func (p *LockfileParser) Parse(filePath string) ([]model.Package, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var packages []model.Package
	providers := p.parseProviders(string(content))

	for _, provider := range providers {
		if provider.Name == "" || provider.Version == "" {
			continue
		}

		metadata := map[string]interface{}{
			"manager": "terraform",
			"source":  provider.Source,
		}

		if len(provider.Hashes) > 0 {
			metadata["hashes"] = provider.Hashes
		}

		cleanVersion := p.cleanTerraformVersion(provider.Version)
		pkg := p.base.CreatePackage(provider.Name, cleanVersion, filePath, true, metadata)
		packages = append(packages, pkg)
	}

	return packages, nil
}

// parseProviders parses providers from HCL content
func (p *LockfileParser) parseProviders(content string) []TerraformProvider {
	var providers []TerraformProvider

	// Regex to match provider blocks
	providerRegex := regexp.MustCompile(`provider\s+"([^"]+)"\s+\{([^}]+)\}`)
	matches := providerRegex.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if len(match) < 3 {
			continue
		}

		providerSource := match[1]
		providerBlock := match[2]

		// Extract provider name from source (e.g., "registry.terraform.io/hashicorp/aws" -> "aws")
		nameParts := strings.Split(providerSource, "/")
		var providerName string
		if len(nameParts) > 0 {
			providerName = nameParts[len(nameParts)-1]
		}

		// Extract version
		versionRegex := regexp.MustCompile(`version\s*=\s*"([^"]+)"`)
		versionMatch := versionRegex.FindStringSubmatch(providerBlock)
		var version string
		if len(versionMatch) > 1 {
			version = versionMatch[1]
		}

		// Extract hashes
		var hashes []string
		hashRegex := regexp.MustCompile(`"(h1:[^"]+)"`)
		hashMatches := hashRegex.FindAllStringSubmatch(providerBlock, -1)
		for _, hashMatch := range hashMatches {
			if len(hashMatch) > 1 {
				hashes = append(hashes, hashMatch[1])
			}
		}

		if providerName != "" && version != "" {
			providers = append(providers, TerraformProvider{
				Name:    providerName,
				Version: version,
				Source:  providerSource,
				Hashes:  hashes,
			})
		}
	}

	return providers
}

// cleanTerraformVersion cleans and normalizes Terraform version strings
func (p *LockfileParser) cleanTerraformVersion(version string) string {
	version = strings.TrimSpace(version)
	prefixes := []string{"~>", ">=", "<=", ">", "<", "="}
	return p.base.GetVersionCleaner().CleanVersion(version, prefixes)
}
