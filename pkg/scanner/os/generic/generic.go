package generic

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/pkg/stream"
	"github.com/carbonetes/diggity/pkg/types"
)

var (
	RelatedPath = "usr/bin/"
)

const Type string = "binary"

func CheckRelatedFile(file string) (string, bool, bool) {
	if strings.Contains(file, RelatedPath) {
		return Type, true, true
	}

	return "", false, false
}

func Scan(data interface{}) interface{} {
	generic, ok := data.(types.Generic)
	if !ok {
		return nil
	}

	ro := generic.ROData

	// Look for a string that matches a version format.
	versionRegex := regexp.MustCompile(`\d+\.\d+\.\d+(-[0-9A-Za-z-]+(\.[0-9A-Za-z-]+)*)?(\+[0-9A-Za-z-]+(\.[0-9A-Za-z-]+)*)?`)
	for _, s := range ro {
		version := versionRegex.FindString(s)
		if version != "" {
			component := types.NewComponent(filepath.Base(generic.Name), version, Type, generic.Path, "", nil)
			cpes := cpe.NewCPE23(component.Name, component.Name, component.Version, Type)
			if len(cpes) > 0 {
				component.CPEs = append(component.CPEs, cpes...)
			}
			stream.AddComponent(component)
			return nil
		}
	}

	return data
}
