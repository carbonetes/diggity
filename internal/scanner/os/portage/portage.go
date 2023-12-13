package portage

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/carbonetes/diggity/internal/logger"
	"github.com/carbonetes/diggity/pkg/stream"
	"github.com/carbonetes/diggity/pkg/types"
)

const Type string = "gentoo"

var (
	RelatedPath = "var/db/pkg/"
	log         = logger.GetLogger()
)

func CheckRelatedFile(file string) (string, bool, bool) {
	if strings.Contains(RelatedPath, file) {
		return Type, true, false
	}
	return "", false, false
}

// TODO: Subject for thorough review and testing
func Scan(data interface{}) interface{} {
	manifest, ok := data.(types.ManifestFile)
	if !ok {
		log.Fatal("Portage Handler received unknown type")
	}

	if len(manifest.Path) == 0 {
		return nil
	}

	target := filepath.Dir(manifest.Path)
	name, version := parseNameVersion(target)
	if len(name) == 0 || len(version) == 0 {
		return nil
	}

	component := types.NewComponent(name, version, Type, manifest.Path, "", nil)
	stream.AddComponent(component)

	return data
}

func parseNameVersion(pkg string) (name string, version string) {
	// parse version
	r := regexp.MustCompile(`[0-9].*`)
	pkgBase := filepath.Base(pkg)
	version = r.FindString(pkgBase)

	// parse name
	namePath := strings.Split(pkg, RelatedPath)[1]
	name = strings.Replace(namePath, "-"+version, "", -1)

	return name, version
}
