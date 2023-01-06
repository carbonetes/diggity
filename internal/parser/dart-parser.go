package parser

import (
	"errors"

	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/internal/file"
	"github.com/carbonetes/diggity/internal/model"

	"os"

	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

const (
	pubspecYaml = "pubspec.yaml"
	pub         = "pub"
	dart        = "dart"
)

// DartMetadata metadata
type DartMetadata map[string]interface{}

// FindDartPackagesFromContent - find dart packages from content
func FindDartPackagesFromContent() {
	if parserEnabled(dart) {
		for _, content := range file.Contents {
			if filepath.Base(content.Path) == pubspecYaml {
				if err := parseDartPackages(content); err != nil {
					err = errors.New("dart-parser: " + err.Error())
					Errors = append(Errors, &err)
				}
			}
		}
	}
	defer WG.Done()
}

// Parse dart package metadata
func parseDartPackages(location *model.Location) error {
	var licenses []string = make([]string, 0)
	byteValue, err := os.ReadFile(location.Path)
	if err != nil {
		return err
	}
	metadata := make(DartMetadata)

	if err := yaml.Unmarshal([]byte(byteValue), &metadata); err != nil {
		return err
	}

	_package := new(model.Package)
	_package.ID = uuid.NewString()
	_package.Name = metadata["name"].(string)
	_package.Type = pub
	_package.Path = metadata["name"].(string)

	//check if version exist, if not set default of 0.0.0
	if val, ok := metadata["version"].(string); ok {
		_package.Version = val
	} else {
		_package.Version = "0.0.0"
	}

	_package.Locations = append(_package.Locations, model.Location{
		Path:      TrimUntilLayer(*location),
		LayerHash: location.LayerHash,
	})

	if val, ok := metadata["description"].(string); ok {
		_package.Description = val
	}

	if val, ok := metadata["license"].(string); ok {
		licenses = append(licenses, val)
	} else {
		licenses = append(licenses, "BSD 3-Clause")
	}
	_package.Licenses = licenses

	//parse CPE
	if val, ok := metadata["author"].(string); ok {
		cpe.NewCPE23(_package, strings.TrimSpace(val), _package.Name, _package.Version)
	} else if val, ok := metadata["authors"].(string); ok {
		cpe.NewCPE23(_package, strings.TrimSpace(val), _package.Name, _package.Version)
	} else {
		cpe.NewCPE23(_package, _package.Name, _package.Name, _package.Version)
	}

	parseDartPURL(_package)
	_package.Metadata = metadata
	Packages = append(Packages, _package)
	return nil
}

// Parse PURL
func parseDartPURL(_package *model.Package) {
	_package.PURL = model.PURL(scheme + ":" + "dart" + "/" + _package.Name + "@" + _package.Version)
}
