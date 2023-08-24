package cocoapods

import (
	"strings"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/model/metadata"
	"github.com/google/uuid"
)

func newPackage(pod interface{}, podFileLockFileMetadata *metadata.PodFileLockMetadata) *model.Package {
	var pkg model.Package
	pkg.ID = uuid.NewString()

	//check metadata in podfile lock file
	var pods string
	switch cp := pod.(type) {
	case string:
		pods = cp
	case map[string]interface{}:
		podVal := pod.(map[string]interface{})
		for podsAll := range podVal {
			pods = podsAll
		}
	}

	splits := strings.Split(pods, " ")
	name := splits[0]
	version := strings.TrimSuffix(strings.TrimPrefix(splits[1], "("), ")")

	pkg.Name = name
	pkg.Version = version
	pkg.Type = Type
	pkg.Path = name
	basepodname := strings.Split(name, "/")[0]
	generateCpes(&pkg)
	setPurl(&pkg)
	if val, ok := podFileLockFileMetadata.SpecChecksums[basepodname]; ok {
		pkg.Metadata = metadata.PodFileLockMetadataCheckSums{
			Checksums: val,
		}
	}
	return &pkg
}

// Parse PURL
func setPurl(pkg *model.Package) {
	pkg.PURL = model.PURL("pkg" + ":" + "cocoapods" + "/" + pkg.Name + "@" + pkg.Version)
}
