package scan

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/scan/formats"
	"github.com/carbonetes/diggity/pkg/scan/parsers"
	"github.com/carbonetes/diggity/pkg/scan/registry"
	"github.com/carbonetes/diggity/pkg/scan/types"
)

// Engine coordinates the scanning process using detected scanners
type Engine struct {
	registry *registry.Registry
}

// NewEngine creates a new scanning engine with the given configuration
func NewEngine(cfg interface{}) *Engine {
	engine := &Engine{
		registry: registry.NewRegistry(),
	}

	// Auto-register default scanners
	engine.registerDefaultScanners()

	return engine
}

// ScanTarget scans the specified target and returns comprehensive results
func (e *Engine) ScanTarget(ctx context.Context, target types.ScanTarget) (*model.ComprehensiveScanResult, error) {
	start := time.Now()

	// Run all scanners on the target and collect results
	var allResults []*types.ScanResult
	var allErrors []string

	scanners := e.registry.GetAllScanners()
	for _, scanner := range scanners {
		result, err := scanner.Scan(ctx, target)
		if err != nil {
			allErrors = append(allErrors, fmt.Sprintf("%s scanner error: %v", scanner.Name(), err))
			continue
		}

		if result != nil && len(result.Packages) > 0 {
			allResults = append(allResults, result)
		}
	}

	// Convert to comprehensive result with actual package data
	return e.convertScanResultsToComprehensive(target, allResults, allErrors, time.Since(start)), nil
}

// ScanWithSpecificParsers scans using only the specified parser types
func (e *Engine) ScanWithSpecificParsers(ctx context.Context, target types.ScanTarget, parserTypes []string) (*model.ComprehensiveScanResult, error) {
	start := time.Now()

	// Run only specified scanners on the target
	var allResults []*types.ScanResult
	var allErrors []string

	scanners := e.registry.GetAllScanners()
	typeSet := make(map[string]bool)
	for _, pType := range parserTypes {
		typeSet[pType] = true
	}

	for _, scanner := range scanners {
		if !typeSet[scanner.Type()] {
			continue
		}

		result, err := scanner.Scan(ctx, target)
		if err != nil {
			allErrors = append(allErrors, fmt.Sprintf("%s scanner error: %v", scanner.Name(), err))
			continue
		}

		if result != nil && len(result.Packages) > 0 {
			allResults = append(allResults, result)
		}
	}

	// Convert to comprehensive result
	return e.convertScanResultsToComprehensive(target, allResults, allErrors, time.Since(start)), nil
}

// ScanDirectory scans a directory and returns results organized by file type
func (e *Engine) ScanDirectory(ctx context.Context, dirPath string) (map[string]*model.ComprehensiveScanResult, error) {
	files, err := e.walkDirectory(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to walk directory: %w", err)
	}

	results := make(map[string]*model.ComprehensiveScanResult)
	scanners := e.registry.GetAllScanners()

	// Group files by scanner type
	fileGroups := make(map[string][]string)
	for _, file := range files {
		for _, scanner := range scanners {
			if relevant, err := scanner.CheckFile(file); err == nil && relevant {
				fileGroups[scanner.Type()] = append(fileGroups[scanner.Type()], file)
				break // Only assign file to first matching scanner
			}
		}
	}

	// Scan each group
	for scannerType, groupFiles := range fileGroups {
		if len(groupFiles) == 0 {
			continue
		}

		// Create a target for this group
		target := types.ScanTarget{
			Type: types.TargetTypeDirectory,
			Path: dirPath,
			Metadata: map[string]string{
				"scanner_type": scannerType,
				"file_count":   fmt.Sprintf("%d", len(groupFiles)),
			},
		}

		result, err := e.ScanWithSpecificParsers(ctx, target, []string{scannerType})
		if err != nil {
			continue
		}

		results[scannerType] = result
	}

	return results, nil
}

// GetAvailableScanners returns the names of all registered scanners
func (e *Engine) GetAvailableScanners() []string {
	scanners := e.registry.GetAllScanners()
	names := make([]string, len(scanners))
	for i, scanner := range scanners {
		names[i] = scanner.Name()
	}
	return names
}

// GetScannerDetails returns detailed information about all registered scanners
func (e *Engine) GetScannerDetails() map[string]types.ScannerInfo {
	scanners := e.registry.GetAllScanners()
	details := make(map[string]types.ScannerInfo)

	for _, scanner := range scanners {
		details[scanner.Name()] = types.ScannerInfo{
			Name:        scanner.Name(),
			Type:        scanner.Type(),
			Description: fmt.Sprintf("%s package scanner", scanner.Type()),
		}
	}

	return details
}

// GetRegisteredScannerCount returns the number of registered scanners
func (e *Engine) GetRegisteredScannerCount() int {
	return len(e.registry.GetAllScanners())
}

// GenerateSBOM creates an SBOM from scan results
func (e *Engine) GenerateSBOM(target types.ScanTarget, format string) ([]byte, error) {
	// First scan to get the data
	ctx := context.Background()
	comprehensive, err := e.ScanTarget(ctx, target)
	if err != nil {
		return nil, fmt.Errorf("failed to scan target: %w", err)
	}

	// Use the formats package to generate SBOM
	formatters := map[string]types.SBOMFormat{
		"cyclonedx": formats.NewCycloneDXFormat("json"),
		"basic":     &formats.BasicFormat{},
	}

	formatter, exists := formatters[format]
	if !exists {
		return nil, fmt.Errorf("unsupported SBOM format: %s", format)
	}

	// Get actual scan results for SBOM generation
	var scanResults []types.ScanResult
	scanners := e.registry.GetAllScanners()
	for _, scanner := range scanners {
		result, err := scanner.Scan(ctx, target)
		if err != nil {
			continue
		}
		if len(result.Packages) > 0 {
			scanResults = append(scanResults, *result)
		}
	}

	// Create EngineResult for formatter
	engineResult := &types.EngineResult{
		Results:       scanResults,
		TotalPackages: len(scanResults),
		Duration:      comprehensive.Duration,
		Metadata:      make(map[string]interface{}),
	}

	// Generate SBOM
	sbomData, err := formatter.Generate(engineResult)
	if err != nil {
		return nil, fmt.Errorf("failed to generate SBOM: %w", err)
	}

	// Convert to bytes (this will depend on the format)
	switch data := sbomData.(type) {
	case []byte:
		return data, nil
	case string:
		return []byte(data), nil
	default:
		// Try JSON marshaling as fallback
		return json.Marshal(data)
	}
}

// getFilesFromTarget extracts files to scan from the target
func (e *Engine) getFilesFromTarget(target types.ScanTarget) ([]string, error) {
	switch target.Type {
	case types.TargetTypeDirectory:
		return e.walkDirectory(target.Path)
	case types.TargetTypeFile:
		return []string{target.Path}, nil
	case types.TargetTypeImage:
		// For now, treat as directory - would need container image extraction
		return e.walkDirectory(target.Path)
	case types.TargetTypeArchive:
		// For now, treat as directory - would need archive extraction
		return e.walkDirectory(target.Path)
	default:
		return []string{target.Path}, nil
	}
}

// walkDirectory recursively walks a directory and returns file paths
func (e *Engine) walkDirectory(dirPath string) ([]string, error) {
	var files []string

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})

	return files, err
}

// registerDefaultScanners registers built-in scanners
func (e *Engine) registerDefaultScanners() {
	// Create modular parsers factory
	modularParsers := parsers.NewModularParsers()

	// Register all available modular parsers
	allParsers := modularParsers.GetAllModularParsers()
	for _, parser := range allParsers {
		e.registry.Register(parser)
	}
}

// convertScanResultsToComprehensive converts scan results to comprehensive result with actual package data
func (e *Engine) convertScanResultsToComprehensive(target types.ScanTarget, results []*types.ScanResult, errors []string, duration time.Duration) *model.ComprehensiveScanResult {
	comprehensive := &model.ComprehensiveScanResult{
		Target:    target,
		Duration:  duration,
		Timestamp: time.Now(),
		Metadata:  make(map[string]interface{}),
		Errors:    errors,
	}

	// Aggregate packages from all scanners
	var allPackages []model.Package
	scannerStats := make(map[string]int)

	for _, result := range results {
		allPackages = append(allPackages, result.Packages...)
		scannerStats[result.Scanner] = len(result.Packages)

		// Merge any additional errors
		comprehensive.Errors = append(comprehensive.Errors, result.Errors...)
	}

	// Create package result with actual data
	if len(allPackages) > 0 {
		comprehensive.PackageResult = &model.ScanResult{
			ComponentsFound: len(allPackages),
			// Additional fields would be set here based on the actual scan data
		}
	}

	// Add scanner statistics to metadata
	comprehensive.Metadata["scanner_stats"] = scannerStats
	comprehensive.Metadata["total_packages"] = len(allPackages)
	comprehensive.Metadata["scanners_used"] = len(results)

	return comprehensive
}

// GetEngineStatus returns comprehensive status information about the engine
func (e *Engine) GetEngineStatus() map[string]interface{} {
	scanners := e.registry.GetAllScanners()

	// Group scanners by category
	programLanguages := []string{}
	systemPackages := []string{}
	other := []string{}

	for _, scanner := range scanners {
		switch scanner.Type() {
		case "python", "npm", "java", "go", "cargo", "php", "ruby", "dotnet", "swift", "dart", "r":
			programLanguages = append(programLanguages, scanner.Type())
		case "apt", "dpkg", "rpm", "alpine":
			systemPackages = append(systemPackages, scanner.Type())
		default:
			other = append(other, scanner.Type())
		}
	}

	return map[string]interface{}{
		"total_parsers":         len(scanners),
		"programming_languages": len(programLanguages),
		"system_packages":       len(systemPackages),
		"other_parsers":         len(other),
		"parser_categories": map[string][]string{
			"programming_languages": programLanguages,
			"system_packages":       systemPackages,
			"other":                 other,
		},
		"all_parser_types": e.GetAvailableScanners(),
	}
}
