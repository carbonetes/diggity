package cocoapods

import (
	"path/filepath"
	"slices"
	"strings"

	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/pkg/cdx"
	"github.com/carbonetes/diggity/pkg/cdx/component"
	"github.com/carbonetes/diggity/pkg/types"
)

const Type string = "cocoapods"

var Manifests = []string{"Podfile.lock"}

func CheckRelatedFile(file string) (string, bool, bool) {
	if slices.Contains(Manifests, filepath.Base(file)) {
		return Type, true, true
	}
	return "", false, false
}

func Scan(data interface{}) interface{} {
	manifest, ok := data.(types.ManifestFile)
	if !ok {
		log.Error("Cocoapods Handler received unknown type")
	}

	metadata := readManifestFile(manifest.Content)
	for _, pod := range metadata.Pods {
		var pods string
		switch c := pod.(type) {
		case string:
			pods = c
		case map[string]interface{}:
			val := pod.(map[string]interface{})
			for all := range val {
				pods = all
			}
		}

		attributes := strings.Split(pods, " ")
		name, version := attributes[0], strings.TrimSuffix(strings.TrimPrefix(attributes[1], "("), ")")

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

	return data
}
