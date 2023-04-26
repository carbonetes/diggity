package dart

import (
	"errors"

	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/model/metadata"
	"github.com/carbonetes/diggity/pkg/parser/bom"
	"github.com/carbonetes/diggity/pkg/parser/util"

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

// Metadata metadata
type Metadata map[string]interface{}

var dartlockFileMetadata metadata.PubspecLockPackage

// FindDartPackagesFromContent - find dart packages from content
func FindDartPackagesFromContent(req *bom.ParserRequirements) {
	if util.ParserEnabled(dart, req.Arguments.EnabledParsers) {
		for _, content := range *req.Contents {
			if filepath.Base(content.Path) == pubspecYaml {
				if err := parseDartPackages(&content, req.Result.Packages); err != nil {
					err = errors.New("dart-parser: " + err.Error())
					*req.Errors = append(*req.Errors, err)
				}
			}
			if filepath.Base(content.Path) == pubspecLock {
				if err := parseDartPackagesLock(&content, req.Result.Packages); err != nil {
					err = errors.New("dart-parser: " + err.Error())
					*req.Errors = append(*req.Errors, err)
				}
			}
		}
	}
	defer req.WG.Done()
}

// Parse dart package metadata
func parseDartPackages(location *model.Location, pkgs *[]model.Package) error {
	var licenses []string = make([]string, 0)
	byteValue, err := os.ReadFile(location.Path)
	if err != nil {
		return err
	}
	metadata := make(Metadata)

	if err := yaml.Unmarshal([]byte(byteValue), &metadata); err != nil {
		return err
	}

	pkg := new(model.Package)
	pkg.ID = uuid.NewString()
	pkg.Name = metadata["name"].(string)
	pkg.Type = pub
	pkg.Path = metadata["name"].(string)

	//check if version exist, if not set default of 0.0.0
	if val, ok := metadata["version"].(string); ok {
		pkg.Version = val
	} else {
		pkg.Version = "0.0.0"
	}

	pkg.Locations = append(pkg.Locations, model.Location{
		Path:      util.TrimUntilLayer(*location),
		LayerHash: location.LayerHash,
	})

	if val, ok := metadata["description"].(string); ok {
		pkg.Description = val
	}

	if val, ok := metadata["license"].(string); ok {
		licenses = append(licenses, val)
	}
	pkg.Licenses = licenses

	//parse CPE
	if val, ok := metadata["author"].(string); ok {
		cpe.NewCPE23(pkg, strings.TrimSpace(val), pkg.Name, pkg.Version)
	} else if val, ok := metadata["authors"].(string); ok {
		cpe.NewCPE23(pkg, strings.TrimSpace(val), pkg.Name, pkg.Version)
	} else {
		cpe.NewCPE23(pkg, pkg.Name, pkg.Name, pkg.Version)
	}

	parseDartPURL(pkg)
	pkg.Metadata = metadata
	*pkgs = append(*pkgs, *pkg)
	return nil
}

// Parse dart packages metadata - lock file
func parseDartPackagesLock(location *model.Location, pkgs *[]model.Package) error {
	byteValue, err := os.ReadFile(location.Path)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal([]byte(byteValue), &dartlockFileMetadata); err != nil {
		return err
	}

	for _, cPackage := range dartlockFileMetadata.Packages {
		pkg := new(model.Package)
		pkg.ID = uuid.NewString()
		pkg.Name = cPackage.Description.Name
		pkg.Version = cPackage.Version
		pkg.Type = pub
		pkg.Path = cPackage.Description.Name
		pkg.Locations = append(pkg.Locations, model.Location{
			Path:      util.TrimUntilLayer(*location),
			LayerHash: location.LayerHash,
		})
		cpe.NewCPE23(pkg, pkg.Name, pkg.Name, pkg.Version)
		parseDartPURL(pkg)
		pkg.Metadata = cPackage
		if pkg.Name != "" {
			*pkgs = append(*pkgs, *pkg)
		}
	}
	return nil
}

// Parse PURL
func parseDartPURL(pkg *model.Package) {
	pkg.PURL = model.PURL("pkg" + ":" + "dart" + "/" + pkg.Name + "@" + pkg.Version)
}
