package scan

import (
	"github.com/carbonetes/diggity/pkg/scan/types"
)

// Re-export essential types for external use
type (
	Scanner    = types.Scanner
	ScanTarget = types.ScanTarget
	ScanResult = types.ScanResult
	TargetType = types.TargetType
)

// Re-export target type constants
const (
	TargetTypeFile      = types.TargetTypeFile
	TargetTypeDirectory = types.TargetTypeDirectory
	TargetTypeImage     = types.TargetTypeImage
	TargetTypeArchive   = types.TargetTypeArchive
)

// NewScanTarget creates a new scan target for the given path
func NewScanTarget(targetType TargetType, path string) ScanTarget {
	return ScanTarget{
		Type:     targetType,
		Path:     path,
		Metadata: make(map[string]string),
	}
}
