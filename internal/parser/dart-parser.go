package parser

import (
	"errors"

	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/internal/file"
	"github.com/carbonetes/diggity/internal/model"
	"github.com/carbonetes/diggity/internal/model/metadata"

	"os"

	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

const (
	pubspecYaml = "pubspec.yaml"
	pubspecLock = "pubspec.lock"
	pub         = "pub"
	dart        = "dart"
)

// DartMetadata metadata
type DartMetadata map[string]interface{}

var dartlockFileMetadata metadata.PubspecLockPackage

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
			if filepath.Base(content.Path) == pubspecLock {
				if err := parseDartPackagesLock(content); err != nil {
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

// Parse dart packages metadata - lock file
func parseDartPackagesLock(location *model.Location) error {
	byteValue, err := os.ReadFile(location.Path)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal([]byte(byteValue), &dartlockFileMetadata); err != nil {
		return err
	}

	for _, cPackage := range dartlockFileMetadata.Packages {
		_package := new(model.Package)
		_package.ID = uuid.NewString()
		_package.Name = cPackage.Description.Name
		_package.Version = cPackage.Version
		_package.Type = pub
		_package.Path = cPackage.Description.Name
		_package.Locations = append(_package.Locations, model.Location{
			Path:      TrimUntilLayer(*location),
			LayerHash: location.LayerHash,
		})
		cpe.NewCPE23(_package, _package.Name, _package.Name, _package.Version)
		parseDartPURL(_package)
		_package.Metadata = cPackage
		if _package.Name != "" {
			Packages = append(Packages, _package)
		}
	}
	return nil
}

// Parse PURL
func parseDartPURL(_package *model.Package) {
	_package.PURL = model.PURL(scheme + ":" + "dart" + "/" + _package.Name + "@" + _package.Version)
}
