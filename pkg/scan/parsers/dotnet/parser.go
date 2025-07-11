package dotnet

import (
	"context"
	"time"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/scan/parsers/common"
	"github.com/carbonetes/diggity/pkg/scan/types"
)

// Parser implements package parsing for .NET projects
type Parser struct {
	*common.BaseParser
	packageConfig *PackageConfigParser
	packagesLock  *PackagesLockParser
	projectFile   *ProjectFileParser
}

// New creates a new .NET parser
func New() *Parser {
	baseParser := common.NewBaseParser("dotnet-parser", "dotnet")

	return &Parser{
		BaseParser:    baseParser,
		packageConfig: NewPackageConfigParser(baseParser),
		packagesLock:  NewPackagesLockParser(baseParser),
		projectFile:   NewProjectFileParser(baseParser),
	}
}

// CheckFile determines if a file is relevant for .NET parsing
func (p *Parser) CheckFile(path string) (bool, error) {
	dotnetFiles := []string{
		"packages.config",
		"packages.lock.json",
	}

	// Check for project files
	if p.GetFileChecker().CheckFileByExtension(path, []string{".csproj", ".vbproj", ".fsproj"}) {
		return true, nil
	}

	return p.GetFileChecker().CheckFileByName(path, dotnetFiles), nil
}

// Scan performs .NET package scanning and parsing
func (p *Parser) Scan(ctx context.Context, target types.ScanTarget) (*types.ScanResult, error) {
	start := time.Now()
	p.LogStart(target)

	result := p.CreateScanResult(target)

	// Find .NET files in the target directory
	dotnetFiles, err := p.GetFileWalker().FindFiles(target.Path, p.CheckFile)
	if err != nil {
		p.LogError(result, "Failed to find .NET files", err)
		return result, nil
	}

	if len(dotnetFiles) == 0 {
		return result, nil
	}

	// Parse each .NET file
	var allPackages []model.Package
	for _, file := range dotnetFiles {
		packages, err := p.parseDotNetFile(file)
		if err != nil {
			p.LogError(result, "Failed to parse "+file, err)
			continue
		}
		allPackages = append(allPackages, packages...)
	}

	// Finalize results
	metadata := map[string]interface{}{
		"dotnet_files_processed": len(dotnetFiles),
	}
	p.FinalizeScanResult(result, allPackages, start, metadata)

	return result, nil
}

// parseDotNetFile parses a specific .NET file using the appropriate sub-parser
func (p *Parser) parseDotNetFile(filePath string) ([]model.Package, error) {
	checker := p.GetFileChecker()

	if checker.CheckFileByName(filePath, []string{"packages.config"}) {
		return p.packageConfig.Parse(filePath)
	}

	if checker.CheckFileByName(filePath, []string{"packages.lock.json"}) {
		return p.packagesLock.Parse(filePath)
	}

	if checker.CheckFileByExtension(filePath, []string{".csproj", ".vbproj", ".fsproj"}) {
		return p.projectFile.Parse(filePath)
	}

	return []model.Package{}, nil
}
