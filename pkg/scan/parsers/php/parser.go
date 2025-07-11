package php

import (
	"context"
	"time"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/scan/parsers/common"
	"github.com/carbonetes/diggity/pkg/scan/types"
)

// Parser implements package parsing for PHP projects
type Parser struct {
	*common.BaseParser
	composerJson *ComposerJsonParser
	composerLock *ComposerLockParser
}

// New creates a new PHP parser
func New() *Parser {
	baseParser := common.NewBaseParser("php-parser", "php")

	return &Parser{
		BaseParser:   baseParser,
		composerJson: NewComposerJsonParser(baseParser),
		composerLock: NewComposerLockParser(baseParser),
	}
}

// CheckFile determines if a file is relevant for php parsing
func (p *Parser) CheckFile(path string) (bool, error) {
	phpFiles := []string{
		"composer.json",
		"composer.lock",
	}

	return p.GetFileChecker().CheckFileByName(path, phpFiles), nil
}

// Scan performs php package scanning and parsing
func (p *Parser) Scan(ctx context.Context, target types.ScanTarget) (*types.ScanResult, error) {
	start := time.Now()
	p.LogStart(target)

	result := p.CreateScanResult(target)

	// Find php files in the target directory
	phpFiles, err := p.GetFileWalker().FindFiles(target.Path, p.CheckFile)
	if err != nil {
		p.LogError(result, "Failed to find php files", err)
		return result, nil
	}

	if len(phpFiles) == 0 {
		return result, nil
	}

	// Parse each php file
	var allPackages []model.Package
	for _, file := range phpFiles {
		packages, err := p.parsePHPFile(file)
		if err != nil {
			p.LogError(result, "Failed to parse "+file, err)
			continue
		}
		allPackages = append(allPackages, packages...)
	}

	// Finalize results
	metadata := map[string]interface{}{
		"php_files_processed": len(phpFiles),
	}
	p.FinalizeScanResult(result, allPackages, start, metadata)

	return result, nil
}

// parsePHPFile parses a specific php file using the appropriate sub-parser
func (p *Parser) parsePHPFile(filePath string) ([]model.Package, error) {
	checker := p.GetFileChecker()

	if checker.CheckFileByName(filePath, []string{"composer.json"}) {
		return p.composerJson.Parse(filePath)
	}

	if checker.CheckFileByName(filePath, []string{"composer.lock"}) {
		return p.composerLock.Parse(filePath)
	}

	return []model.Package{}, nil
}
