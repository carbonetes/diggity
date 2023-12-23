package nix

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/pkg/stream"
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

	if strings.Contains(filepath.Base(manifest.Path), ".nix") || strings.Contains(filepath.Base(manifest.Path), ".drv") {
		return nil
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
		return nil
	}

	// Parse the package name and version
	metadata := parseNixPath(target)
	if metadata == nil {
		return nil
	}

	if metadata.Name == "" || metadata.Version == "" {
		return nil
	}

	component := types.NewComponent(metadata.Name, metadata.Version, Type, manifest.Path, "", metadata)
	component.PURL = fmt.Sprintf("pkg:nix/%s@%s", metadata.Name, metadata.Version)
	cpes := cpe.NewCPE23(component.Name, component.Name, component.Version, Type)
	if len(cpes) > 0 {
		component.CPEs = append(component.CPEs, cpes...)
	}
	stream.AddComponent(component)

	return data
}
