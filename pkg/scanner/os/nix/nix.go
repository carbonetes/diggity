package nix

import (
	"path/filepath"
	"strings"

	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/internal/helper"
	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/pkg/cdx"
	"github.com/carbonetes/diggity/pkg/cdx/component"
	"github.com/carbonetes/diggity/pkg/types"
)

const Type string = "nix"

var RelatedPath = "nix/store/"

func CheckRelatedFile(file string) (string, bool, bool) {
	if strings.Contains(RelatedPath, file) {
		log.Debugf("Found %s file", file)
		return Type, true, false
	}
	return "", false, false
}

func Scan(data interface{}) interface{} {
	manifest, ok := data.(types.ManifestFile)
	if !ok {
		log.Error("Nix Handler received unknown type")
		return nil
	}

	scan(manifest)

	return data
}

func scan(manifest types.ManifestFile) {
	if strings.Contains(filepath.Base(manifest.Path), ".nix") || strings.Contains(filepath.Base(manifest.Path), ".drv") {
		return
	}

	separator := "/"
	if strings.Contains(manifest.Path, "\\") {
		separator = "\\"
	}

	// Get the package name version
	paths := strings.Split(manifest.Path, separator)
	var target string
	for index, path := range paths {
		if path == "nix" {
			if paths[index+1] == "store" && index+2 < len(paths) {
				target = paths[index+2]
				break

			}
		}
	}

	if target == "" {
		return
	}

	// Parse the package name and version
	metadata := parseNixPath(target)
	if metadata == nil {
		return
	}

	if metadata.Name == "" || metadata.Version == "" {
		return
	}

	c := component.New(metadata.Name, metadata.Version, Type)

	cpes := cpe.NewCPE23(c.Name, c.Name, c.Version, Type)
	if len(cpes) > 0 {
		for _, cpe := range cpes {
			component.AddCPE(c, cpe)
		}
	}

	component.AddOrigin(c, manifest.Path)
	component.AddType(c, Type)

	rawMetadata, err := helper.ToJSON(metadata)
	if err != nil {
		log.Errorf("Error converting metadata to JSON: %s", err)
	}

	if len(rawMetadata) > 0 {
		component.AddRawMetadata(c, rawMetadata)
	}

	cdx.AddComponent(c)

}
