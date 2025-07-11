package npm

import (
	"context"
	"time"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/scan/parsers/common"
	"github.com/carbonetes/diggity/pkg/scan/types"
)

// Parser implements package parsing for Node.js/npm projects
type Parser struct {
	*common.BaseParser
	packageJson *PackageJsonParser
	packageLock *PackageLockParser
	yarnLock    *YarnLockParser
}

// New creates a new NPM parser
func New() *Parser {
	baseParser := common.NewBaseParser("npm-parser", "npm")

	return &Parser{
		BaseParser:  baseParser,
		packageJson: NewPackageJsonParser(baseParser),
		packageLock: NewPackageLockParser(baseParser),
		yarnLock:    NewYarnLockParser(baseParser),
	}
}

// CheckFile determines if a file is relevant for npm parsing
func (p *Parser) CheckFile(path string) (bool, error) {
	npmFiles := []string{
		"package.json",
		"package-lock.json",
		"yarn.lock",
	}

	return p.GetFileChecker().CheckFileByName(path, npmFiles), nil
}

// Scan performs npm package scanning and parsing
func (p *Parser) Scan(ctx context.Context, target types.ScanTarget) (*types.ScanResult, error) {
	start := time.Now()
	p.LogStart(target)

	result := p.CreateScanResult(target)

	// Find npm files in the target directory
	npmFiles, err := p.GetFileWalker().FindFiles(target.Path, p.CheckFile)
	if err != nil {
		p.LogError(result, "Failed to find npm files", err)
		return result, nil
	}

	if len(npmFiles) == 0 {
		return result, nil
	}

	// Parse each npm file
	var allPackages []model.Package
	for _, file := range npmFiles {
		packages, err := p.parseNPMFile(file)
		if err != nil {
			p.LogError(result, "Failed to parse "+file, err)
			continue
		}
		allPackages = append(allPackages, packages...)
	}

	// Finalize results
	metadata := map[string]interface{}{
		"npm_files_processed": len(npmFiles),
	}
	p.FinalizeScanResult(result, allPackages, start, metadata)

	return result, nil
}

// parseNPMFile parses a specific npm file using the appropriate sub-parser
func (p *Parser) parseNPMFile(filePath string) ([]model.Package, error) {
	checker := p.GetFileChecker()

	if checker.CheckFileByName(filePath, []string{"package.json"}) {
		return p.packageJson.Parse(filePath)
	}

	if checker.CheckFileByName(filePath, []string{"package-lock.json"}) {
		return p.packageLock.Parse(filePath)
	}

	if checker.CheckFileByName(filePath, []string{"yarn.lock"}) {
		return p.yarnLock.Parse(filePath)
	}

	return []model.Package{}, nil
}
