package npm

import (
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/pkg/cdx"
	"github.com/carbonetes/diggity/pkg/cdx/component"
	"github.com/carbonetes/diggity/pkg/types"
)

const Type string = "npm"

var (
	Manifests        = []string{"package.json", "package-lock.json", ".package.json", ".package-lock.json", "yarn.lock", "pnpm-lock.yaml"}
	packageNameRegex = regexp.MustCompile(`^/?([^(]*)(?:\(.*\))*$`)
)

func CheckRelatedFile(file string) (string, bool, bool) {
	if slices.Contains(Manifests, filepath.Base(file)) {
		return Type, true, true
	}
	return "", false, false
}

func Scan(data interface{}) interface{} {
	manifest, ok := data.(types.ManifestFile)
	if !ok {
		log.Error("NPM Handler received unknown type")
	}

	if strings.Contains(manifest.Path, "package.json") {
		metadata := readManifestFile(manifest.Content)
		if metadata.Name == "" || metadata.Version == "" {
			return nil
		}

		c := component.New(metadata.Name, metadata.Version, Type)

		cpes := cpe.NewCPE23(c.Name, c.Name, c.Version, Type)
		if len(cpes) > 0 {
			for _, cpe := range cpes {
				component.AddCPE(c, cpe)
			}
		}

		component.AddOrigin(c, manifest.Path)
		component.AddType(c, Type)

		switch metadata.License.(type) {
		case string:
			component.AddLicense(c, metadata.License.(string))
		case map[string]interface{}:
			license := metadata.License.(map[string]interface{})
			if _, ok := license["type"]; ok {
				component.AddLicense(c, license["type"].(string))
			}
		}

		cdx.AddComponent(c)

	} else if strings.Contains(manifest.Path, "package-lock.json") {
		metadata := readPackageLockfile(manifest.Content)
		if len(metadata.Dependencies) == 0 {
			return nil
		}
		for name, dependency := range metadata.Dependencies {
			if name == "" || dependency.Version == "" {
				continue
			}

			c := component.New(name, dependency.Version, Type)

			cpes := cpe.NewCPE23(c.Name, c.Name, c.Version, Type)
			if len(cpes) > 0 {
				for _, cpe := range cpes {
					component.AddCPE(c, cpe)
				}
			}

			component.AddOrigin(c, manifest.Path)
			component.AddType(c, Type)

			cdx.AddComponent(c)
		}
	} else if strings.Contains(manifest.Path, "yarn.lock") {
		packages, err := ParseYarnLock(manifest.Content)
		if err != nil {
			log.Errorf("Error parsing yarn.lock: %s", err)
		}
		for name, info := range packages {
			if name == "" {
				continue
			}

			c := component.New(name, info.Version, Type)

			cpes := cpe.NewCPE23(c.Name, c.Name, c.Version, Type)
			if len(cpes) > 0 {
				for _, cpe := range cpes {
					component.AddCPE(c, cpe)
				}
			}

			component.AddOrigin(c, manifest.Path)
			component.AddType(c, Type)

			cdx.AddComponent(c)
		}
	} else if strings.Contains(manifest.Path, "pnpm-lock.yaml") {
		metadata := readPnpmLockfile(manifest.Content)
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

			cdx.AddComponent(c)
		}

		separator := "/"
		lockfileVersion, err := strconv.ParseFloat(metadata.LockFileVersion, 64)
		if err != nil {
			log.Errorf("Error parsing lockfile version: %s", err)
		}

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

			cpes := cpe.NewCPE23(c.Name, c.Name, c.Version, Type)
			if len(cpes) > 0 {
				for _, cpe := range cpes {
					component.AddCPE(c, cpe)
				}
			}

			component.AddOrigin(c, manifest.Path)
			component.AddType(c, Type)

			cdx.AddComponent(c)
		}
	}

	return data
}
