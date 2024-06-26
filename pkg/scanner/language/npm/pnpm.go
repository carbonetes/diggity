package npm

import (
	"strconv"
	"strings"

	"github.com/CycloneDX/cyclonedx-go"
	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/pkg/cdx/component"
	"github.com/carbonetes/diggity/pkg/types"
)

func scanPnpmLockfile(payload types.Payload) *[]cyclonedx.Component {
	components := []cyclonedx.Component{}
	manifest := payload.Body.(types.ManifestFile)
	metadata := readPnpmLockfile(manifest.Content)
	if metadata == nil {
		return nil
	}

	for name, info := range metadata.Dependencies {
		if name == "" {
			continue
		}

		var version string
		switch t := info.(type) {
		case string:
			version = strings.SplitN(t, "()", 2)[0]
		case map[string]interface{}:
			val, ok := t["version"].(string)
			if !ok {
				break
			}
			version = strings.SplitN(val, "(", 2)[0]
		}
		if len(version) == 0 {
			continue
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

		if len(payload.Layer) > 0 {
			component.AddLayer(c, payload.Layer)
		}

		components = append(components, *c)
	}

	separator := "/"
	lockfileVersion, _ := strconv.ParseFloat(metadata.LockFileVersion, 64)
	if lockfileVersion >= 6.0 {
		separator = "@"
	}

	for id := range metadata.Packages {
		id = packageNameRegex.ReplaceAllString(id, "$1")
		id = strings.TrimPrefix(id, "/")
		props := strings.Split(id, separator)

		name := strings.Join(props[:len(props)-1], separator)
		version := props[len(props)-1]

		c := component.New(name, version, Type)

		if len(c.Name) == 0 || len(c.Version) == 0 {
			continue
		}

		cpes := cpe.NewCPE23(c.Name, c.Name, c.Version, Type)
		if len(cpes) > 0 {
			for _, cpe := range cpes {
				component.AddCPE(c, cpe)
			}
		}

		component.AddOrigin(c, manifest.Path)
		component.AddType(c, Type)

		if len(payload.Layer) > 0 {
			component.AddLayer(c, payload.Layer)
		}

		components = append(components, *c)
	}

	return &components
}
