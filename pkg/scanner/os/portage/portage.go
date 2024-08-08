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

var (
	RelatedPath = "/db/pkg/"
	RelatedFile = "CONTENTS"
)

func CheckRelatedFile(file string) (string, bool, bool) {
	if strings.Contains(file, RelatedPath) {
		if filepath.Base(file) == RelatedFile {
			return Type, true, true
		}
	}
	return "", false, false
}

// TODO: Subject for thorough review and testing
func Scan(data interface{}) interface{} {
	payload, ok := data.(types.Payload)
	if !ok {
		log.Debug("Portage Handler received unknown type")
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

	if len(file.Path) == 0 {
		return
	}

	target := filepath.Dir(file.Path)
	log.Debugf("Scanning %s", target)
	name, version := parseNameVersion(target)
	log.Debugf("Name: %s, Version: %s", name, version)
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

	component.AddOrigin(c, file.Path)
	component.AddType(c, Type)

	// no metadata

	if len(payload.Layer) > 0 {
		component.AddLayer(c, payload.Layer)
	}

	cdx.AddComponent(c, payload.Address)
}

func parseNameVersion(pkg string) (name string, version string) {
	// parse version
	r := regexp.MustCompile(`[0-9].*`)
	pkgBase := filepath.Base(pkg)
	version = r.FindString(pkgBase)

	// parse name
	name = strings.Replace(filepath.Base(pkg), "-"+version, "", -1)

	return name, version
}
