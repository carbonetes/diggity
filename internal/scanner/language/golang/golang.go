package golang

import (
	"path/filepath"
	"slices"

	"github.com/carbonetes/diggity/internal/logger"
	"github.com/carbonetes/diggity/pkg/stream"
	"github.com/carbonetes/diggity/pkg/types"
	"golang.org/x/mod/modfile"
)

const Type string = "golang"

var (
	Manifests = []string{"go.mod"}
	log       = logger.GetLogger()
)

func CheckRelatedFile(file string) (string, bool) {
	if slices.Contains(Manifests, filepath.Base(file)) {
		return Type, true
	}
	return "", false
}

func Scan(data interface{}) interface{} {
	manifest, ok := data.(types.ManifestFile)
	if !ok {
		log.Fatal("Go Modules Handler received unknown type")
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
