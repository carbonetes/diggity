package types

import (
	"context"
	"time"

	"github.com/carbonetes/diggity/pkg/artifact"
	"github.com/carbonetes/diggity/pkg/model"
)

// TargetType represents the type of scan target
type TargetType string

const (
	TargetTypeFile      TargetType = "file"
	TargetTypeDirectory TargetType = "directory"
	TargetTypeImage     TargetType = "image"
	TargetTypeArchive   TargetType = "archive"
)

// ScanTarget represents what should be scanned
type ScanTarget struct {
	Type     TargetType        `json:"type"`
	Path     string            `json:"path"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

// Scanner defines the interface for all scanners
type Scanner interface {
	// Name returns the scanner name
	Name() string

	// Type returns the scanner type (e.g., "npm", "cargo", "dpkg")
	Type() string

	// CheckFile determines if a file is relevant for this scanner
	CheckFile(path string) (bool, error)

	// Scan performs the actual scanning
	Scan(ctx context.Context, target ScanTarget) (*ScanResult, error)
}

// ScanResult contains the results of a scan operation
type ScanResult struct {
	Scanner         string                 `json:"scanner"`
	Target          ScanTarget             `json:"target"`
	Packages        []model.Package        `json:"packages"`
	Secrets         []model.Secret         `json:"secrets,omitempty"`
	Licenses        []model.License        `json:"licenses,omitempty"`
	ComponentsFound int                    `json:"components_found"`
	Duration        time.Duration          `json:"duration"`
	Metadata        map[string]interface{} `json:"metadata"`
	Errors          []string               `json:"errors,omitempty"`
}

// ScannerInfo provides detailed information about a scanner
type ScannerInfo struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

// ScanConfig contains configuration for scanning operations
type ScanConfig struct {
	IncludeSecrets  bool          `json:"include_secrets"`
	IncludeLicenses bool          `json:"include_licenses"`
	ExcludePaths    []string      `json:"exclude_paths,omitempty"`
	MaxDepth        int           `json:"max_depth"`
	Timeout         time.Duration `json:"timeout"`
}

// DefaultScanConfig returns a default scan configuration
func DefaultScanConfig() *ScanConfig {
	return &ScanConfig{
		IncludeSecrets:  false,
		IncludeLicenses: true,
		ExcludePaths:    []string{".git", "node_modules", ".vscode"},
		MaxDepth:        -1, // unlimited
		Timeout:         5 * time.Minute,
	}
}

// EngineResult represents the final result of a scan engine operation
type EngineResult struct {
	Artifact      *artifact.Artifact     `json:"artifact"`
	Results       []ScanResult           `json:"results"`
	TotalPackages int                    `json:"total_packages"`
	TotalSecrets  int                    `json:"total_secrets"`
	TotalLicenses int                    `json:"total_licenses"`
	Duration      time.Duration          `json:"duration"`
	SBOMs         map[string]interface{} `json:"sboms,omitempty"`
	Metadata      map[string]interface{} `json:"metadata"`
	Errors        []string               `json:"errors,omitempty"`
}

// SBOMFormat defines the interface for SBOM format generators
type SBOMFormat interface {
	Name() string
	Generate(result *EngineResult) (interface{}, error)
}
