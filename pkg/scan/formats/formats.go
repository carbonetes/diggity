package formats

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"time"

	"github.com/CycloneDX/cyclonedx-go"
	"github.com/carbonetes/diggity/pkg/scan/types"
)

// CycloneDXFormat implements SBOMFormat for CycloneDX
type CycloneDXFormat struct {
	format string // json or xml
}

// NewCycloneDXFormat creates a new CycloneDX format handler
func NewCycloneDXFormat(format string) *CycloneDXFormat {
	return &CycloneDXFormat{format: format}
}

// Name returns the format name
func (f *CycloneDXFormat) Name() string {
	return fmt.Sprintf("cyclonedx-%s", f.format)
}

// Generate generates CycloneDX SBOM from scan results
func (f *CycloneDXFormat) Generate(result *types.EngineResult) (interface{}, error) {
	bom := f.createBOM(result)

	switch f.format {
	case "json":
		return f.generateJSON(bom)
	case "xml":
		return f.generateXML(bom)
	default:
		return nil, fmt.Errorf("unsupported CycloneDX format: %s", f.format)
	}
}

func (f *CycloneDXFormat) createBOM(result *types.EngineResult) *cyclonedx.BOM {
	bom := &cyclonedx.BOM{
		BOMFormat:   "CycloneDX",
		SpecVersion: cyclonedx.SpecVersion1_5,
		Version:     1,
		Metadata: &cyclonedx.Metadata{
			Timestamp: time.Now().Format(time.RFC3339),
		},
		Components: &[]cyclonedx.Component{},
	}

	// Add components from scan results
	components := make([]cyclonedx.Component, 0)

	for _, scanResult := range result.Results {
		for _, pkg := range scanResult.Packages {
			comp := cyclonedx.Component{
				Type:       cyclonedx.ComponentTypeLibrary,
				BOMRef:     fmt.Sprintf("%s@%s", pkg.Name, pkg.Version),
				Name:       pkg.Name,
				Version:    pkg.Version,
				PackageURL: pkg.PackageURL,
			}
			components = append(components, comp)
		}
	}

	bom.Components = &components
	return bom
}

func (f *CycloneDXFormat) generateJSON(bom *cyclonedx.BOM) (interface{}, error) {
	data, err := json.MarshalIndent(bom, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return string(data), nil
}

func (f *CycloneDXFormat) generateXML(bom *cyclonedx.BOM) (interface{}, error) {
	data, err := xml.MarshalIndent(bom, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal XML: %w", err)
	}
	return string(data), nil
}

// BasicFormat implements a simple JSON format for scan results
type BasicFormat struct{}

// NewBasicFormat creates a new basic format handler
func NewBasicFormat() *BasicFormat {
	return &BasicFormat{}
}

// Name returns the format name
func (f *BasicFormat) Name() string {
	return "json"
}

// Generate generates a basic JSON representation of scan results
func (f *BasicFormat) Generate(result *types.EngineResult) (interface{}, error) {
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return string(data), nil
}

// GetFormatByName returns a format instance by name
func GetFormatByName(name string) (types.SBOMFormat, error) {
	switch name {
	case "cyclonedx-json":
		return NewCycloneDXFormat("json"), nil
	case "cyclonedx-xml":
		return NewCycloneDXFormat("xml"), nil
	case "json":
		return NewBasicFormat(), nil
	default:
		return nil, fmt.Errorf("unsupported format: %s", name)
	}
}

// GetAvailableFormats returns all available format names
func GetAvailableFormats() []string {
	return []string{
		"cyclonedx-json",
		"cyclonedx-xml",
		"json",
	}
}
