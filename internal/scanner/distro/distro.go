package distro

import (
	"log"
	"slices"

	"github.com/carbonetes/diggity/pkg/stream"
	"github.com/carbonetes/diggity/pkg/types"
)

var (
	Type      = "distro"
	Manifests = []string{"etc/os-release"}
	Handler   = func(data interface{}) interface{} {
		data, ok := data.(types.ManifestFile)
		if !ok {
			log.Fatal("Distro handler received unknown type")
		}
		distro, err := parseRelease(data.(types.ManifestFile))
		if err != nil {
			log.Fatal(err.Error())
		}

		if distro == nil {
			log.Fatal("Distro handler cannot parse manifest file")
		}

		return data
	}
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

	if distro == nil {
		log.Fatal("Distro handler cannot parse manifest file")
	}
	stream.SetDistro(*distro)
	return data
}

func CheckRelatedFile(file string) (string, bool) {
	if slices.Contains(Manifests, file) {
		return Type, true
	}
	return "", false
}
