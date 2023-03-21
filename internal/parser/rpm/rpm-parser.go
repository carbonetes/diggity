package rpm

/*
========== RPM PARSER ==========
Applicable to OS with RPM as its Package Manager such as:
-CENT OS
-RHEL
-Fedora
-openSUSE
*/

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/internal/file"
	"github.com/carbonetes/diggity/internal/parser/bom"
	"github.com/carbonetes/diggity/internal/parser/util"
	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/model/metadata"

	"strings"

	"github.com/google/uuid"
	rpmdb "github.com/knqyf263/go-rpmdb/pkg"
)

const (
	rpmType = "rpm"
)

var (
	rpmPackagesPath = filepath.Join("rpm", "Packages")
)

// FindRpmPackagesFromContent Find rpm/Packages in the file content.
func FindRpmPackagesFromContent() {
	// Get RPM Information if rpm/Packages is found
	if util.ParserEnabled(rpmType) {
		for _, content := range file.Contents {
			if strings.Contains(content.Path, rpmPackagesPath) {
				if err := readRpmContent(content); err != nil {
					err = errors.New("rpm-parser: " + err.Error())
					bom.Errors = append(bom.Errors, &err)
				}
			}
		}
	}
	defer bom.WG.Done()
}

// Read RPM package information from rpm/Packages
func readRpmContent(location *model.Location) error {

	// Open and Get rpm/Packages data
	db, err := rpmdb.Open(location.Path)
	if err != nil {
		return err
	}
	rpmPkgList, err := db.ListPackages()
	if err != nil {
		return err
	}

	for _, rpmPkg := range rpmPkgList {
		// Get RPM package contents
		_package := new(model.Package)
		_package = initRpmPackage(_package, location, rpmPkg)

		bom.Packages = append(bom.Packages, _package)
	}
	return nil
}

// Initialize RPM package contents
func initRpmPackage(p *model.Package, location *model.Location, rpmPkg *rpmdb.PackageInfo) *model.Package {
	p.ID = uuid.NewString()
	p.Type = rpmType
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
func parseRpmPackageURL(_package *model.Package, architecture string) {
	_package.PURL = model.PURL("pkg" + ":" + rpmType + "/" + _package.Name + "@" + _package.Version + "?arch=" + architecture)
}

// Initialize RPM metadata values from content
func initFinalRpmMetadata(_package *model.Package, rpmPkg *rpmdb.PackageInfo) {
	_package.Metadata = metadata.RPMMetadata{
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

// Format licenses
func formatLicenses(_package *model.Package, licenses string) {
	if len(licenses) > 0 && licenses != " " {
		if strings.Contains(licenses, " and ") {
			_package.Licenses = strings.Split(licenses, " and ")
		} else if strings.Contains(licenses, " or ") {
			_package.Licenses = strings.Split(licenses, " or ")
		} else {
			_package.Licenses = []string{licenses}
		}
	} else {
		_package.Licenses = []string{}
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
