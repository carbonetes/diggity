package pacman

import (
	"context"
	"path/filepath"
	"strings"
	"time"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/scan/parsers/common"
	"github.com/carbonetes/diggity/pkg/scan/types"
)

// Parser implements package parsing for Pacman-based systems (Arch Linux)
type Parser struct {
	*common.BaseParser
	pacmanDb   *PacmanDbParser
	pkgbuild   *PkgbuildParser
	mirrorlist *MirrorlistParser
}

// New creates a new Pacman parser
func New() *Parser {
	baseParser := common.NewBaseParser("pacman-parser", "pacman")

	return &Parser{
		BaseParser: baseParser,
		pacmanDb:   NewPacmanDbParser(baseParser),
		pkgbuild:   NewPkgbuildParser(baseParser),
		mirrorlist: NewMirrorlistParser(baseParser),
	}
}

// CheckFile determines if a file is relevant for Pacman parsing
func (p *Parser) CheckFile(path string) (bool, error) {
	fileName := filepath.Base(path)

	// Check specific paths for Pacman-related files
	switch {
	case fileName == "desc" && strings.Contains(path, "/var/lib/pacman/local/"):
		return true, nil
	case fileName == "PKGBUILD":
		return true, nil
	case fileName == "mirrorlist" && strings.Contains(path, "/etc/pacman.d/"):
		return true, nil
	case fileName == "pacman.conf" && strings.Contains(path, "/etc/"):
		return true, nil
	}

	return false, nil
}

// Scan performs Pacman package scanning and parsing
func (p *Parser) Scan(ctx context.Context, target types.ScanTarget) (*types.ScanResult, error) {
	start := time.Now()
	p.LogStart(target)

	result := p.CreateScanResult(target)

	// Find Pacman files in the target directory
	pacmanFiles, err := p.GetFileWalker().FindFiles(target.Path, p.CheckFile)
	if err != nil {
		p.LogError(result, "Failed to find Pacman files", err)
		return result, nil
	}

	if len(pacmanFiles) == 0 {
		return result, nil
	}

	// Parse each Pacman file
	var allPackages []model.Package
	for _, file := range pacmanFiles {
		packages, err := p.parsePacmanFile(file)
		if err != nil {
			p.LogError(result, "Failed to parse "+file, err)
			continue
		}
		allPackages = append(allPackages, packages...)
	}

	// Finalize results
	metadata := map[string]interface{}{
		"pacman_files_processed": len(pacmanFiles),
	}
	p.FinalizeScanResult(result, allPackages, start, metadata)

	return result, nil
}

// parsePacmanFile parses a specific Pacman file using the appropriate sub-parser
func (p *Parser) parsePacmanFile(filePath string) ([]model.Package, error) {
	fileName := filepath.Base(filePath)

	switch {
	case fileName == "desc" && strings.Contains(filePath, "/var/lib/pacman/local/"):
		return p.pacmanDb.Parse(filePath)
	case fileName == "PKGBUILD":
		return p.pkgbuild.Parse(filePath)
	case fileName == "mirrorlist" && strings.Contains(filePath, "/etc/pacman.d/"):
		return p.mirrorlist.Parse(filePath)
	}

	return []model.Package{}, nil
}
