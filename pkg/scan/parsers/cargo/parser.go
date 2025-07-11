package cargo

import (
	"context"
	"time"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/scan/parsers/common"
	"github.com/carbonetes/diggity/pkg/scan/types"
)

// Parser implements package parsing for Rust/Cargo projects
type Parser struct {
	*common.BaseParser
	cargoToml *CargoTomlParser
	cargoLock *CargoLockParser
}

// New creates a new Cargo parser
func New() *Parser {
	baseParser := common.NewBaseParser("cargo-parser", "cargo")

	return &Parser{
		BaseParser: baseParser,
		cargoToml:  NewCargoTomlParser(baseParser),
		cargoLock:  NewCargoLockParser(baseParser),
	}
}

// CheckFile determines if a file is relevant for cargo parsing
func (p *Parser) CheckFile(path string) (bool, error) {
	cargoFiles := []string{
		"Cargo.toml",
		"Cargo.lock",
	}

	return p.GetFileChecker().CheckFileByName(path, cargoFiles), nil
}

// Scan performs cargo package scanning and parsing
func (p *Parser) Scan(ctx context.Context, target types.ScanTarget) (*types.ScanResult, error) {
	start := time.Now()
	p.LogStart(target)

	result := p.CreateScanResult(target)

	// Find cargo files in the target directory
	cargoFiles, err := p.GetFileWalker().FindFiles(target.Path, p.CheckFile)
	if err != nil {
		p.LogError(result, "Failed to find cargo files", err)
		return result, nil
	}

	if len(cargoFiles) == 0 {
		return result, nil
	}

	// Parse each cargo file
	var allPackages []model.Package
	for _, file := range cargoFiles {
		packages, err := p.parseCargoFile(file)
		if err != nil {
			p.LogError(result, "Failed to parse "+file, err)
			continue
		}
		allPackages = append(allPackages, packages...)
	}

	// Finalize results
	metadata := map[string]interface{}{
		"cargo_files_processed": len(cargoFiles),
	}
	p.FinalizeScanResult(result, allPackages, start, metadata)

	return result, nil
}

// parseCargoFile parses a specific cargo file using the appropriate sub-parser
func (p *Parser) parseCargoFile(filePath string) ([]model.Package, error) {
	checker := p.GetFileChecker()

	if checker.CheckFileByName(filePath, []string{"Cargo.toml"}) {
		return p.cargoToml.Parse(filePath)
	}

	if checker.CheckFileByName(filePath, []string{"Cargo.lock"}) {
		return p.cargoLock.Parse(filePath)
	}

	return []model.Package{}, nil
}
