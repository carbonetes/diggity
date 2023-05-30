package rpm

import (
	"fmt"
	"strings"

	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/parser/util"
	"github.com/google/uuid"
	rpmdb "github.com/knqyf263/go-rpmdb/pkg"
)

// Initialize RPM package contents
func initRpmPackage(p *model.Package, location *model.Location, rpmPkg *rpmdb.PackageInfo) *model.Package {
	p.ID = uuid.NewString()
	p.Type = Type
	p.Path = rpmPackagesPath
	p.Name = rpmPkg.Name
	p.Version = fmt.Sprintf("%+v-%+v", rpmPkg.Version, rpmPkg.Release)
	p.Description = rpmPkg.Summary

	// get licenses
	formatLicenses(p, rpmPkg.License)

	// get locations
	p.Locations = append(p.Locations, model.Location{
		Path:      util.TrimUntilLayer(*location),
		LayerHash: location.LayerHash,
	})

	// get purl
	parseRpmPackageURL(p, rpmPkg.Arch)

	// set and fill final metadata
	initFinalRpmMetadata(p, rpmPkg)
	// p.Metadata = rpmPkg

	// format version
	var cpeVersion string
	if rpmPkg.EpochNum() != 0 {
		cpeVersion = fmt.Sprintf("%+v\\:%+v", rpmPkg.EpochNum(), p.Version)
		p.Version = fmt.Sprintf("%+v:%+v", rpmPkg.EpochNum(), p.Version)
	} else {
		cpeVersion = p.Version
	}

	// get CPEs
	cpe.NewCPE23(p, formatVendor(rpmPkg.Vendor), rpmPkg.Name, cpeVersion)

	return p
}

// Parse PURL
func parseRpmPackageURL(pkg *model.Package, architecture string) {
	pkg.PURL = model.PURL("pkg" + ":" + Type + "/" + pkg.Name + "@" + pkg.Version + "?arch=" + architecture)
}

// Format licenses
func formatLicenses(pkg *model.Package, licenses string) {
	if len(licenses) > 0 && licenses != " " {
		if strings.Contains(licenses, " and ") {
			pkg.Licenses = strings.Split(licenses, " and ")
		} else if strings.Contains(licenses, " or ") {
			pkg.Licenses = strings.Split(licenses, " or ")
		} else {
			pkg.Licenses = []string{licenses}
		}
	} else {
		pkg.Licenses = []string{}
	}
}

// Format vendor for CPEs
func formatVendor(vendor string) string {
	switch vendor {
	case "CentOS":
		return "centos"
	case "Red Hat, Inc.":
		return "redhat"
	case "Fedora Project":
		return "fedoraproject"
	default:
		return strings.TrimSpace(strings.ToLower(vendor))
	}
}
