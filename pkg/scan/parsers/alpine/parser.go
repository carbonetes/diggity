package alpine

import (
	"context"
	"path/filepath"
	"strings"
	"time"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/scan/parsers/common"
	"github.com/carbonetes/diggity/pkg/scan/types"
)

// Parser implements package parsing for Alpine Linux (APK)
type Parser struct {
	*common.BaseParser
	apkDb    *ApkDbParser
	apkIndex *ApkIndexParser
	apkbuild *ApkbuildParser
}

// New creates a new Alpine parser
func New() *Parser {
	baseParser := common.NewBaseParser("alpine-parser", "apk")

	return &Parser{
		BaseParser: baseParser,
		apkDb:      NewApkDbParser(baseParser),
		apkIndex:   NewApkIndexParser(baseParser),
		apkbuild:   NewApkbuildParser(baseParser),
	}
}

// CheckFile determines if a file is relevant for Alpine parsing
func (p *Parser) CheckFile(path string) (bool, error) {
	fileName := filepath.Base(path)

	// Check specific paths for Alpine-related files
	switch {
	case fileName == "installed" && strings.Contains(path, "/lib/apk/db/"):
		return true, nil
	case fileName == "APKINDEX.tar.gz":
		return true, nil
	case fileName == "APKBUILD":
		return true, nil
	case strings.HasPrefix(fileName, "APKINDEX") && strings.Contains(path, "/var/cache/apk/"):
		return true, nil
	}

	return false, nil
}

// Scan performs Alpine package scanning and parsing
func (p *Parser) Scan(ctx context.Context, target types.ScanTarget) (*types.ScanResult, error) {
	start := time.Now()
	p.LogStart(target)

	result := p.CreateScanResult(target)

	// Find Alpine files in the target directory
	apkFiles, err := p.GetFileWalker().FindFiles(target.Path, p.CheckFile)
	if err != nil {
		p.LogError(result, "Failed to find Alpine files", err)
		return result, nil
	}

	if len(apkFiles) == 0 {
		return result, nil
	}

	// Parse each Alpine file
	var allPackages []model.Package
	for _, file := range apkFiles {
		packages, err := p.parseApkFile(file)
		if err != nil {
			p.LogError(result, "Failed to parse "+file, err)
			continue
		}
		allPackages = append(allPackages, packages...)
	}

	// Finalize results
	metadata := map[string]interface{}{
		"apk_files_processed": len(apkFiles),
	}
	p.FinalizeScanResult(result, allPackages, start, metadata)

	return result, nil
}

// parseApkFile parses a specific Alpine file using the appropriate sub-parser
func (p *Parser) parseApkFile(filePath string) ([]model.Package, error) {
	fileName := filepath.Base(filePath)

	switch {
	case fileName == "installed" && strings.Contains(filePath, "/lib/apk/db/"):
		return p.apkDb.Parse(filePath)
	case strings.HasPrefix(fileName, "APKINDEX"):
		return p.apkIndex.Parse(filePath)
	case fileName == "APKBUILD":
		return p.apkbuild.Parse(filePath)
	}

	return []model.Package{}, nil
}
