package npm

import (
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/carbonetes/diggity/internal/logger"
	"github.com/carbonetes/diggity/pkg/stream"
	"github.com/carbonetes/diggity/pkg/types"
)

const Type string = "npm"

var (
	Manifests        = []string{"package.json", "package-lock.json", ".package.json", ".package-lock.json", "yarn.lock", "pnpm-lock.yaml"}
	log              = logger.GetLogger()
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

		component := types.NewComponent(metadata.Name, metadata.Version, Type, manifest.Path, metadata.Description, metadata)
		switch metadata.License.(type) {
		case string:
			component.Licenses = append(component.Licenses, metadata.License.(string))
		case map[string]interface{}:
			license := metadata.License.(map[string]interface{})
			if _, ok := license["type"]; ok {
				component.Licenses = append(component.Licenses, license["type"].(string))
			}
		}
		stream.AddComponent(component)
	} else if strings.Contains(manifest.Path, "package-lock.json") {
		metadata := readPackageLockfile(manifest.Content)
		if len(metadata.Dependencies) == 0 {
			return nil
		}
		for name, dependency := range metadata.Dependencies {
			if name == "" || dependency.Version == "" {
				continue
			}
			component := types.NewComponent(name, dependency.Version, Type, manifest.Path, "", dependency)
			stream.AddComponent(component)
		}
	} else if strings.Contains(manifest.Path, "yarn.lock") {
		metadata := readYarnLockfile(manifest.Content)
		for name, info := range metadata.Dependencies {
			if name == "" {
				continue
			}
			component := types.NewComponent(name, info.Version, Type, manifest.Path, "", info)
			component.PURL = "pkg:npm/" + name + "@" + info.Version
			stream.AddComponent(component)
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
			component := types.NewComponent(name, version, Type, manifest.Path, "", info)
			component.PURL = "pkg:npm/" + name + "@" + version
			stream.AddComponent(component)
		}

		separator := "/"
		lockfileVersion, err := strconv.ParseFloat(metadata.LockFileVersion, 64)
		if err != nil {
			log.Errorf("Error parsing lockfile version: %s", err)
		}

		if lockfileVersion >= 6.0 {
			separator = "@"
		}

		for id, pkg := range metadata.Packages {
			id = packageNameRegex.ReplaceAllString(id, "$1")
			id = strings.TrimPrefix(id, "/")
			props := strings.Split(id, separator)

			name := strings.Join(props[:len(props)-1], separator)
			version := props[len(props)-1]

			component := types.NewComponent(name, version, Type, manifest.Path, "", pkg)
			stream.AddComponent(component)
		}
	}

	return data
}
