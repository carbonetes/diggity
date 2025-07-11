package golang

import (
	"bufio"
	"os"
	"strings"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/scan/parsers/common"
)

// GoSumParser handles parsing of go.sum files
type GoSumParser struct {
	base *common.BaseParser
}

// NewGoSumParser creates a new go.sum parser
func NewGoSumParser(base *common.BaseParser) *GoSumParser {
	return &GoSumParser{base: base}
}

// Parse parses go.sum files
func (g *GoSumParser) Parse(filePath string) ([]model.Package, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var packages []model.Package
	scanner := bufio.NewScanner(file)

	seen := make(map[string]bool) // Track already seen packages

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines
		if line == "" {
			continue
		}

		pkg := g.parseSumLine(line, filePath)
		if pkg != nil {
			// Avoid duplicates (go.sum often has multiple entries per module)
			key := pkg.Name + "@" + pkg.Version
			if !seen[key] {
				packages = append(packages, *pkg)
				seen[key] = true
			}
		}
	}

	return packages, scanner.Err()
}

// parseSumLine parses a line from go.sum
func (g *GoSumParser) parseSumLine(line, filePath string) *model.Package {
	// go.sum format: module version hash
	parts := strings.Fields(line)
	if len(parts) < 3 {
		return nil
	}

	modulePath := parts[0]
	version := parts[1]
	hash := parts[2]

	// Skip /go.mod entries (these are module metadata, not packages)
	if strings.HasSuffix(version, "/go.mod") {
		return nil
	}

	metadata := map[string]interface{}{
		"manager": "go",
		"hash":    hash,
		"locked":  true,
	}

	pkg := g.base.CreatePackage(modulePath, version, filePath, false, metadata)
	return &pkg
}
