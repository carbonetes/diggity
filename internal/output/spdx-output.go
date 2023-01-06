// Package output - references
// Reference: https://spdx.github.io/spdx-spec/v2.2.2/
// Composition Reference: https://spdx.github.io/spdx-spec/v2.2.2/composition-of-an-SPDX-document/
package output

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/carbonetes/diggity/internal/model"
	"github.com/carbonetes/diggity/internal/model/output"
	spdx "github.com/carbonetes/diggity/internal/output/spdx-utils"
	"github.com/carbonetes/diggity/internal/parser"
)

// PrintSpdxJSON Print Packages in SPDX-JSON format
func printSpdxJSON() {
	spdxJSON, err := GetSpdxJSON(parser.Arguments.Image)

	if err != nil {
		panic(err)
	}

	if len(*parser.Arguments.OutputFile) > 0 {
		saveResultToFile(string(spdxJSON))
	} else {
		fmt.Printf("%+v", string(spdxJSON))
	}
}

// GetSpdxJSON Init SPDX-JSON Output
func GetSpdxJSON(image *string) ([]byte, error) {
	result := output.SpdxJSONDocument{
		SPDXID:      spdx.Ref + spdx.Doc,
		Name:        spdx.FormatName(image),
		SpdxVersion: spdx.Version,
		CreationInfo: output.CreationInfo{
			Created:            time.Now().UTC(),
			Creators:           spdx.CreateInfo,
			LicenseListVersion: spdx.LicenseListVersion,
		},
		DataLicense:       spdx.DataLicense,
		DocumentNamespace: spdx.FormatNamespace(spdx.FormatName(image)),
		SpdxJSONPackages:  spdxJSONPackages(parser.Packages),
	}
	return json.MarshalIndent(result, "", " ")
}

// spdxJSONPackages Get Packages in SPDX-JSON format
func spdxJSONPackages(packages []*model.Package) (spdxJSONPkgs []output.SpdxJSONPackage) {
	// Sort packages alphabetically
	sort.Slice(packages, func(i, j int) bool {
		return packages[i].Name < packages[j].Name
	})

	for _, p := range packages {
		spdxJSONPkgs = append(spdxJSONPkgs, output.SpdxJSONPackage{
			SpdxID:           spdx.Ref + p.ID,
			Name:             p.Name,
			Description:      p.Description,
			DownloadLocation: spdx.DownloadLocation(p),
			LicenseConcluded: spdx.LicensesDeclared(p),
			ExternalRefs:     spdx.ExternalRefs(p),
			FilesAnalyzed:    false, // If false, indicates packages that represent metadata or URI references to a project, product, artifact, distribution or a component.
			Homepage:         spdx.Homepage(p),
			LicenseDeclared:  spdx.LicensesDeclared(p),
			Originator:       spdx.Originator(p),
			SourceInfo:       spdx.SourceInfo(p),
			VersionInfo:      p.Version,
			Copyright:        spdx.NoAssertion,
		})
	}
	return spdxJSONPkgs
}

// PrintSpdxTagValue Print Packages in SPDX-TAG_VALUE format
func printSpdxTagValue() {
	spdxTagValues := GetSpdxTagValues()

	if len(*parser.Arguments.OutputFile) > 0 {
		saveResultToFile(stringSliceToString(spdxTagValues))
	} else {
		fmt.Printf("%+v", stringSliceToString(spdxTagValues))
	}
}

func stringSliceToString(slice []string) string {
	result := ""
	for _, s := range slice {
		result += fmt.Sprintln(s)
	}
	return result
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
		spdx.Version,                            // SPDXVersion
		spdx.DataLicense,                        // DataLicense
		spdx.Ref+spdx.Doc,                       // SPDXID
		spdx.FormatName(parser.Arguments.Image), // DocumentName
		spdx.FormatNamespace(spdx.FormatName(parser.Arguments.Image)), // DocumentNamespace
		spdx.LicenseListVersion,               // LicenseListVersion
		spdx.Creator,                          // Creator: Organization
		spdx.Tool,                             // Creator: Tool
		time.Now().UTC().Format(time.RFC3339), // Created
	))

	// Sort packages alphabetically
	sortPackages()

	// Parse Package Information to SPDX-TAG-VALUE Format
	for _, p := range parser.Packages {
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
			p.Name,                   // Package
			p.Name,                   // PackageName
			spdx.FormatTagID(p),      // SPDXID
			p.Version,                // PackageVersion
			spdx.DownloadLocation(p), // PackageDownloadLocation
			false,                    // FilesAnalyzed
			spdx.LicensesDeclared(p), // PackageLicenseConcluded
			spdx.LicensesDeclared(p), // PackageLicenseDeclared
			spdx.NoAssertion,         // PackageCopyrightText
		))

		for _, ref := range spdx.ExternalRefs(p) {
			spdxTagValues = append(spdxTagValues, fmt.Sprintf(
				"ExternalRef: %s %s %s",
				ref.ReferenceCategory,
				ref.ReferenceType,
				ref.ReferenceLocator,
			))
		}
	}

	return spdxTagValues
}
