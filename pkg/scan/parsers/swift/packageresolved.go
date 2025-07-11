package swift

import (
	"encoding/json"
	"io/ioutil"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/scan/parsers/common"
)

// PackageResolvedParser handles parsing of Package.resolved files
type PackageResolvedParser struct {
	base *common.BaseParser
}

// NewPackageResolvedParser creates a new Package.resolved parser
func NewPackageResolvedParser(base *common.BaseParser) *PackageResolvedParser {
	return &PackageResolvedParser{base: base}
}

// PackageResolvedV1 represents the v1 format of Package.resolved
type PackageResolvedV1 struct {
	Object struct {
		Pins []PinV1 `json:"pins"`
	} `json:"object"`
	Version int `json:"version"`
}

// PinV1 represents a dependency pin in v1 format
type PinV1 struct {
	Package       string  `json:"package"`
	RepositoryURL string  `json:"repositoryURL"`
	State         StateV1 `json:"state"`
}

// StateV1 represents the state of a dependency in v1 format
type StateV1 struct {
	Branch   string `json:"branch,omitempty"`
	Revision string `json:"revision"`
	Version  string `json:"version,omitempty"`
}

// PackageResolvedV2 represents the v2 format of Package.resolved
type PackageResolvedV2 struct {
	Pins    []PinV2 `json:"pins"`
	Version int     `json:"version"`
}

// PinV2 represents a dependency pin in v2 format
type PinV2 struct {
	Identity string  `json:"identity"`
	Kind     string  `json:"kind"`
	Location string  `json:"location"`
	State    StateV2 `json:"state"`
}

// StateV2 represents the state of a dependency in v2 format
type StateV2 struct {
	Branch   string `json:"branch,omitempty"`
	Revision string `json:"revision"`
	Version  string `json:"version,omitempty"`
}

// Parse parses Package.resolved files
func (p *PackageResolvedParser) Parse(filePath string) ([]model.Package, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var packages []model.Package

	// Try parsing as v2 format first
	var v2Resolved PackageResolvedV2
	if err := json.Unmarshal(content, &v2Resolved); err == nil && v2Resolved.Version == 2 {
		for _, pin := range v2Resolved.Pins {
			metadata := map[string]interface{}{
				"manager":  "swift-package-manager",
				"location": pin.Location,
				"kind":     pin.Kind,
			}

			version := pin.State.Version
			if version == "" && pin.State.Branch != "" {
				version = pin.State.Branch
			}

			pkg := p.base.CreatePackage(pin.Identity, version, filePath, true, metadata)
			packages = append(packages, pkg)
		}
		return packages, nil
	}

	// Try parsing as v1 format
	var v1Resolved PackageResolvedV1
	if err := json.Unmarshal(content, &v1Resolved); err == nil && v1Resolved.Version == 1 {
		for _, pin := range v1Resolved.Object.Pins {
			metadata := map[string]interface{}{
				"manager":       "swift-package-manager",
				"repositoryURL": pin.RepositoryURL,
			}

			version := pin.State.Version
			if version == "" && pin.State.Branch != "" {
				version = pin.State.Branch
			}

			pkg := p.base.CreatePackage(pin.Package, version, filePath, true, metadata)
			packages = append(packages, pkg)
		}
		return packages, nil
	}

	return packages, nil
}
