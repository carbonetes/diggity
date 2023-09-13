package rpm

import (
	"fmt"
	"strings"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/google/uuid"
	rpmdb "github.com/knqyf263/go-rpmdb/pkg"
)

// Initialize RPM package contents
func initRpmPackage(location *model.Location, rpmPkg *rpmdb.PackageInfo) *model.Package {
	p := model.Package{
		ID:          uuid.NewString(),
		Type:        Type,
		Name:        rpmPkg.Name,
		Version:     fmt.Sprintf("%+v-%+v", rpmPkg.Version, rpmPkg.Release),
		Description: rpmPkg.Summary,
	}

	// Get licenses
	p.Licenses = formatLicenses(rpmPkg.License)

	// Get purl
	p.PURL = model.PURL(fmt.Sprintf("pkg:%s/%s@%s?arch=%s", Type, p.Name, p.Version, rpmPkg.Arch))

	// Set and fill final metadata
	initFinalRpmMetadata(&p, rpmPkg)

	// Format version and get CPEs
	// p.CPEs = []string{cpe.NewCPE23(formatVendor(rpmPkg.Vendor), rpmPkg.Name, cpeVersion)}
	generateCpes(&p, rpmPkg.Name, formatVendor(rpmPkg.Vendor), formatCPEVersion(rpmPkg))
	return &p
}

// Format licenses
func formatLicenses(licensesGroup string) []string {
	if len(licensesGroup) > 0 && licensesGroup != " " {
		licenses := []string{}
		subgroups := strings.Split(licensesGroup, " and ")

		for _, group := range subgroups {
			group = strings.TrimSuffix(strings.TrimPrefix(group, "("), ")")
			licenses = append(licenses, strings.Split(group, " or ")...)
		}

		return deduplicateLicenses(licenses)
	}

	return []string{licensesGroup}
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

// Format CPE version
func formatCPEVersion(rpmPkg *rpmdb.PackageInfo) string {
	if rpmPkg.EpochNum() != 0 {
		return fmt.Sprintf("%+v:%+v", rpmPkg.EpochNum(), rpmPkg.Version)
	}
	return rpmPkg.Version
}

// Deduplicate licenses
func deduplicateLicenses(licenses []string) []string {
	uniqueLicenses := make(map[string]struct{})

	for _, license := range licenses {
		uniqueLicenses[license] = struct{}{}
	}

	deduplicated := make([]string, 0, len(uniqueLicenses))
	for license := range uniqueLicenses {
		deduplicated = append(deduplicated, license)
	}

	return deduplicated
}
