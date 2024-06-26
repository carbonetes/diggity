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
	if err != nil {
		return nil // skip
	}

	if len(lockfile) == 0 {
		return nil
	}

	for id, pkg := range lockfile {
		id = strings.TrimSpace(id)
		if strings.Contains(id, ",") {
			pkgs := strings.Split(id, ",")
			if len(pkgs) == 0 {
				continue
			}
			log.Debug("Found multiple packages in Yarn lockfile: ", pkgs)
			for _, p := range pkgs {
				if p == pkg.Resolution {
					continue
				}

				// log.Debug("Scanning Yarn package: ", p)

				name, version := parseYarnPackageName(p), parseYarnVersion(p)
				if len(name) == 0 || len(version) == 0 {
					continue
				}

				metadata := pkg
				metadata.Version = version

				c := component.New(name, pkg.Version, Type)

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

				// log.Debug("Adding Yarn package to CycloneDX BOM: ", c.Name, c.Version)

				// cdx.AddComponent(c, payload.Address)

				components = append(components, *c)
			}
		}

		if pkg.Resolution == "" || pkg.Version == "" {
			continue
		}

		name := parseYarnPackageName(pkg.Resolution)

		var version string
		if pkg.Version == "" {
			version = parseYarnVersion(pkg.Resolution)
		} else {
			version = pkg.Version
		}

		c := component.New(name, version, Type)
		cpes := cpe.NewCPE23(c.Name, c.Name, c.Version, Type)
		if len(cpes) > 0 {
			for _, cpe := range cpes {
				component.AddCPE(c, cpe)
			}
		}

		component.AddOrigin(c, manifest.Path)
		component.AddType(c, Type)

		rawMetadata, err := helper.ToJSON(pkg)
		if err != nil {
			log.Debugf("Error converting metadata to JSON: %s", err)
		}

		if len(rawMetadata) > 0 {
			component.AddRawMetadata(c, rawMetadata)
		}

		if len(payload.Layer) > 0 {
			component.AddLayer(c, payload.Layer)
		}

		if len(payload.Layer) > 0 {
			component.AddLayer(c, payload.Layer)
		}

		// cdx.AddComponent(c, payload.Address)
		components = append(components, *c)
	}

	return &components
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
