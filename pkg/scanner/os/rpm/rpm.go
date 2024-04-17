package rpm

import (
	"fmt"
	"path/filepath"
	"slices"

	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/internal/helper"
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

func CheckRelatedFiles(file string) (string, bool, bool) {
	if slices.Contains(RelatedPaths, filepath.Dir(file)) {
		if slices.Contains(ManifestFiles, filepath.Base(file)) {
			return Type, true, true
		}

	}
	return "", false, false
}

func Scan(data interface{}) interface{} {
	payload, ok := data.(types.Payload)
	if !ok {
		log.Error("RPM Handler received unknown type")
		return nil
	}

	scan(payload)

	return data
}

func scan(payload types.Payload) {
	rpmdb := payload.Body.(types.RpmDB)

	if len(rpmdb.PackageInfos) == 0 {
		return
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

		rawMetadata, err := helper.ToJSON(pkgInfo)
		if err != nil {
			log.Errorf("Error converting metadata to JSON: %s", err)
		}

		if len(rawMetadata) > 0 {
			component.AddRawMetadata(c, rawMetadata)
		}

		cdx.AddComponent(c, payload.Address)
	}
}
