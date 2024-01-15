package distro

import (
	"slices"

	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/pkg/stream"
	"github.com/carbonetes/diggity/pkg/types"
)

var (
	Type      = "distro"
	Manifests = []string{"etc/os-release"}
)

func Scan(data interface{}) interface{} {
	data, ok := data.(types.ManifestFile)
	if !ok {
		log.Error("Distro handler received unknown type")
	}
	distro, err := parseRelease(data.(types.ManifestFile))
	if err != nil {
		log.Error(err.Error())
	}

	stream.SetDistro(distro)
	return data
}

func CheckRelatedFile(file string) (string, bool, bool) {
	if slices.Contains(Manifests, file) {
		return Type, true, true
	}
	return "", false, false
}
