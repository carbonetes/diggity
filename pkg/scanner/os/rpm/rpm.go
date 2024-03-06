package rpm

import (
	"fmt"
	"path/filepath"
	"slices"

	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/pkg/cdx"
	"github.com/carbonetes/diggity/pkg/cdx/component"
	"github.com/carbonetes/diggity/pkg/types"
)

const Type string = "rpm"

var (
	ManifestFiles = []string{"Packages", "Packages.db", "rpmdb.sqlite"}
	RelatedPaths  = []string{"var\\lib\\rpm", "usr\\lib\\rpm", "etc\\rpm"}
)

func Scan(data interface{}) interface{} {
	rpmdb, ok := data.(types.RpmDB)

	if !ok {
		log.Error("RPM Handler received unknown type")
	}

	if len(rpmdb.PackageInfos) == 0 {
		return nil
	}

	for _, pkgInfo := range rpmdb.PackageInfos {

		if pkgInfo.Name == "" || pkgInfo.Version == "" {
			continue
		}

		version := fmt.Sprintf("%+v-%+v", pkgInfo.Version, pkgInfo.Release)

		c := component.New(pkgInfo.Name, version, Type)

		cpes := cpe.NewCPE23(c.Name, c.Name, version, Type)
		if len(cpes) > 0 {
			for _, cpe := range cpes {
				component.AddCPE(c, cpe)
			}
		}

		component.AddOrigin(c, rpmdb.Path)
		component.AddType(c, Type)

		licenses := formatLicenses(pkgInfo.License)

		if len(licenses) > 0 {
			for _, license := range licenses {
				component.AddLicense(c, license)
			}
		}

		cdx.AddComponent(c)
	}

	return data
}

func CheckRelatedFiles(file string) (string, bool, bool) {
	if slices.Contains(RelatedPaths, filepath.Dir(file)) {
		if slices.Contains(ManifestFiles, filepath.Base(file)) {
			return Type, true, true
		}

	}
	return "", false, false
}
