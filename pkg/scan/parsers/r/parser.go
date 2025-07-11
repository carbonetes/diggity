package r

import (
	"context"
	"path/filepath"
	"time"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/scan/parsers/common"
	"github.com/carbonetes/diggity/pkg/scan/types"
)

// Parser implements package parsing for R projects
type Parser struct {
	*common.BaseParser
	description *DescriptionParser
	renv        *RenvParser
}

// New creates a new R parser
func New() *Parser {
	baseParser := common.NewBaseParser("r-parser", "r")

	return &Parser{
		BaseParser:  baseParser,
		description: NewDescriptionParser(baseParser),
		renv:        NewRenvParser(baseParser),
	}
}

// CheckFile checks if a file should be scanned by this parser
func (p *Parser) CheckFile(filePath string) (bool, error) {
	fileName := filepath.Base(filePath)

	rFiles := []string{
		"DESCRIPTION",
		"renv.lock",
	}

	return p.GetFileChecker().CheckFileByName(fileName, rFiles), nil
}

// Scan scans the target for R packages
func (p *Parser) Scan(ctx context.Context, target types.ScanTarget) (*types.ScanResult, error) {
	start := time.Now()
	result := p.CreateScanResult(target)
	var allPackages []model.Package

	// Find all relevant R files
	files, err := p.GetFileWalker().FindFiles(target.Path, p.CheckFile)
	if err != nil {
		return result, err
	}

	for _, file := range files {
		var packages []model.Package
		var parseErr error

		fileName := filepath.Base(file)

		switch fileName {
		case "DESCRIPTION":
			packages, parseErr = p.description.Parse(file)
		case "renv.lock":
			packages, parseErr = p.renv.Parse(file)
		}

		if parseErr != nil {
			// Log error but continue processing other files
			continue
		}

		allPackages = append(allPackages, packages...)
	}

	metadata := map[string]interface{}{
		"r_files_processed": len(files),
	}
	p.FinalizeScanResult(result, allPackages, start, metadata)
	return result, nil
}
