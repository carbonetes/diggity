package dpkg

import (
	"context"
	"path/filepath"
	"strings"
	"time"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/scan/parsers/common"
	"github.com/carbonetes/diggity/pkg/scan/types"
)

// Parser implements package parsing for DPKG (Debian Package Manager)
type Parser struct {
	*common.BaseParser
	status  *StatusParser
	info    *InfoParser
	control *ControlParser
}

// New creates a new DPKG parser
func New() *Parser {
	baseParser := common.NewBaseParser("dpkg-parser", "dpkg")

	return &Parser{
		BaseParser: baseParser,
		status:     NewStatusParser(baseParser),
		info:       NewInfoParser(baseParser),
		control:    NewControlParser(baseParser),
	}
}

// CheckFile determines if a file is relevant for DPKG parsing
func (p *Parser) CheckFile(path string) (bool, error) {
	fileName := filepath.Base(path)

	// Check specific paths for DPKG-related files
	switch {
	case fileName == "status" && strings.Contains(path, "/var/lib/dpkg/"):
		return true, nil
	case fileName == "control" && strings.Contains(path, "DEBIAN/"):
		return true, nil
	case strings.Contains(path, "/var/lib/dpkg/info/") && (strings.HasSuffix(fileName, ".list") || strings.HasSuffix(fileName, ".md5sums") || strings.HasSuffix(fileName, ".conffiles")):
		return true, nil
	}

	return false, nil
}

// Scan performs DPKG package scanning and parsing
func (p *Parser) Scan(ctx context.Context, target types.ScanTarget) (*types.ScanResult, error) {
	start := time.Now()
	p.LogStart(target)

	result := p.CreateScanResult(target)

	// Find DPKG files in the target directory
	dpkgFiles, err := p.GetFileWalker().FindFiles(target.Path, p.CheckFile)
	if err != nil {
		p.LogError(result, "Failed to find DPKG files", err)
		return result, nil
	}

	if len(dpkgFiles) == 0 {
		return result, nil
	}

	// Parse each DPKG file
	var allPackages []model.Package
	for _, file := range dpkgFiles {
		packages, err := p.parseDpkgFile(file)
		if err != nil {
			p.LogError(result, "Failed to parse "+file, err)
			continue
		}
		allPackages = append(allPackages, packages...)
	}

	// Finalize results
	metadata := map[string]interface{}{
		"dpkg_files_processed": len(dpkgFiles),
	}
	p.FinalizeScanResult(result, allPackages, start, metadata)

	return result, nil
}

// parseDpkgFile parses a specific DPKG file using the appropriate sub-parser
func (p *Parser) parseDpkgFile(filePath string) ([]model.Package, error) {
	fileName := filepath.Base(filePath)

	switch {
	case fileName == "status" && strings.Contains(filePath, "/var/lib/dpkg/"):
		return p.status.Parse(filePath)
	case fileName == "control" && strings.Contains(filePath, "DEBIAN/"):
		return p.control.Parse(filePath)
	case strings.Contains(filePath, "/var/lib/dpkg/info/"):
		return p.info.Parse(filePath)
	}

	return []model.Package{}, nil
}
