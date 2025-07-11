package apt

import (
	"context"
	"path/filepath"
	"strings"
	"time"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/scan/parsers/common"
	"github.com/carbonetes/diggity/pkg/scan/types"
)

// Parser implements package parsing for APT package manager
type Parser struct {
	*common.BaseParser
	lists   *ListsParser
	cache   *CacheParser
	sources *SourcesParser
}

// New creates a new APT parser
func New() *Parser {
	baseParser := common.NewBaseParser("apt-parser", "apt")

	return &Parser{
		BaseParser: baseParser,
		lists:      NewListsParser(baseParser),
		cache:      NewCacheParser(baseParser),
		sources:    NewSourcesParser(baseParser),
	}
}

// CheckFile determines if a file is relevant for APT parsing
func (p *Parser) CheckFile(path string) (bool, error) {
	fileName := filepath.Base(path)

	// Check specific paths for APT-specific files (not DPKG)
	switch {
	case fileName == "Packages" && strings.Contains(path, "/var/lib/apt/lists/"):
		return true, nil
	case fileName == "sources.list" && strings.Contains(path, "/etc/apt/"):
		return true, nil
	case strings.HasSuffix(fileName, ".list") && strings.Contains(path, "/etc/apt/sources.list.d/"):
		return true, nil
	case strings.Contains(path, "/var/cache/apt/"):
		return true, nil
	}

	return false, nil
}

// Scan performs APT package scanning and parsing
func (p *Parser) Scan(ctx context.Context, target types.ScanTarget) (*types.ScanResult, error) {
	start := time.Now()
	p.LogStart(target)

	result := p.CreateScanResult(target)

	// Find APT files in the target directory
	aptFiles, err := p.GetFileWalker().FindFiles(target.Path, p.CheckFile)
	if err != nil {
		p.LogError(result, "Failed to find APT files", err)
		return result, nil
	}

	if len(aptFiles) == 0 {
		return result, nil
	}

	// Parse each APT file
	var allPackages []model.Package
	for _, file := range aptFiles {
		packages, err := p.parseAptFile(file)
		if err != nil {
			p.LogError(result, "Failed to parse "+file, err)
			continue
		}
		allPackages = append(allPackages, packages...)
	}

	// Finalize results
	metadata := map[string]interface{}{
		"apt_files_processed": len(aptFiles),
	}
	p.FinalizeScanResult(result, allPackages, start, metadata)

	return result, nil
}

// parseAptFile parses a specific APT file using the appropriate sub-parser
func (p *Parser) parseAptFile(filePath string) ([]model.Package, error) {
	fileName := filepath.Base(filePath)

	switch {
	case fileName == "Packages" && strings.Contains(filePath, "/var/lib/apt/lists/"):
		return p.lists.Parse(filePath)
	case fileName == "sources.list" && strings.Contains(filePath, "/etc/apt/"):
		return p.sources.Parse(filePath)
	case strings.HasSuffix(fileName, ".list") && strings.Contains(filePath, "/etc/apt/sources.list.d/"):
		return p.sources.Parse(filePath)
	case strings.Contains(filePath, "/var/cache/apt/"):
		return p.cache.Parse(filePath)
	}

	return []model.Package{}, nil
}
