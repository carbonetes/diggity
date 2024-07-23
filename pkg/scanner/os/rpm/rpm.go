package rpm

import (
	"fmt"
	"slices"
	"strings"

	"github.com/CycloneDX/cyclonedx-go"
	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/internal/helper"
	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/pkg/cdx"
	"github.com/carbonetes/diggity/pkg/cdx/component"
	"github.com/carbonetes/diggity/pkg/cdx/dependency"
	"github.com/carbonetes/diggity/pkg/types"
)

const Type string = "rpm"

var (
	ManifestFiles = []string{"rpm/Packages", "rpm/Packages.db", "rpm/rpmdb.sqlite"}
)

func CheckRelatedFiles(file string) (string, bool, bool) {
	for _, manifest := range ManifestFiles {
		if strings.Contains(file, manifest) {
			return Type, true, true
		}
	}

	return "", false, false
}

func Scan(data interface{}) interface{} {
	payload, ok := data.(types.Payload)
	if !ok {
		log.Debug("RPM Handler received unknown type")
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

		// Remove unnecessary fields
		pkgInfo.BaseNames = nil
		pkgInfo.FileDigests = nil
		pkgInfo.DirNames = nil
		pkgInfo.DirIndexes = nil
		pkgInfo.FileFlags = nil
		pkgInfo.FileModes = nil
		pkgInfo.GroupNames = nil
		pkgInfo.UserNames = nil

		rawMetadata, err := helper.ToJSON(pkgInfo)
		if err != nil {
			log.Debugf("Error converting metadata to JSON: %s", err)
		}

		if len(rawMetadata) > 0 {
			component.AddRawMetadata(c, rawMetadata)
		}

		if len(payload.Layer) > 0 {
			component.AddLayer(c, payload.Layer)
		}

		qm := make(map[string]string)
		if pkgInfo.Arch != "" {
			qm["arch"] = pkgInfo.Arch
		}

		if pkgInfo.SourceRpm != "" {
			qm["upstream"] = pkgInfo.SourceRpm
		}

		if pkgInfo.Vendor != "" {
			c.Publisher = pkgInfo.Vendor
		}

		if pkgInfo.License != "" {
			c.Licenses = &cyclonedx.Licenses{
				{
					License: &cyclonedx.License{
						ID: pkgInfo.License,
					},
				},
			}
		}

		component.AddRefQualifier(c, qm)

		dependencyNode := &cyclonedx.Dependency{
			Ref:          c.BOMRef,
			Dependencies: &[]string{},
		}

		if len(pkgInfo.Requires) > 0 {
			for _, dep := range pkgInfo.Requires {
				for _, p := range rpmdb.PackageInfos {
					for _, r := range p.Provides {
						if strings.Contains(dep, r) {
							if !slices.Contains(*dependencyNode.Dependencies, p.Name) {
								*dependencyNode.Dependencies = append(*dependencyNode.Dependencies, p.Name)
							}
						}
					}
				}
			}
		}

		if len(*dependencyNode.Dependencies) > 0 {
			dependency.AddDependency(payload.Address, dependencyNode)
		}

		cdx.AddComponent(c, payload.Address)
	}
}
