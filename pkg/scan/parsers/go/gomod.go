package golang

import (
	"bufio"
	"os"
	"regexp"
	"strings"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/scan/parsers/common"
)

// GoModParser handles parsing of go.mod files
type GoModParser struct {
	base *common.BaseParser
}

// NewGoModParser creates a new go.mod parser
func NewGoModParser(base *common.BaseParser) *GoModParser {
	return &GoModParser{base: base}
}

// Parse parses go.mod files
func (g *GoModParser) Parse(filePath string) ([]model.Package, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var packages []model.Package
	scanner := bufio.NewScanner(file)

	inRequireBlock := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}

		// Check for require block
		if strings.HasPrefix(line, "require (") {
			inRequireBlock = true
			continue
		}

		if inRequireBlock && line == ")" {
			inRequireBlock = false
			continue
		}

		// Parse dependency line
		pkg := g.parseDependencyLine(line, filePath, inRequireBlock)
		if pkg != nil {
			packages = append(packages, *pkg)
		}
	}

	return packages, scanner.Err()
}

// parseDependencyLine parses a dependency line from go.mod
func (g *GoModParser) parseDependencyLine(line, filePath string, inRequireBlock bool) *model.Package {
	// Remove "require " prefix if not in block
	if !inRequireBlock && strings.HasPrefix(line, "require ") {
		line = strings.TrimPrefix(line, "require ")
	}

	// Pattern: module_path version [// indirect]
	pattern := `^([^\s]+)\s+([^\s]+)(?:\s+//\s*(indirect))?`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(line)

	if len(matches) >= 3 {
		modulePath := matches[1]
		version := matches[2]
		isIndirect := len(matches) > 3 && matches[3] == "indirect"

		metadata := map[string]interface{}{
			"manager":  "go",
			"indirect": isIndirect,
		}

		pkg := g.base.CreatePackage(modulePath, version, filePath, false, metadata)
		return &pkg
	}

	return nil
}
