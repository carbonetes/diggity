package python

import (
	"strings"

	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/parser/util"
	"github.com/google/uuid"
)

// Initialize python package
func initPythonPackages(metadata map[string]interface{}, location *model.Location) *model.Package {
	p := new(model.Package)
	p.ID = uuid.NewString()
	p.Type = Type
	p.Locations = append(p.Locations, model.Location{
		Path:      util.TrimUntilLayer(*location),
		LayerHash: location.LayerHash,
	})

	// parse name and version based on metadata
	if _, ok := metadata["Name"]; ok {
		p.Name = metadata["Name"].(string)
		p.Version = metadata["Version"].(string)
		p.Path = metadata["Name"].(string)
	} else {
		p.Name = metadata["name"].(string)
		p.Version = metadata["version"].(string)
		p.Path = metadata["name"].(string)
	}

	// check first if description exist in metadata
	if val, ok := metadata["description"].(string); ok {
		p.Description = val
	}

	// check first if license exist in metadata
	if val, ok := metadata["License"]; ok {
		p.Licenses = append(p.Licenses, val.(string))
	} else {
		p.Licenses = []string{}
	}

	p.Type = Type

	// parse PURL
	setPurl(p)
	filesPath := strings.Split(location.Path, pythonPackage)[0]
	filesPath = filesPath + pythonRecord
	err := parseMetadataFiles(metadata, filesPath)
	if _, ok := metadata["Files"]; ok && err == nil {
		tmpLocation := new(model.Location)
		tmpLocation.LayerHash = location.LayerHash
		tmpLocation.Path = filesPath
		p.Locations = append(p.Locations, model.Location{
			Path:      util.TrimUntilLayer(*tmpLocation),
			LayerHash: location.LayerHash,
		})
	}
	p.Metadata = metadata

	// parse CPE
	if val, ok := metadata["Author"].(string); ok {
		if val == unknownField {
			val = p.Name
		}
		cpe.NewCPE23(p, strings.TrimSpace(val), p.Name, p.Version)
	} else {
		cpe.NewCPE23(p, p.Name, p.Name, p.Version)
	}

	return p
}

// Parse PURL
func setPurl(pkg *model.Package) {
	pkg.PURL = model.PURL("pkg" + ":" + pypi + "/" + pkg.Name + "@" + pkg.Version)
}
