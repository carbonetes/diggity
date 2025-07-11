package swift

import (
	"context"
	"path/filepath"
	"time"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/scan/parsers/common"
	"github.com/carbonetes/diggity/pkg/scan/types"
)

// Parser implements package parsing for Swift projects
type Parser struct {
	*common.BaseParser
	packageSwift    *PackageSwiftParser
	packageResolved *PackageResolvedParser
	podfile         *PodfileParser
}

// New creates a new Swift parser
func New() *Parser {
	baseParser := common.NewBaseParser("swift-parser", "swift")

	return &Parser{
		BaseParser:      baseParser,
		packageSwift:    NewPackageSwiftParser(baseParser),
		packageResolved: NewPackageResolvedParser(baseParser),
		podfile:         NewPodfileParser(baseParser),
	}
}

// CheckFile checks if a file should be scanned by this parser
func (p *Parser) CheckFile(filePath string) (bool, error) {
	fileName := filepath.Base(filePath)

	swiftFiles := []string{
		"Package.swift",
		"Package.resolved",
		".package.resolved",
		"Podfile",
		"Podfile.lock",
	}

	return p.GetFileChecker().CheckFileByName(fileName, swiftFiles), nil
}

// Scan scans the target for Swift packages
func (p *Parser) Scan(ctx context.Context, target types.ScanTarget) (*types.ScanResult, error) {
	start := time.Now()
	result := p.CreateScanResult(target)
	var allPackages []model.Package

	// Find all relevant Swift files
	files, err := p.GetFileWalker().FindFiles(target.Path, p.CheckFile)
	if err != nil {
		return result, err
	}

	for _, file := range files {
		var packages []model.Package
		var parseErr error

		fileName := filepath.Base(file)

		switch fileName {
		case "Package.swift":
			packages, parseErr = p.packageSwift.Parse(file)
		case "Package.resolved", ".package.resolved":
			packages, parseErr = p.packageResolved.Parse(file)
		case "Podfile":
			packages, parseErr = p.podfile.Parse(file)
		case "Podfile.lock":
			packages, parseErr = p.podfile.ParseLock(file)
		}

		if parseErr != nil {
			// Log error but continue processing other files
			continue
		}

		allPackages = append(allPackages, packages...)
	}

	metadata := map[string]interface{}{
		"swift_files_processed": len(files),
	}
	p.FinalizeScanResult(result, allPackages, start, metadata)
	return result, nil
}
