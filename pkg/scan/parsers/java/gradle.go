package java

import (
	"bufio"
	"os"
	"regexp"
	"strings"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/scan/parsers/common"
)

// GradleParser handles parsing of build.gradle files
type GradleParser struct {
	base *common.BaseParser
}

// NewGradleParser creates a new Gradle parser
func NewGradleParser(base *common.BaseParser) *GradleParser {
	return &GradleParser{base: base}
}

// Parse parses build.gradle files
func (g *GradleParser) Parse(filePath string) ([]model.Package, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var packages []model.Package
	scanner := bufio.NewScanner(file)

	parser := &gradleBlockParser{}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if g.shouldSkipLine(line) {
			continue
		}

		parser.processLine(line)

		if parser.inDependenciesBlock() {
			pkg := g.parseDependencyLine(line, filePath)
			if pkg != nil {
				packages = append(packages, *pkg)
			}
		}
	}

	return packages, scanner.Err()
}

// gradleBlockParser tracks parsing state for Gradle files
type gradleBlockParser struct {
	inDepsBlock bool
	blockDepth  int
}

// processLine processes a line and updates parsing state
func (p *gradleBlockParser) processLine(line string) {
	if strings.Contains(line, "dependencies") && strings.Contains(line, "{") {
		p.inDepsBlock = true
		p.blockDepth = 1
		return
	}

	if p.inDepsBlock {
		p.blockDepth += strings.Count(line, "{")
		p.blockDepth -= strings.Count(line, "}")

		if p.blockDepth <= 0 {
			p.inDepsBlock = false
		}
	}
}

// inDependenciesBlock returns true if currently in dependencies block
func (p *gradleBlockParser) inDependenciesBlock() bool {
	return p.inDepsBlock && p.blockDepth > 0
}

// shouldSkipLine checks if a line should be skipped
func (g *GradleParser) shouldSkipLine(line string) bool {
	return line == "" || strings.HasPrefix(line, "//") || strings.HasPrefix(line, "*")
}

// parseDependencyLine parses a single dependency line from build.gradle
func (g *GradleParser) parseDependencyLine(line, filePath string) *model.Package {
	// Common Gradle dependency patterns
	patterns := []string{
		// implementation 'group:artifact:version'
		`(?:implementation|api|compile|testImplementation|testCompile|runtimeOnly|compileOnly)\s+['"]([\w\.-]+):([\w\.-]+):([^'"]+)['"]`,
		// implementation group: 'group', name: 'artifact', version: 'version'
		`(?:implementation|api|compile|testImplementation|testCompile|runtimeOnly|compileOnly)\s+group:\s*['"]([^'"]+)['"],\s*name:\s*['"]([^'"]+)['"],\s*version:\s*['"]([^'"]+)['"]`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(line)

		if len(matches) >= 4 {
			groupID := matches[1]
			artifactID := matches[2]
			version := matches[3]

			// Create package name as groupId:artifactId
			name := groupID + ":" + artifactID

			// Check if it's a dev dependency
			isDev := g.isDevDependency(line)

			metadata := map[string]interface{}{
				"manager":     "gradle",
				"group_id":    groupID,
				"artifact_id": artifactID,
				"scope":       g.extractScope(line),
			}

			pkg := g.base.CreatePackage(name, version, filePath, isDev, metadata)
			return &pkg
		}
	}

	return nil
}

// isDevDependency checks if the dependency is for development/testing
func (g *GradleParser) isDevDependency(line string) bool {
	devKeywords := []string{"test", "Test", "debug", "Debug"}
	for _, keyword := range devKeywords {
		if strings.Contains(line, keyword) {
			return true
		}
	}
	return false
}

// extractScope extracts the scope/configuration from the dependency line
func (g *GradleParser) extractScope(line string) string {
	scopes := []string{"implementation", "api", "compile", "testImplementation", "testCompile", "runtimeOnly", "compileOnly"}
	for _, scope := range scopes {
		if strings.Contains(line, scope) {
			return scope
		}
	}
	return "implementation"
}
