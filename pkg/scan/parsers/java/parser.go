package java

import (
	"context"
	"time"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/scan/parsers/common"
	"github.com/carbonetes/diggity/pkg/scan/types"
)

// Parser implements package parsing for Java projects
type Parser struct {
	*common.BaseParser
	maven  *MavenParser
	gradle *GradleParser
}

// New creates a new Java parser
func New() *Parser {
	baseParser := common.NewBaseParser("java-parser", "java")

	return &Parser{
		BaseParser: baseParser,
		maven:      NewMavenParser(baseParser),
		gradle:     NewGradleParser(baseParser),
	}
}

// CheckFile determines if a file is relevant for java parsing
func (p *Parser) CheckFile(path string) (bool, error) {
	javaFiles := []string{
		"pom.xml",
		"build.gradle",
		"build.gradle.kts",
	}

	return p.GetFileChecker().CheckFileByName(path, javaFiles), nil
}

// Scan performs java package scanning and parsing
func (p *Parser) Scan(ctx context.Context, target types.ScanTarget) (*types.ScanResult, error) {
	start := time.Now()
	p.LogStart(target)

	result := p.CreateScanResult(target)

	// Find java files in the target directory
	javaFiles, err := p.GetFileWalker().FindFiles(target.Path, p.CheckFile)
	if err != nil {
		p.LogError(result, "Failed to find java files", err)
		return result, nil
	}

	if len(javaFiles) == 0 {
		return result, nil
	}

	// Parse each java file
	var allPackages []model.Package
	for _, file := range javaFiles {
		packages, err := p.parseJavaFile(file)
		if err != nil {
			p.LogError(result, "Failed to parse "+file, err)
			continue
		}
		allPackages = append(allPackages, packages...)
	}

	// Finalize results
	metadata := map[string]interface{}{
		"java_files_processed": len(javaFiles),
	}
	p.FinalizeScanResult(result, allPackages, start, metadata)

	return result, nil
}

// parseJavaFile parses a specific java file using the appropriate sub-parser
func (p *Parser) parseJavaFile(filePath string) ([]model.Package, error) {
	checker := p.GetFileChecker()

	if checker.CheckFileByName(filePath, []string{"pom.xml"}) {
		return p.maven.Parse(filePath)
	}

	if checker.CheckFileByName(filePath, []string{"build.gradle", "build.gradle.kts"}) {
		return p.gradle.Parse(filePath)
	}

	return []model.Package{}, nil
}
