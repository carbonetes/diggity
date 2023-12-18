package rpm

import (
	"github.com/carbonetes/diggity/pkg/stream"
	"github.com/carbonetes/diggity/pkg/types"
)

// Read RPM package information from rpm db
func readRpmDb(rpmdb types.RpmDB) {
	if rpmdb.PackageInfos == nil {
		return
	}

	if len(rpmdb.PackageInfos) == 0 {
		return
	}

	for _, pkgInfo := range rpmdb.PackageInfos {
		component := newComponent(pkgInfo)
		if component.Name == "" || component.Version == "" {
			continue
		}
		stream.AddComponent(component)
	}
}
