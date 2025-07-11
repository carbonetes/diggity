package r

import (
	"bufio"
	"os"
	"regexp"
	"strings"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/scan/parsers/common"
)

// DescriptionParser handles parsing of DESCRIPTION files
type DescriptionParser struct {
	base *common.BaseParser
}

// NewDescriptionParser creates a new DESCRIPTION parser
func NewDescriptionParser(base *common.BaseParser) *DescriptionParser {
	return &DescriptionParser{base: base}
}

// RDependency represents a parsed R dependency
type RDependency struct {
	Name    string
	Version string
}

// Parse parses DESCRIPTION files
func (p *DescriptionParser) Parse(filePath string) ([]model.Package, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var packages []model.Package
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Parse Imports, Depends, Suggests, LinkingTo
		if strings.HasPrefix(line, "Imports:") ||
			strings.HasPrefix(line, "Depends:") ||
			strings.HasPrefix(line, "Suggests:") ||
			strings.HasPrefix(line, "LinkingTo:") {

			dependencyType := strings.TrimSuffix(strings.Split(line, ":")[0], "")
			dependencyList := strings.TrimSpace(strings.SplitN(line, ":", 2)[1])

			// Handle multi-line dependencies
			for scanner.Scan() {
				nextLine := scanner.Text()
				if strings.HasPrefix(nextLine, " ") || strings.HasPrefix(nextLine, "\t") {
					dependencyList += " " + strings.TrimSpace(nextLine)
				} else {
					// Put the line back by creating a new scanner
					break
				}
			}

			deps := p.parseDependencyList(dependencyList)
			for _, dep := range deps {
				if dep.Name == "R" {
					// Skip R base dependency
					continue
				}

				metadata := map[string]interface{}{
					"manager":         "cran",
					"dependency_type": strings.ToLower(dependencyType),
				}

				cleanVersion := p.cleanRVersion(dep.Version)
				pkg := p.base.CreatePackage(dep.Name, cleanVersion, filePath, false, metadata)
				packages = append(packages, pkg)
			}
		}
	}

	return packages, nil
}

// parseDependencyList parses a comma-separated list of R dependencies
func (p *DescriptionParser) parseDependencyList(depList string) []RDependency {
	var dependencies []RDependency

	// Split by comma, but be careful of parentheses
	deps := strings.Split(depList, ",")

	for _, dep := range deps {
		dep = strings.TrimSpace(dep)
		if dep == "" {
			continue
		}

		// Parse package name and version constraint
		// Format: packagename (>= version)
		re := regexp.MustCompile(`^([^\s\(]+)(?:\s*\(([^)]+)\))?`)
		matches := re.FindStringSubmatch(dep)

		if len(matches) > 1 {
			name := matches[1]
			version := ""
			if len(matches) > 2 {
				version = matches[2]
			}

			dependencies = append(dependencies, RDependency{
				Name:    name,
				Version: version,
			})
		}
	}

	return dependencies
}

// cleanRVersion cleans and normalizes R version strings
func (p *DescriptionParser) cleanRVersion(version string) string {
	version = strings.TrimSpace(version)
	prefixes := []string{">=", "<=", ">", "<", "=", "=="}
	return p.base.GetVersionCleaner().CleanVersion(version, prefixes)
}
