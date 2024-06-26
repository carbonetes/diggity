package npm

import (
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/CycloneDX/cyclonedx-go"
	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/pkg/cdx"
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
		log.Debug("NPM Handler received unknown type")
		return nil
	}

	scan(payload)

	return data
}

func scan(payload types.Payload) {
	manifest := payload.Body.(types.ManifestFile)
	log.Debug("Scanning NPM manifest file: ", manifest.Path)

	components := []cyclonedx.Component{}
	switch filepath.Base(manifest.Path) {
	case "package.json":
		result := scanPackageJSON(payload)
		if result != nil && len(*result) > 0 {
			components = append(components, *result...)
		}
	case "package-lock.json":
		result := scanPackageLockfile(payload)
		if result != nil && len(*result) > 0 {
			components = append(components, *result...)
		}
	case "yarn.lock":
		result := scanYarnLockile(payload)
		if result != nil && len(*result) > 0 {
			components = append(components, *result...)
		}
	case "pnpm-lock.yaml":
		result := scanPnpmLockfile(payload)
		if result != nil && len(*result) > 0 {
			components = append(components, *result...)
		}
	}

	if len(components) > 0 {
		cdx.AddComponents(&components, payload.Address)
	}

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
