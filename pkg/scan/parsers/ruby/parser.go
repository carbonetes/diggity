package ruby

import (
	"context"
	"time"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/scan/parsers/common"
	"github.com/carbonetes/diggity/pkg/scan/types"
)

// Parser implements package parsing for Ruby projects
type Parser struct {
	*common.BaseParser
	gemfile     *GemfileParser
	gemfileLock *GemfileLockParser
	gemspec     *GemspecParser
}

// New creates a new Ruby parser
func New() *Parser {
	baseParser := common.NewBaseParser("ruby-parser", "ruby")

	return &Parser{
		BaseParser:  baseParser,
		gemfile:     NewGemfileParser(baseParser),
		gemfileLock: NewGemfileLockParser(baseParser),
		gemspec:     NewGemspecParser(baseParser),
	}
}

// CheckFile determines if a file is relevant for ruby parsing
func (p *Parser) CheckFile(path string) (bool, error) {
	rubyFiles := []string{
		"Gemfile",
		"Gemfile.lock",
	}

	// Check for .gemspec files
	if p.GetFileChecker().CheckFileByExtension(path, []string{".gemspec"}) {
		return true, nil
	}

	return p.GetFileChecker().CheckFileByName(path, rubyFiles), nil
}

// Scan performs ruby package scanning and parsing
func (p *Parser) Scan(ctx context.Context, target types.ScanTarget) (*types.ScanResult, error) {
	start := time.Now()
	p.LogStart(target)

	result := p.CreateScanResult(target)

	// Find ruby files in the target directory
	rubyFiles, err := p.GetFileWalker().FindFiles(target.Path, p.CheckFile)
	if err != nil {
		p.LogError(result, "Failed to find ruby files", err)
		return result, nil
	}

	if len(rubyFiles) == 0 {
		return result, nil
	}

	// Parse each ruby file
	var allPackages []model.Package
	for _, file := range rubyFiles {
		packages, err := p.parseRubyFile(file)
		if err != nil {
			p.LogError(result, "Failed to parse "+file, err)
			continue
		}
		allPackages = append(allPackages, packages...)
	}

	// Finalize results
	metadata := map[string]interface{}{
		"ruby_files_processed": len(rubyFiles),
	}
	p.FinalizeScanResult(result, allPackages, start, metadata)

	return result, nil
}

// parseRubyFile parses a specific ruby file using the appropriate sub-parser
func (p *Parser) parseRubyFile(filePath string) ([]model.Package, error) {
	checker := p.GetFileChecker()

	if checker.CheckFileByName(filePath, []string{"Gemfile"}) {
		return p.gemfile.Parse(filePath)
	}

	if checker.CheckFileByName(filePath, []string{"Gemfile.lock"}) {
		return p.gemfileLock.Parse(filePath)
	}

	if checker.CheckFileByExtension(filePath, []string{".gemspec"}) {
		return p.gemspec.Parse(filePath)
	}

	return []model.Package{}, nil
}
