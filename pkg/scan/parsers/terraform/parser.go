package terraform

import (
	"context"
	"path/filepath"
	"time"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/scan/parsers/common"
	"github.com/carbonetes/diggity/pkg/scan/types"
)

// Parser implements package parsing for Terraform projects
type Parser struct {
	*common.BaseParser
	lockfile *LockfileParser
}

// New creates a new Terraform parser
func New() *Parser {
	baseParser := common.NewBaseParser("terraform-parser", "terraform")

	return &Parser{
		BaseParser: baseParser,
		lockfile:   NewLockfileParser(baseParser),
	}
}

// CheckFile checks if a file should be scanned by this parser
func (p *Parser) CheckFile(filePath string) (bool, error) {
	fileName := filepath.Base(filePath)

	terraformFiles := []string{
		".terraform.lock.hcl",
		"terraform.lock.hcl",
	}

	return p.GetFileChecker().CheckFileByName(fileName, terraformFiles), nil
}

// Scan scans the target for Terraform providers
func (p *Parser) Scan(ctx context.Context, target types.ScanTarget) (*types.ScanResult, error) {
	start := time.Now()
	result := p.CreateScanResult(target)
	var allPackages []model.Package

	// Find all relevant Terraform files
	files, err := p.GetFileWalker().FindFiles(target.Path, p.CheckFile)
	if err != nil {
		return result, err
	}

	for _, file := range files {
		var packages []model.Package
		var parseErr error

		fileName := filepath.Base(file)
		if fileName == ".terraform.lock.hcl" || fileName == "terraform.lock.hcl" {
			packages, parseErr = p.lockfile.Parse(file)
		}

		if parseErr != nil {
			// Log error but continue processing other files
			continue
		}

		allPackages = append(allPackages, packages...)
	}

	metadata := map[string]interface{}{
		"terraform_files_processed": len(files),
	}
	p.FinalizeScanResult(result, allPackages, start, metadata)
	return result, nil
}
