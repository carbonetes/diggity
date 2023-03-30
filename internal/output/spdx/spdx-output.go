package spdx

// Package output - references
// Reference: https://spdxutils.github.io/spdx-spec/v2.2.2/
// Composition Reference: https://spdxutils.github.io/spdx-spec/v2.2.2/composition-of-an-SPDX-document/

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/carbonetes/diggity/internal/output/save"
	spdxutils "github.com/carbonetes/diggity/internal/output/spdx/spdx-utils"
	"github.com/carbonetes/diggity/internal/output/util"
	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/parser/bom"
	spdxcommon "github.com/spdx/tools-golang/spdx/common"
	spdx22 "github.com/spdx/tools-golang/spdx/v2_2"
)

// PrintSpdxJSON Print Packages in SPDX-JSON format
func PrintSpdxJSON() {
	spdxJSON, err := GetSpdxJSON(bom.Arguments.Image)

	if err != nil {
		panic(err)
	}

	if len(*bom.Arguments.OutputFile) > 0 {
		save.ResultToFile(string(spdxJSON))
	} else {
		fmt.Printf("%+v\n", string(spdxJSON))
	}
}

// GetSpdxJSON Init SPDX-JSON Output
func GetSpdxJSON(image *string) ([]byte, error) {
	result := spdx22.Document{
		SPDXIdentifier: spdxutils.Ref + spdxutils.Doc,
		DocumentName:   spdxutils.FormatName(image),
		SPDXVersion:    spdxutils.Version,
		CreationInfo: &spdx22.CreationInfo{
			Created:            time.Now().UTC().String(),
			Creators:           spdxutils.CreateInfo,
			LicenseListVersion: spdxutils.LicenseListVersion,
		},
		DataLicense:       spdxutils.DataLicense,
		DocumentNamespace: spdxutils.FormatNamespace(spdxutils.FormatName(image)),
		Packages:          spdxJSONPackages(bom.Packages),
	}
	return json.MarshalIndent(result, "", " ")
}

// spdxJSONPackages Get Packages in SPDX-JSON format
func spdxJSONPackages(packages []*model.Package) (spdxJSONPkgs []*spdx22.Package) {
	// Sort packages alphabetically
	sort.Slice(packages, func(i, j int) bool {
		return packages[i].Name < packages[j].Name
	})

	for _, p := range packages {
		spdxPkg := &spdx22.Package{
			PackageSPDXIdentifier:     spdxcommon.ElementID(spdxutils.Ref + p.ID),
			PackageName:               p.Name,
			PackageDescription:        p.Description,
			PackageDownloadLocation:   spdxutils.DownloadLocation(p),
			PackageLicenseConcluded:   spdxutils.LicensesDeclared(p),
			PackageExternalReferences: spdxutils.ExternalRefs(p),
			FilesAnalyzed:             false, // If false, indicates packages that represent metadata or URI references to a project, product, artifact, distribution or a component.
			PackageHomePage:           spdxutils.Homepage(p),
			PackageLicenseDeclared:    spdxutils.LicensesDeclared(p),
			PackageSourceInfo:         spdxutils.SourceInfo(p),
			PackageVersion:            p.Version,
			PackageCopyrightText:      spdxutils.NoAssertion,
		}

		originatorType, originator := spdxutils.Originator(p)
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

// PrintSpdxTagValue Print Packages in SPDX-TAG_VALUE format
func PrintSpdxTagValue() {
	spdxTagValues := GetSpdxTagValues()

	if len(*bom.Arguments.OutputFile) > 0 {
		save.ResultToFile(stringSliceToString(spdxTagValues))
	} else {
		fmt.Printf("%+v", stringSliceToString(spdxTagValues))
	}
}

// GetSpdxTagValues Parse SPDX-TAG_VALUE format
func GetSpdxTagValues() (spdxTagValues []string) {
	// Init Document Creation Information
	spdxTagValues = append(spdxTagValues, fmt.Sprintf(
		"SPDXVersion: %s\n"+
			"DataLicense: %s\n"+
			"SPDXID: %s\n"+
			"DocumentName: %s\n"+
			"DocumentNamespace: %s\n"+
			"LicenseListVersion: %s\n"+
			"Creator: %s\n"+
			"Creator: %s\n"+
			"Created: %+v",
		spdxutils.Version,                         // SPDXVersion
		spdxutils.DataLicense,                     // DataLicense
		spdxutils.Ref+spdxutils.Doc,               // SPDXID
		spdxutils.FormatName(bom.Arguments.Image), // DocumentName
		spdxutils.FormatNamespace(spdxutils.FormatName(bom.Arguments.Image)), // DocumentNamespace
		spdxutils.LicenseListVersion,                                         // LicenseListVersion
		spdxutils.Creator,                                                    // Creator: Organization
		spdxutils.Tool,                                                       // Creator: Tool
		time.Now().UTC().Format(time.RFC3339),                                // Created
	))

	// Sort packages alphabetically
	util.SortPackages()

	// Parse Package Information to SPDX-TAG-VALUE Format
	for _, p := range bom.Packages {
		spdxTagValues = append(spdxTagValues, fmt.Sprintf(
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
			p.Name,                        // Package
			p.Name,                        // PackageName
			spdxutils.FormatTagID(p),      // SPDXID
			p.Version,                     // PackageVersion
			spdxutils.DownloadLocation(p), // PackageDownloadLocation
			false,                         // FilesAnalyzed
			spdxutils.LicensesDeclared(p), // PackageLicenseConcluded
			spdxutils.LicensesDeclared(p), // PackageLicenseDeclared
			spdxutils.NoAssertion,         // PackageCopyrightText
		))

		for _, ref := range spdxutils.ExternalRefs(p) {
			spdxTagValues = append(spdxTagValues, fmt.Sprintf(
				"ExternalRef: %s %s %s",
				ref.Category,
				ref.RefType,
				ref.Locator,
			))
		}
	}

	return spdxTagValues
}

// convert spdx-tag-values to single string
func stringSliceToString(slice []string) string {
	result := ""
	for _, s := range slice {
		result += fmt.Sprintln(s)
	}
	return result
}
