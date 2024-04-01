package golang

import (
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/internal/helper"
	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/pkg/cdx"
	"github.com/carbonetes/diggity/pkg/cdx/component"
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
		goBinary, ok := data.(types.GoBinary)
		if !ok {
			log.Error("Go Modules Handler received unknown type")
			return nil
		}

		scanBinary(goBinary)

		return data
	}

	scan(manifest)

	return data
}

func scan(manifest types.ManifestFile) {
	modFile := readManifestFile(manifest.Content, manifest.Path)
	for _, pkg := range modFile.Require {
		if pkg.Mod.Path == "" || pkg.Mod.Version == "" {
			continue
		}
		if checkIfExcluded(modFile.Exclude, pkg.Mod.Path) {
			continue
		}

		c := component.New(pkg.Mod.Path, pkg.Mod.Version, Type)

		cpes := GenerateCpes(c.Version, SplitPath(c.Name))
		if len(cpes) > 0 {
			for _, cpe := range cpes {
				component.AddCPE(c, cpe)
			}
		}

		rawMetadata, err := helper.ToJSON(pkg)
		if err != nil {
			log.Errorf("Error converting metadata to JSON: %s", err)
		}

		component.AddOrigin(c, manifest.Path)
		component.AddType(c, Type)
		component.AddRawMetadata(c, rawMetadata)

		cdx.AddComponent(c)
	}

	for _, pkg := range modFile.Replace {
		if pkg.New.Path == "" || pkg.New.Version == "" {
			continue
		}
		if checkIfExcluded(modFile.Exclude, pkg.New.Path) {
			continue
		}

		c := component.New(pkg.New.Path, pkg.New.Version, Type)

		cpes := GenerateCpes(c.Version, SplitPath(c.Name))
		if len(cpes) > 0 {
			for _, cpe := range cpes {
				component.AddCPE(c, cpe)
			}
		}

		component.AddOrigin(c, manifest.Path)
		component.AddType(c, Type)

		rawMetadata, err := helper.ToJSON(pkg)
		if err != nil {
			log.Errorf("Error converting metadata to JSON: %s", err)
		}

		if len(rawMetadata) > 0 {
			component.AddRawMetadata(c, rawMetadata)
		}

		cdx.AddComponent(c)
	}
}

func scanBinary(goBinary types.GoBinary) {
	buildInfo := goBinary.BuildInfo

	for _, s := range buildInfo.Settings {
		// locate version
		v := parseVersion(s.Value)
		if v != "" {
			c := component.New(goBinary.File, v, Type)

			cpes := GenerateCpes(c.Version, SplitPath(c.Name))
			if len(cpes) > 0 {
				for _, cpe := range cpes {
					component.AddCPE(c, cpe)
				}
			}

			component.AddOrigin(c, goBinary.Path)
			component.AddType(c, Type)

			rawMetadata, err := helper.ToJSON(s)
			if err != nil {
				log.Errorf("Error converting metadata to JSON: %s", err)
			}

			if len(rawMetadata) > 0 {
				component.AddRawMetadata(c, rawMetadata)
			}

			cdx.AddComponent(c)
			break
		}
	}

	for _, dep := range buildInfo.Deps {
		if dep.Path == "" || dep.Version == "" {
			continue
		}

		c := component.New(dep.Path, dep.Version, Type)

		cpes := GenerateCpes(c.Version, SplitPath(c.Name))
		if len(cpes) > 0 {
			for _, cpe := range cpes {
				component.AddCPE(c, cpe)
			}
		}

		component.AddOrigin(c, goBinary.Path)
		component.AddType(c, Type)

		rawMetadata, err := helper.ToJSON(dep)
		if err != nil {
			log.Errorf("Error converting metadata to JSON: %s", err)
		}

		if len(rawMetadata) > 0 {
			component.AddRawMetadata(c, rawMetadata)
		}

		cdx.AddComponent(c)
	}
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

var pattern = regexp.MustCompile(`\d+\.\d+\.\d+`)

func parseVersion(s string) string {
	return pattern.FindString(s)
}
