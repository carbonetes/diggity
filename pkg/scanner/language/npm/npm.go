package npm

import (
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/internal/helper"
	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/pkg/cdx"
	"github.com/carbonetes/diggity/pkg/cdx/component"
	"github.com/carbonetes/diggity/pkg/types"
	"github.com/hashicorp/go-version"
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
	payload, ok := data.(types.Payload)
	if !ok {
		log.Error("NPM Handler received unknown type")
		return nil
	}

	scan(payload)

	return data
}

func scan(payload types.Payload) {
	manifest := payload.Body.(types.ManifestFile)
	if strings.Contains(manifest.Path, "package.json") {
		metadata := readManifestFile(manifest.Content)

		devDependencies := metadata.DevDependencies
		if len(devDependencies) > 0 {
			for name, version := range devDependencies {
				n := parseYarnPackageName(name)
				v := cleanVersion(version.(string))
				if n == "" || v == "" {
					continue
				}

				if !validateVersion(v) {
					continue
				}

				c := component.New(n, v, Type)

				cpes := cpe.NewCPE23(c.Name, c.Name, c.Version, Type)
				if len(cpes) > 0 {
					for _, cpe := range cpes {
						component.AddCPE(c, cpe)
					}
				}

				component.AddOrigin(c, manifest.Path)
				component.AddType(c, Type)

				cdx.AddComponent(c, payload.Address)
			}
		}
		dependencies := metadata.Dependencies
		if len(dependencies) > 0 {
			for name, version := range dependencies {
				n := parseYarnPackageName(name)
				v := cleanVersion(version.(string))
				if n == "" || v == "" {
					continue
				}

				if !validateVersion(v) {
					continue
				}

				c := component.New(n, v, Type)

				cpes := cpe.NewCPE23(c.Name, c.Name, c.Version, Type)
				if len(cpes) > 0 {
					for _, cpe := range cpes {
						component.AddCPE(c, cpe)
					}
				}

				component.AddOrigin(c, manifest.Path)
				component.AddType(c, Type)

				cdx.AddComponent(c, payload.Address)
			}
		}

		if metadata.Name == "" || metadata.Version == "" {
			return
		}

		n := cleanName(metadata.Name)
		v := cleanVersion(metadata.Version)

		c := component.New(n, v, Type)

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

		rawMetadata, err := helper.ToJSON(metadata)
		if err != nil {
			log.Errorf("Error converting metadata to JSON: %s", err)
		}

		if len(rawMetadata) > 0 {
			component.AddRawMetadata(c, rawMetadata)
		}

		cdx.AddComponent(c, payload.Address)

	} else if strings.Contains(manifest.Path, "package-lock.json") {
		metadata := readPackageLockfile(manifest.Content)
		if len(metadata.Dependencies) == 0 {
			return
		}
		for name, dependency := range metadata.Dependencies {
			if name == "" || dependency.Version == "" {
				continue
			}

			c := component.New(parseYarnPackageName(name), cleanVersion(dependency.Version), Type)

			cpes := cpe.NewCPE23(c.Name, c.Name, c.Version, Type)
			if len(cpes) > 0 {
				for _, cpe := range cpes {
					component.AddCPE(c, cpe)
				}
			}

			component.AddOrigin(c, manifest.Path)
			component.AddType(c, Type)

			cdx.AddComponent(c, payload.Address)
		}
	} else if strings.Contains(manifest.Path, "yarn.lock") {
		lockfile, err := parseYarnLock(manifest.Content)
		if err != nil {
			log.Errorf("Error parsing yarn.lock: %s", err)
		}
		for id, pkg := range lockfile {
			if strings.Contains(id, ",") {
				pkgs := strings.Split(id, ",")
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
						log.Errorf("Error converting metadata to JSON: %s", err)
					}

					if len(rawMetadata) > 0 {
						component.AddRawMetadata(c, rawMetadata)
					}

					cdx.AddComponent(c, payload.Address)
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
				log.Errorf("Error converting metadata to JSON: %s", err)
			}

			if len(rawMetadata) > 0 {
				component.AddRawMetadata(c, rawMetadata)
			}

			cdx.AddComponent(c, payload.Address)
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

			cdx.AddComponent(c, payload.Address)
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

			cdx.AddComponent(c, payload.Address)
		}
	}
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

func cleanName(name string) string {
	return strings.TrimSpace(strings.ReplaceAll(name, "@", ""))
}

func cleanVersion(version string) string {
	version = strings.TrimSpace(strings.ReplaceAll(version, "^", ""))
	version = strings.TrimSpace(strings.ReplaceAll(version, "~", ""))
	return version
}

func validateVersion(v string) bool {
	_, err := version.NewVersion(v)
	return err == nil
}
