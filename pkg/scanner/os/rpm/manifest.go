package rpm

import (
	"github.com/carbonetes/diggity/internal/cpe"
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
		cpes := cpe.NewCPE23(component.Name, component.Name, component.Version, Type)
		if len(cpes) > 0 {
			component.CPEs = append(component.CPEs, cpes...)
		}
		stream.AddComponent(component)
	}
}
