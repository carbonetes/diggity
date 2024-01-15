package golang

import (
	"path/filepath"
	"slices"
	"strings"

	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/pkg/stream"
	"github.com/carbonetes/diggity/pkg/types"
	"golang.org/x/mod/modfile"
)

const Type string = "golang"

var Manifests = []string{"go.mod"}

func CheckRelatedFile(file string) (string, bool, bool) {
	if slices.Contains(Manifests, filepath.Base(file)) {
		return Type, true, true
	}
	return "", false, false
}

func Scan(data interface{}) interface{} {
	manifest, ok := data.(types.ManifestFile)
	if !ok {
		log.Error("Go Modules Handler received unknown type")
	}

	modFile := readManifestFile(manifest.Content, manifest.Path)
	for _, pkg := range modFile.Require {
		if pkg.Mod.Path == "" || pkg.Mod.Version == "" {
			continue
		}
		if checkIfExcluded(modFile.Exclude, pkg.Mod.Path) {
			continue
		}
		component := types.NewComponent(pkg.Mod.Path, pkg.Mod.Version, Type, manifest.Path, "", pkg)
		cpes := GenerateCpes(component.Version, SplitPath(component.Name))
		if len(cpes) > 0 {
			component.CPEs = append(component.CPEs, cpes...)
		}
		stream.AddComponent(component)
	}

	for _, pkg := range modFile.Replace {
		if pkg.New.Path == "" || pkg.New.Version == "" {
			continue
		}
		if checkIfExcluded(modFile.Exclude, pkg.New.Path) {
			continue
		}
		component := types.NewComponent(pkg.New.Path, pkg.New.Version, Type, manifest.Path, "", pkg)
		cpes := GenerateCpes(component.Version, SplitPath(component.Name))
		if len(cpes) > 0 {
			component.CPEs = append(component.CPEs, cpes...)
		}
		stream.AddComponent(component)
	}

	return data
}

func checkIfExcluded(excludeList []*modfile.Exclude, name string) bool {
	for _, exclude := range excludeList {
		if exclude.Mod.Path == name {
			return true
		}
	}
	return false
}

// Split go package path
func SplitPath(path string) []string {
	return strings.Split(path, "/")
}

func GenerateCpes(version string, paths []string) []string {
	var cpes []string
	// check if cpePaths only contains the product
	if len(paths) > 1 {
		cpes = cpe.NewCPE23(paths[len(paths)-2], paths[len(paths)-1], FormatVersion(version), Type)
	} else {
		cpes = cpe.NewCPE23("", paths[0], FormatVersion(version), Type)
	}
	return cpes
}

// Format Version String
func FormatVersion(version string) string {
	if strings.Contains(version, "(") && strings.Contains(version, ")") {
		version = strings.Replace(version, "(", "", -1)
		version = strings.Replace(version, ")", "", -1)
	}
	return version
}
