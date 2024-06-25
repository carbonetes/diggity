package spdx

import (
	"github.com/carbonetes/diggity/pkg/types"
	spdxcommon "github.com/spdx/tools-golang/spdx/v2/common"
	spdx23 "github.com/spdx/tools-golang/spdx/v2/v2_3"
)

func ToSPDX23Packages(components []types.Component) (packages []*spdx23.Package) {
	for _, c := range components {
		pkg := ToSPDX23Package(c)
		packages = append(packages, pkg)
	}
	return packages
}

func ToSPDX23Package(c types.Component) *spdx23.Package {
	return &spdx23.Package{
		PackageSPDXIdentifier:     spdxcommon.ElementID(Ref + c.PURL),
		PackageName:               c.Name,
		PackageDescription:        c.Description,
		PackageDownloadLocation:   DownloadLocation(c),
		PackageVersion:            c.Version,
		PackageLicenseConcluded:   LicensesDeclared(c),
		PackageExternalReferences: ExternalRefs(c),
		FilesAnalyzed:             true,
		PackageHomePage:           Homepage(c),
		PackageSourceInfo:         SourceInfo(c),
	}
}
