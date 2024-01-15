package spdx

import (
	"time"

	"github.com/carbonetes/diggity/pkg/types"
	spdxcommon "github.com/spdx/tools-golang/spdx/v2/common"
	spdx23 "github.com/spdx/tools-golang/spdx/v2/v2_3"
)

// Convert diggity sbom format to spdx 2.3 document
func ToSPDX23(sbom types.SBOM, input string) *spdx23.Document {
	return &spdx23.Document{
		SPDXIdentifier:    Ref + Doc,
		SPDXVersion:       spdx23.Version,
		DocumentName:      input,
		DocumentNamespace: FormatNamespace(input),
		DataLicense:       DataLicense,

		CreationInfo: &spdx23.CreationInfo{
			Created:  time.Now().UTC().String(),
			Creators: Creators,
		},
		Packages: ToSPDX23Packages(sbom.Components),
	}
}

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
