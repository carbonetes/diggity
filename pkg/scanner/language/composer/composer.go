package composer

import (
	"path/filepath"
	"slices"
	"strings"

	"github.com/CycloneDX/cyclonedx-go"
	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/internal/helper"
	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/pkg/cdx"
	"github.com/carbonetes/diggity/pkg/cdx/component"
	"github.com/carbonetes/diggity/pkg/types"
	"github.com/golistic/urn"
)

const Type string = "composer"

var Manifests = []string{"composer.lock"}

func CheckRelatedFile(file string) (string, bool, bool) {
	if slices.Contains(Manifests, filepath.Base(file)) {
		return Type, true, true
	}
	return "", false, false
}

func Scan(data interface{}) interface{} {
	payload, ok := data.(types.Payload)
	if !ok {
		log.Debug("Composer Handler received unknown type")
		return nil
	}

	scan(payload)

	return data
}

func scan(payload types.Payload) {
	file, ok := payload.Body.(types.ManifestFile)
	if !ok {
		log.Debugf("Failed to convert payload body to manifest file")
		return
	}

	metadata := readManifestFile(file.Content)

	processPackages(metadata.Packages, file.Path, payload.Address)
	processDevPackages(metadata.PackagesDev, file.Path, payload.Layer, payload.Address)
}

func processPackages(packages []ComposerPackage, filePath string, address *urn.URN) {
	for _, pkg := range packages {
		if pkg.Name == "" {
			continue
		}

		c := component.New(pkg.Name, pkg.Version, Type)
		addCPEs(c, c.Name, c.Name, c.Version)
		addCommonAttributes(c, filePath, pkg)
		cdx.AddComponent(c, address)
	}
}

func processDevPackages(packagesDev []ComposerPackage, filePath, layer string, address *urn.URN) {
	for _, pkg := range packagesDev {
		if pkg.Name == "" {
			continue
		}

		c := component.New(pkg.Name, pkg.Version, Type)
		props := strings.Split(pkg.Name, "/")
		vendor, product := props[0], props[1]
		addCPEs(c, vendor, product, pkg.Version)
		addCommonAttributes(c, filePath, pkg)

		if len(layer) > 0 {
			component.AddLayer(c, layer)
		}

		cdx.AddComponent(c, address)
	}
}

func addCPEs(c *cyclonedx.Component, vendor, product, version string) {
	cpes := cpe.NewCPE23(vendor, product, version, Type)
	if len(cpes) > 0 {
		for _, cpe := range cpes {
			component.AddCPE(c, cpe)
		}
	}
}

func addCommonAttributes(c *cyclonedx.Component, filePath string, pkg ComposerPackage) {
	component.AddOrigin(c, filePath)
	component.AddType(c, Type)

	rawMetadata, err := helper.ToJSON(pkg)
	if err != nil {
		log.Debugf("Error converting metadata to JSON: %s", err)
	}

	if len(rawMetadata) > 0 {
		component.AddRawMetadata(c, rawMetadata)
	}

	if len(pkg.License) > 0 {
		for _, license := range pkg.License {
			component.AddLicense(c, license)
		}
	}
}
