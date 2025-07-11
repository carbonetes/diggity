package dart

import (
	"context"
	"path/filepath"
	"time"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/scan/parsers/common"
	"github.com/carbonetes/diggity/pkg/scan/types"
)

// Parser implements package parsing for Dart/Flutter projects
type Parser struct {
	*common.BaseParser
	pubspec *PubspecParser
}

// New creates a new Dart parser
func New() *Parser {
	baseParser := common.NewBaseParser("dart-parser", "dart")

	return &Parser{
		BaseParser: baseParser,
		pubspec:    NewPubspecParser(baseParser),
	}
}

// CheckFile checks if a file should be scanned by this parser
func (p *Parser) CheckFile(filePath string) (bool, error) {
	fileName := filepath.Base(filePath)

	dartFiles := []string{
		"pubspec.yaml",
		"pubspec.lock",
	}

	return p.GetFileChecker().CheckFileByName(fileName, dartFiles), nil
}

// Scan scans the target for Dart packages
func (p *Parser) Scan(ctx context.Context, target types.ScanTarget) (*types.ScanResult, error) {
	start := time.Now()
	result := p.CreateScanResult(target)
	var allPackages []model.Package

	// Find all relevant Dart files
	files, err := p.GetFileWalker().FindFiles(target.Path, p.CheckFile)
	if err != nil {
		return result, err
	}

	for _, file := range files {
		var packages []model.Package
		var parseErr error

		fileName := filepath.Base(file)

		switch fileName {
		case "pubspec.yaml":
			packages, parseErr = p.pubspec.Parse(file)
		case "pubspec.lock":
			packages, parseErr = p.pubspec.ParseLock(file)
		}

		if parseErr != nil {
			// Log error but continue processing other files
			continue
		}

		allPackages = append(allPackages, packages...)
	}

	metadata := map[string]interface{}{
		"dart_files_processed": len(files),
	}
	p.FinalizeScanResult(result, allPackages, start, metadata)
	return result, nil
}
