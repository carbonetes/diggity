package distro

import (
	"slices"

	"github.com/carbonetes/diggity/internal/logger"
	"github.com/carbonetes/diggity/pkg/stream"
	"github.com/carbonetes/diggity/pkg/types"
)

var (
	Type      = "distro"
	log       = logger.GetLogger()
	Manifests = []string{"etc/os-release"}
)

func Scan(data interface{}) interface{} {
	data, ok := data.(types.ManifestFile)
	if !ok {
		log.Fatal("Distro handler received unknown type")
	}
	distro, err := parseRelease(data.(types.ManifestFile))
	if err != nil {
		log.Fatal(err.Error())
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
