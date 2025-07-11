package golang

import (
	"context"
	"time"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/scan/parsers/common"
	"github.com/carbonetes/diggity/pkg/scan/types"
)

// Parser implements package parsing for Go projects
type Parser struct {
	*common.BaseParser
	goMod *GoModParser
	goSum *GoSumParser
}

// New creates a new Go parser
func New() *Parser {
	baseParser := common.NewBaseParser("go-parser", "go")

	return &Parser{
		BaseParser: baseParser,
		goMod:      NewGoModParser(baseParser),
		goSum:      NewGoSumParser(baseParser),
	}
}

// CheckFile determines if a file is relevant for go parsing
func (p *Parser) CheckFile(path string) (bool, error) {
	goFiles := []string{
		"go.mod",
		"go.sum",
	}

	return p.GetFileChecker().CheckFileByName(path, goFiles), nil
}

// Scan performs go package scanning and parsing
func (p *Parser) Scan(ctx context.Context, target types.ScanTarget) (*types.ScanResult, error) {
	start := time.Now()
	p.LogStart(target)

	result := p.CreateScanResult(target)

	// Find go files in the target directory
	goFiles, err := p.GetFileWalker().FindFiles(target.Path, p.CheckFile)
	if err != nil {
		p.LogError(result, "Failed to find go files", err)
		return result, nil
	}

	if len(goFiles) == 0 {
		return result, nil
	}

	// Parse each go file
	var allPackages []model.Package
	for _, file := range goFiles {
		packages, err := p.parseGoFile(file)
		if err != nil {
			p.LogError(result, "Failed to parse "+file, err)
			continue
		}
		allPackages = append(allPackages, packages...)
	}

	// Finalize results
	metadata := map[string]interface{}{
		"go_files_processed": len(goFiles),
	}
	p.FinalizeScanResult(result, allPackages, start, metadata)

	return result, nil
}

// parseGoFile parses a specific go file using the appropriate sub-parser
func (p *Parser) parseGoFile(filePath string) ([]model.Package, error) {
	checker := p.GetFileChecker()

	if checker.CheckFileByName(filePath, []string{"go.mod"}) {
		return p.goMod.Parse(filePath)
	}

	if checker.CheckFileByName(filePath, []string{"go.sum"}) {
		return p.goSum.Parse(filePath)
	}

	return []model.Package{}, nil
}
