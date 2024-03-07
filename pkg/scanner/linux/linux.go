package linux

import (
	"slices"

	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/pkg/types"
)

var (
	Releases []types.OSRelease
	Type     = "osrelease"
	// Add more os release files here if needed
	Manifests = []string{"etc/os-release", "usr/lib/os-release", "etc/lsb-release", "etc/centos-release", "etc/redhat-release", "etc/debian_version", "etc/alpine-release", "etc/SuSE-release", "etc/gentoo-release", "etc/arch-release", "etc/oracle-release"}
)

func Scan(data interface{}) interface{} {
	data, ok := data.(types.ManifestFile)
	if !ok {
		log.Error("Distro handler received unknown type")
	}

	Releases = append(Releases, parse(data.(types.ManifestFile)))

	return data
}

func CheckRelatedFile(file string) (string, bool, bool) {
	if slices.Contains(Manifests, file) {
		return Type, true, true
	}
	return "", false, false
}
