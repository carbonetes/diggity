package portage

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/pkg/cdx"
	"github.com/carbonetes/diggity/pkg/cdx/component"
	"github.com/carbonetes/diggity/pkg/types"
)

const Type string = "gentoo"

var RelatedPath = "var/db/pkg/"

func CheckRelatedFile(file string) (string, bool, bool) {
	if strings.Contains(RelatedPath, file) {
		return Type, true, false
	}
	return "", false, false
}

// TODO: Subject for thorough review and testing
func Scan(data interface{}) interface{} {
	payload, ok := data.(types.Payload)
	if !ok {
		log.Error("Portage Handler received unknown type")
		return nil
	}

	scan(payload)

	return data
}

func scan(payload types.Payload) {
	manifest := payload.Body.(types.ManifestFile)
	
	if len(manifest.Path) == 0 {
		return
	}

	target := filepath.Dir(manifest.Path)
	name, version := parseNameVersion(target)
	if len(name) == 0 || len(version) == 0 {
		return
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

	// no metadata

	cdx.AddComponent(c, payload.Address)
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
