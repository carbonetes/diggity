package rpm

import (
	"context"
	"path/filepath"
	"strings"
	"time"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/scan/parsers/common"
	"github.com/carbonetes/diggity/pkg/scan/types"
)

// Parser implements package parsing for RPM-based systems (Red Hat/CentOS/Fedora)
type Parser struct {
	*common.BaseParser
	rpmDb    *RpmDbParser
	yumRepos *YumReposParser
	dnfRepos *DnfReposParser
	rpmSpec  *RpmSpecParser
}

// New creates a new RPM parser
func New() *Parser {
	baseParser := common.NewBaseParser("rpm-parser", "rpm")

	return &Parser{
		BaseParser: baseParser,
		rpmDb:      NewRpmDbParser(baseParser),
		yumRepos:   NewYumReposParser(baseParser),
		dnfRepos:   NewDnfReposParser(baseParser),
		rpmSpec:    NewRpmSpecParser(baseParser),
	}
}

// CheckFile determines if a file is relevant for RPM parsing
func (p *Parser) CheckFile(path string) (bool, error) {
	fileName := filepath.Base(path)

	// Check specific paths for RPM-related files
	switch {
	case fileName == "Packages" && strings.Contains(path, "/var/lib/rpm/"):
		return true, nil
	case strings.HasSuffix(fileName, ".repo") && strings.Contains(path, "/etc/yum.repos.d/"):
		return true, nil
	case strings.HasSuffix(fileName, ".repo") && strings.Contains(path, "/etc/dnf/repos.d/"):
		return true, nil
	case strings.HasSuffix(fileName, ".spec"):
		return true, nil
	case fileName == "repodata" && strings.Contains(path, "/repodata/"):
		return true, nil
	}

	return false, nil
}

// Scan performs RPM package scanning and parsing
func (p *Parser) Scan(ctx context.Context, target types.ScanTarget) (*types.ScanResult, error) {
	start := time.Now()
	p.LogStart(target)

	result := p.CreateScanResult(target)

	// Find RPM files in the target directory
	rpmFiles, err := p.GetFileWalker().FindFiles(target.Path, p.CheckFile)
	if err != nil {
		p.LogError(result, "Failed to find RPM files", err)
		return result, nil
	}

	if len(rpmFiles) == 0 {
		return result, nil
	}

	// Parse each RPM file
	var allPackages []model.Package
	for _, file := range rpmFiles {
		packages, err := p.parseRpmFile(file)
		if err != nil {
			p.LogError(result, "Failed to parse "+file, err)
			continue
		}
		allPackages = append(allPackages, packages...)
	}

	// Finalize results
	metadata := map[string]interface{}{
		"rpm_files_processed": len(rpmFiles),
	}
	p.FinalizeScanResult(result, allPackages, start, metadata)

	return result, nil
}

// parseRpmFile parses a specific RPM file using the appropriate sub-parser
func (p *Parser) parseRpmFile(filePath string) ([]model.Package, error) {
	fileName := filepath.Base(filePath)

	switch {
	case fileName == "Packages" && strings.Contains(filePath, "/var/lib/rpm/"):
		return p.rpmDb.Parse(filePath)
	case strings.HasSuffix(fileName, ".repo") && strings.Contains(filePath, "/etc/yum.repos.d/"):
		return p.yumRepos.Parse(filePath)
	case strings.HasSuffix(fileName, ".repo") && strings.Contains(filePath, "/etc/dnf/repos.d/"):
		return p.dnfRepos.Parse(filePath)
	case strings.HasSuffix(fileName, ".spec"):
		return p.rpmSpec.Parse(filePath)
	}

	return []model.Package{}, nil
}
