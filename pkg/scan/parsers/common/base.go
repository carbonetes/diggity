package common

import (
	"time"

	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/scan/types"
)

// BaseParser provides common functionality for all parsers
type BaseParser struct {
	name        string
	parserType  string
	fileWalker  *FileWalker
	fileCleaner *VersionCleaner
	fileChecker *FileChecker
}

// NewBaseParser creates a new BaseParser instance
func NewBaseParser(name, parserType string) *BaseParser {
	return &BaseParser{
		name:        name,
		parserType:  parserType,
		fileWalker:  NewFileWalker(),
		fileCleaner: NewVersionCleaner(),
		fileChecker: NewFileChecker(),
	}
}

// Name returns the parser name
func (bp *BaseParser) Name() string {
	return bp.name
}

// Type returns the parser type
func (bp *BaseParser) Type() string {
	return bp.parserType
}

// GetFileWalker returns the file walker instance
func (bp *BaseParser) GetFileWalker() *FileWalker {
	return bp.fileWalker
}

// GetVersionCleaner returns the version cleaner instance
func (bp *BaseParser) GetVersionCleaner() *VersionCleaner {
	return bp.fileCleaner
}

// GetFileChecker returns the file checker instance
func (bp *BaseParser) GetFileChecker() *FileChecker {
	return bp.fileChecker
}

// CreateScanResult creates a new scan result with common fields populated
func (bp *BaseParser) CreateScanResult(target types.ScanTarget) *types.ScanResult {
	return &types.ScanResult{
		Scanner:         bp.Name(),
		Target:          target,
		Packages:        []model.Package{},
		ComponentsFound: 0,
		Duration:        0,
		Metadata:        make(map[string]interface{}),
		Errors:          []string{},
	}
}

// FinalizeScanResult finalizes the scan result with computed values
func (bp *BaseParser) FinalizeScanResult(result *types.ScanResult, packages []model.Package, start time.Time, additionalMetadata map[string]interface{}) {
	result.Packages = packages
	result.ComponentsFound = len(packages)
	result.Duration = time.Since(start)

	// Add additional metadata
	for key, value := range additionalMetadata {
		result.Metadata[key] = value
	}

	log.Debugf("%s scan completed: found %d packages", bp.name, len(packages))
}

// LogStart logs the start of parsing for a target
func (bp *BaseParser) LogStart(target types.ScanTarget) {
	log.Debugf("%s parser processing: %s", bp.name, target.Path)
}

// LogError logs an error and adds it to the result
func (bp *BaseParser) LogError(result *types.ScanResult, message string, err error) {
	errorMsg := message
	if err != nil {
		errorMsg = message + ": " + err.Error()
	}
	result.Errors = append(result.Errors, errorMsg)
	log.Debugf("Error in %s parser: %s", bp.name, errorMsg)
}

// CreatePackage creates a standard package model
func (bp *BaseParser) CreatePackage(name, version, sourceFile string, isDev bool, metadata map[string]interface{}) model.Package {
	if metadata == nil {
		metadata = make(map[string]interface{})
	}

	metadata["source_file"] = sourceFile
	metadata["is_dev"] = isDev

	return model.Package{
		Name:     name,
		Version:  version,
		Type:     bp.parserType,
		Language: bp.parserType,
		Metadata: metadata,
	}
}
