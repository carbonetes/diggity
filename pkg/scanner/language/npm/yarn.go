package npm

import (
	"strings"

	"github.com/CycloneDX/cyclonedx-go"
	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/internal/helper"
	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/pkg/cdx/component"
	"github.com/carbonetes/diggity/pkg/types"
)

func scanYarnLockile(payload types.Payload) *[]cyclonedx.Component {
	components := []cyclonedx.Component{}
	manifest := payload.Body.(types.ManifestFile)
	lockfile, err := parseYarnLock(manifest.Content)
	if err != nil || len(lockfile) == 0 {
		return nil // skip
	}

	for id, pkg := range lockfile {
		id = strings.TrimSpace(id)
		if strings.Contains(id, ",") {
			processMultiplePackages(id, pkg, manifest, payload, &components)
		} else {
			processSinglePackage(pkg, manifest, payload, &components)
		}
	}

	return &components
}

func processMultiplePackages(id string, pkg YarnLockfile, manifest types.ManifestFile, payload types.Payload, components *[]cyclonedx.Component) {
	pkgs := strings.Split(id, ",")
	if len(pkgs) == 0 {
		return
	}

	for _, p := range pkgs {
		if p == pkg.Resolution {
			continue
		}

		name, version := parseYarnPackageName(p), parseYarnVersion(p)
		if len(name) == 0 || len(version) == 0 {
			continue
		}

		metadata := pkg
		metadata.Version = version

		c := createYarnComponent(name, pkg.Version, metadata, manifest, payload)
		*components = append(*components, *c)
	}
}

func processSinglePackage(pkg YarnLockfile, manifest types.ManifestFile, payload types.Payload, components *[]cyclonedx.Component) {
	if pkg.Resolution == "" || pkg.Version == "" {
		return
	}

	name := parseYarnPackageName(pkg.Resolution)
	version := pkg.Version
	if version == "" {
		version = parseYarnVersion(pkg.Resolution)
	}

	c := createYarnComponent(name, version, pkg, manifest, payload)
	if len(c.Name) == 0 || len(c.Version) == 0 {
		return
	}

	*components = append(*components, *c)
}

func createYarnComponent(name, version string, metadata YarnLockfile, manifest types.ManifestFile, payload types.Payload) *cyclonedx.Component {
	c := component.New(name, version, Type)

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
		log.Debugf("Error converting metadata to JSON: %s", err)
	}

	if len(rawMetadata) > 0 {
		component.AddRawMetadata(c, rawMetadata)
	}

	if len(payload.Layer) > 0 {
		component.AddLayer(c, payload.Layer)
	}

	return c
}

func parseYarnPackageName(name string) string {
	if !strings.Contains(name, "@") {
		return name
	}

	parts := strings.Split(name, "@")
	if len(parts) == 2 {
		return strings.TrimSpace(strings.Split(name, "@")[0])
	}

	if len(parts) == 3 {
		return strings.TrimSpace(strings.Split(name, "@")[1])
	}

	return name
}

func parseYarnVersion(s string) string {
	if !strings.Contains(s, ":") {
		return s
	}

	parts := strings.Split(s, ":")
	if len(parts) == 2 {
		return strings.TrimSpace(parts[1])
	}

	return strings.TrimSpace(strings.ReplaceAll(parts[len(parts)-1], "^", ""))
}
