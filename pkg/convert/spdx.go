package convert

// Reference: https://spdxutils.github.io/spdx-spec/v2.2.2/
// Composition Reference: https://spdxutils.github.io/spdx-spec/v2.2.2/composition-of-an-SPDX-document/

import (
	"fmt"
	"time"

	spdxutils "github.com/carbonetes/diggity/pkg/convert/spdx_utils"
	"github.com/carbonetes/diggity/pkg/model"
	spdxcommon "github.com/spdx/tools-golang/spdx/v2/common"
	spdx23 "github.com/spdx/tools-golang/spdx/v2/v2_3"
)

func ToSPDX(target string, pkgs *[]model.Package) *spdx23.Document {
	return &spdx23.Document{
		SPDXIdentifier: spdxutils.Ref + spdxutils.Doc,
		DocumentName:   target,
		SPDXVersion:    spdx23.Version,
		CreationInfo: &spdx23.CreationInfo{
			Created:            time.Now().UTC().String(),
			Creators:           spdxutils.CreateInfo,
			LicenseListVersion: spdxutils.LicenseListVersion,
		},
		DataLicense:       spdxutils.DataLicense,
		DocumentNamespace: spdxutils.FormatNamespace(target),
		Packages:          spdxJSONPackages(pkgs),
	}
}

func spdxJSONPackages(packages *[]model.Package) (spdxJSONPkgs []*spdx23.Package) {
	for _, p := range *packages {
		spdxPkg := &spdx23.Package{
			PackageSPDXIdentifier:     spdxcommon.ElementID(spdxutils.Ref + p.ID),
			PackageName:               p.Name,
			PackageDescription:        p.Description,
			PackageDownloadLocation:   spdxutils.DownloadLocation(&p),
			PackageLicenseConcluded:   spdxutils.LicensesDeclared(&p),
			PackageExternalReferences: spdxutils.ExternalRefs(&p),
			FilesAnalyzed:             false, // If false, indicates packages that represent metadata or URI references to a project, product, artifact, distribution or a component.
			PackageHomePage:           spdxutils.Homepage(&p),
			PackageLicenseDeclared:    spdxutils.LicensesDeclared(&p),
			PackageSourceInfo:         spdxutils.SourceInfo(&p),
			PackageVersion:            p.Version,
			PackageCopyrightText:      spdxutils.NoAssertion,
		}

		originatorType, originator := spdxutils.Originator(&p)
		if originatorType != "" && originator != "" {
			spdxPkg.PackageOriginator = &spdxcommon.Originator{
				Originator:     originator,
				OriginatorType: originatorType,
			}
		}

		spdxJSONPkgs = append(spdxJSONPkgs, spdxPkg)
	}
	return spdxJSONPkgs
}

func ToSPDXTagValue(args *model.Arguments, pkgs *[]model.Package) *[]string {
	spdxTagValues := new([]string)
	// Init Document Creation Information
	*spdxTagValues = append(*spdxTagValues, fmt.Sprintf(
		"SPDXVersion: %s\n"+
			"DataLicense: %s\n"+
			"SPDXID: %s\n"+
			"DocumentName: %s\n"+
			"DocumentNamespace: %s\n"+
			"LicenseListVersion: %s\n"+
			"Creator: %s\n"+
			"Creator: %s\n"+
			"Created: %+v",
		spdxutils.Version,                                     // SPDXVersion
		spdxutils.DataLicense,                                 // DataLicense
		spdxutils.Ref+spdxutils.Doc,                           // SPDXID
		spdxutils.FormatName(args),                            // DocumentName
		spdxutils.FormatNamespace(spdxutils.FormatName(args)), // DocumentNamespace
		spdxutils.LicenseListVersion,                          // LicenseListVersion
		spdxutils.Creator,                                     // Creator: Organization
		spdxutils.Tool,                                        // Creator: Tool
		time.Now().UTC().Format(time.RFC3339),                 // Created
	))

	// Parse Package Information to SPDX-TAG-VALUE Format
	for _, p := range *pkgs {
		*spdxTagValues = append(*spdxTagValues, fmt.Sprintf(
			"\n"+
				"##### Package: %s\n"+
				"\n"+
				"PackageName: %s\n"+
				"SPDXID: %s\n"+
				"PackageVersion: %s\n"+
				"PackageDownloadLocation: %s\n"+
				"FilesAnalyzed: %v\n"+
				"PackageLicenseConcluded: %s\n"+
				"PackageLicenseDeclared: %s\n"+
				"PackageCopyrightText: %s",
			p.Name,                         // Package
			p.Name,                         // PackageName
			spdxutils.FormatTagID(&p),      // SPDXID
			p.Version,                      // PackageVersion
			spdxutils.DownloadLocation(&p), // PackageDownloadLocation
			false,                          // FilesAnalyzed
			spdxutils.LicensesDeclared(&p), // PackageLicenseConcluded
			spdxutils.LicensesDeclared(&p), // PackageLicenseDeclared
			spdxutils.NoAssertion,          // PackageCopyrightText
		))

		for _, ref := range spdxutils.ExternalRefs(&p) {
			*spdxTagValues = append(*spdxTagValues, fmt.Sprintf(
				"ExternalRef: %s %s %s",
				ref.Category,
				ref.RefType,
				ref.Locator,
			))
		}
	}

	return spdxTagValues

}
