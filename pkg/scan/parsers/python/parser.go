package python

import (
	"context"
	"time"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/scan/parsers/common"
	"github.com/carbonetes/diggity/pkg/scan/types"
)

// Parser implements package parsing for Python projects
type Parser struct {
	*common.BaseParser
	requirements *RequirementsParser
	pipfile      *PipfileParser
	setupPy      *SetupPyParser
	pyproject    *PyprojectParser
}

// New creates a new Python parser
func New() *Parser {
	baseParser := common.NewBaseParser("python-parser", "python")

	return &Parser{
		BaseParser:   baseParser,
		requirements: NewRequirementsParser(baseParser),
		pipfile:      NewPipfileParser(baseParser),
		setupPy:      NewSetupPyParser(baseParser),
		pyproject:    NewPyprojectParser(baseParser),
	}
}

// CheckFile determines if a file is relevant for python parsing
func (p *Parser) CheckFile(path string) (bool, error) {
	pythonFiles := []string{
		"requirements.txt",
		"requirements-dev.txt",
		"Pipfile",
		"Pipfile.lock",
		"pyproject.toml",
		"setup.py",
		"poetry.lock",
	}

	return p.GetFileChecker().CheckFileByName(path, pythonFiles), nil
}

// Scan performs python package scanning and parsing
func (p *Parser) Scan(ctx context.Context, target types.ScanTarget) (*types.ScanResult, error) {
	start := time.Now()
	p.LogStart(target)

	result := p.CreateScanResult(target)

	// Find python files in the target directory
	pythonFiles, err := p.GetFileWalker().FindFiles(target.Path, p.CheckFile)
	if err != nil {
		p.LogError(result, "Failed to find python files", err)
		return result, nil
	}

	if len(pythonFiles) == 0 {
		return result, nil
	}

	// Parse each python file
	var allPackages []model.Package
	for _, file := range pythonFiles {
		packages, err := p.parsePythonFile(file)
		if err != nil {
			p.LogError(result, "Failed to parse "+file, err)
			continue
		}
		allPackages = append(allPackages, packages...)
	}

	// Finalize results
	metadata := map[string]interface{}{
		"python_files_processed": len(pythonFiles),
	}
	p.FinalizeScanResult(result, allPackages, start, metadata)

	return result, nil
}

// parsePythonFile parses a specific python file using the appropriate sub-parser
func (p *Parser) parsePythonFile(filePath string) ([]model.Package, error) {
	checker := p.GetFileChecker()

	if checker.CheckFileByName(filePath, []string{"requirements.txt", "requirements-dev.txt"}) {
		return p.requirements.Parse(filePath)
	}

	if checker.CheckFileByName(filePath, []string{"Pipfile"}) {
		return p.pipfile.Parse(filePath)
	}

	if checker.CheckFileByName(filePath, []string{"setup.py"}) {
		return p.setupPy.Parse(filePath)
	}

	if checker.CheckFileByName(filePath, []string{"pyproject.toml"}) {
		return p.pyproject.Parse(filePath)
	}

	return []model.Package{}, nil
}
