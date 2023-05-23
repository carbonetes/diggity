package rpm

import (
	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/model/metadata"
	rpmdb "github.com/knqyf263/go-rpmdb/pkg"
)

// Initialize RPM metadata values from content
func initFinalRpmMetadata(pkg *model.Package, rpmPkg *rpmdb.PackageInfo) {
	pkg.Metadata = metadata.RPMMetadata{
		Release:         rpmPkg.Release,
		Architecture:    rpmPkg.Arch,
		SourceRpm:       rpmPkg.SourceRpm,
		License:         rpmPkg.License,
		Size:            rpmPkg.Size,
		Name:            rpmPkg.Name,
		PGP:             rpmPkg.PGP,
		ModularityLabel: rpmPkg.Modularitylabel,
		Summary:         rpmPkg.Summary,
		Vendor:          rpmPkg.Vendor,
		Version:         rpmPkg.Version,
		Epoch:           rpmPkg.EpochNum(),
		DigestAlgorithm: rpmPkg.DigestAlgorithm.String(),
	}
}
