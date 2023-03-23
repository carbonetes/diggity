package spdx

import (
	"path/filepath"
	"testing"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/model/metadata"
	"github.com/carbonetes/diggity/pkg/model/output"
	"github.com/carbonetes/diggity/pkg/parser/alpine"
	"github.com/carbonetes/diggity/pkg/parser/gem"
)

var (
	spdxPackage1 = model.Package{
		ID:      "8fe93afb-86f2-4639-a3eb-6c4e787f210b",
		Name:    "lzo",
		Type:    "rpm",
		Version: "2.08",
		Path:    filepath.Join("var", "lib", "rpm", "Packages"),
		Locations: []model.Location{
			{
				Path:      filepath.Join("var", "lib", "rpm", "Packages"),
				LayerHash: "69a15d957a7a6f77e3fe31f330da5f4b6b582f228917a713a7a9e59449a3f413",
			},
		},
		Description: "Data compression library with very fast (de)compression",
		Licenses: []string{
			"GPLv2+",
		},
		CPEs: []string{
			"cpe:2.3:a:centos:lzo:2.08-14.el8:*:*:*:*:*:*:*",
			"cpe:2.3:a:lzo:lzo:2.08-14.el8:*:*:*:*:*:*:*",
		},
		PURL: model.PURL("pkg:rpm/lzo@2.08?arch=x86_64"),
		Metadata: metadata.RPMMetadata{
			Release:      "14.el8",
			Architecture: "x86_64",
			SourceRpm:    "lzo-2.08-14.el8.src.rpm",
			License:      "GPLv2+",
			Size:         198757,
			Name:         "lzo",
			PGP:          "RSA/SHA256, Tue Jul  2 00:01:31 2019, Key ID 05b555b38483c65d",
			Summary:      "Data compression library with very fast (de)compression",
			Vendor:       "CentOS",
			Version:      "2.08",
		},
	}
	spdxPackage2 = model.Package{
		ID:      "9583e9ec-df1d-484a-b560-8e1415ea92c2",
		Name:    "gitlab.com/yawning/obfs4.git",
		Type:    "go-module",
		Version: "v0.0.0-20220204003609-77af0cba934d",
		Path:    "",
		Locations: []model.Location{
			{
				Path:      filepath.Join("bin", "gost"),
				LayerHash: "0cd4836a36e094e1870a2e6c2578a7ad9d9cb42a7313944a6d05ab72892fc3c3",
			},
		},
		Description: "",
		Licenses:    []string{},
		CPEs: []string{
			"cpe:2.3:a:yawning:obfs4.git:v0.0.0-20220204003609-77af0cba934d:*:*:*:*:*:*:*",
			"cpe:2.3:a:obfs4.git:obfs4.git:v0.0.0-20220204003609-77af0cba934d:*:*:*:*:*:*:*",
		},
		PURL: model.PURL("pkg:go/gitlab.com/yawning/obfs4.git@v0.0.0-20220204003609-77af0cba934d"),
		Metadata: metadata.GoBinMetadata{
			Architecture:     "amd64",
			Compiler:         "gc",
			OS:               "linux",
			GoCompileRelease: "go1.19.1",
			H1Digest:         "h1:tJ8F7ABaQ3p3wjxwXiWSktVDgjZEXkvaRawd2rIq5ws=",
			Path:             "gitlab.com/yawning/obfs4.git",
			Version:          "v0.0.0-20220204003609-77af0cba934d",
		},
	}
	spdxPackage3 = model.Package{
		ID:      "bdbd600f-dbdf-49a1-a329-a339f1123ffd",
		Name:    "scanelf",
		Type:    "apk",
		Version: "1.3.4-r0",
		Path:    filepath.Join("lib", "apk", "db", "installed"),
		Locations: []model.Location{
			{
				Path:      filepath.Join("lib", "apk", "db", "installed"),
				LayerHash: "1288696addccc4013c5bcf61c1b6c38128a7214a0942976792918b51912d90f7",
			},
			{
				Path:      filepath.Join("lib", "apk", "db", "installed"),
				LayerHash: "1288696addccc4013c5bcf61c1b6c38128a7214a0942976792918b51912d90f7",
			},
		},
		Description: "Scan ELF binaries for stuff",
		Licenses: []string{
			"GPL-2.0-only",
		},
		CPEs: []string{
			"cpe:2.3:a:scanelf:scanelf:1.3.4-r0:*:*:*:*:*:*:*",
		},
		PURL: model.PURL("pkg:apk/alpine/scanelf@1.3.4-r0?arch=x86_64\u0026upstream=pax-utils\u0026distro=alpine"),
		Metadata: alpine.Manifest{
			"Architecture":         "x86_64",
			"BuildTimestamp":       "1651005390",
			"GitCommitHashApk":     "d7ae612a3cc5f827289d915783b4cbf8c7207947",
			"License":              "GPL-2.0-only",
			"Maintainer":           "Natanael Copa \u003cncopa@alpinelinux.org\u003e",
			"PackageDescription":   "Scan ELF binaries for stuff",
			"PackageInstalledSize": "94208",
			"PackageName":          "scanelf",
			"PackageOrigin":        "pax-utils",
			"PackageSize":          "36745",
			"PackageURL":           "https://wiki.gentoo.org/wiki/Hardened/PaX_Utilities",
			"PackageVersion":       "1.3.4-r0",
			"Provides":             "cmd:scanelf=1.3.4-r0",
			"PullChecksum":         "Q1Gcqe+ND8DFOlhM3R0o5KyZjR2oE=",
			"PullDependencies":     "pax-utils",
		},
	}
	spdxPackage4 = model.Package{
		ID:      "418ee75b-cb1a-4abe-aad6-d757c7a91610",
		Name:    "scanf",
		Type:    "gem",
		Version: "1.0.0",
		Path:    "",
		Locations: []model.Location{
			{
				Path:      filepath.Join("usr", "share", "gems", "specifications", "default", "scanf-1.0.0.gemspec"),
				LayerHash: "a67d9e51873dfbda0e6af0f9971ccea211405916ede446f52b5e7f3ea9d71fc3",
			},
		},
		Description: "scanf is an implementation of the C function scanf(3).",
		Licenses: []string{
			"BSD2Clause",
		},
		CPEs: []string{
			"cpe:2.3:a:scanf:scanf:1.0.0:*:*:*:*:*:*:*",
		},
		PURL: model.PURL("pkg:gem/scanf@1.0.0"),
		Metadata: gem.Metadata{
			"authors":     []string{"David Alan Black"},
			"bindir":      "exe",
			"date":        "2017-12-11",
			"description": "scanf is an implementation of the C function scanf(3).",
			"email":       "[dblack@superlink.net]",
			"files":       []string{"scanfrb"},
			"homepage":    "https://github.com/ruby/scanf",
			"licenses": []string{
				"BSD2Clause",
			},
			"name":                      "scanf",
			"require_paths":             "[lib]",
			"required_ruby_version":     "Gem::Requirement.new(\u003e= 2.3.0)",
			"required_rubygems_version": "Gem::Requirement.new(\u003e= 0) if s.respond_to? :required_rubygems_version=",
			"rubygems_version":          "2.7.6.2",
			"specification_version":     "4",
			"summary":                   "scanf is an implementation of the C function scanf(3).  if s.respond_to? :specification_version then",
			"version":                   "1.0.0",
		},
	}
)

func TestGetSPDXJSON(t *testing.T) {
	tests := []string{"bom:latest", "smartentry/centos:latest", "buluma/centos:6", "furynix/fedora:29", "test/image:test-tag", "test_image/image:test_tag", "image"}
	for _, test := range tests {
		_, err := GetSpdxJSON(&test)
		if err != nil {
			t.Error("Test Failed: Error occurred while parsing test spdx json.")
		}
	}
}
func TestSpdxJSONPackages(t *testing.T) {
	packages := []*model.Package{&spdxPackage1, &spdxPackage2, &spdxPackage3, &spdxPackage4}
	expected := []output.SpdxJSONPackage{
		{
			SpdxID:           "SPDXRef-9583e9ec-df1d-484a-b560-8e1415ea92c2",
			Name:             "gitlab.com/yawning/obfs4.git",
			Description:      "",
			DownloadLocation: "NOASSERTION",
			LicenseConcluded: "NONE",
			ExternalRefs: []output.ExternalRef{
				{
					ReferenceCategory: "SECURITY",
					ReferenceLocator:  "cpe:2.3:a:yawning:obfs4.git:v0.0.0-20220204003609-77af0cba934d:*:*:*:*:*:*:*",
					ReferenceType:     "cpe23Type",
				},
				{
					ReferenceCategory: "SECURITY",
					ReferenceLocator:  "cpe:2.3:a:obfs4.git:obfs4.git:v0.0.0-20220204003609-77af0cba934d:*:*:*:*:*:*:*",
					ReferenceType:     "cpe23Type",
				},
				{
					ReferenceCategory: "PACKAGE_MANAGER",
					ReferenceLocator:  "pkg:go/gitlab.com/yawning/obfs4.git@v0.0.0-20220204003609-77af0cba934d",
					ReferenceType:     "purl",
				},
			},
			FilesAnalyzed:   false,
			Homepage:        "",
			LicenseDeclared: "NONE",
			Originator:      "",
			SourceInfo:      "Information parsed from go-module information: bin/gost",
			VersionInfo:     "v0.0.0-20220204003609-77af0cba934d",
			Copyright:       "NOASSERTION",
		},
		{
			SpdxID:           "SPDXRef-8fe93afb-86f2-4639-a3eb-6c4e787f210b",
			Name:             "lzo",
			Description:      "Data compression library with very fast (de)compression",
			DownloadLocation: "NOASSERTION",
			LicenseConcluded: "NOASSERTION",
			ExternalRefs: []output.ExternalRef{
				{
					ReferenceCategory: "SECURITY",
					ReferenceLocator:  "cpe:2.3:a:centos:lzo:2.08-14.el8:*:*:*:*:*:*:*",
					ReferenceType:     "cpe23Type",
				},
				{
					ReferenceCategory: "SECURITY",
					ReferenceLocator:  "cpe:2.3:a:lzo:lzo:2.08-14.el8:*:*:*:*:*:*:*",
					ReferenceType:     "cpe23Type",
				},
				{
					ReferenceCategory: "PACKAGE_MANAGER",
					ReferenceLocator:  "pkg:rpm/lzo@2.08?arch=x86_64",
					ReferenceType:     "purl",
				},
			},
			FilesAnalyzed:   false,
			Homepage:        "",
			LicenseDeclared: "NOASSERTION",
			Originator:      "Organization: CentOS",
			SourceInfo:      "Information parsed from RPM DB: var/lib/rpm/Packages",
			VersionInfo:     "2.08",
			Copyright:       "NOASSERTION",
		},
		{
			SpdxID:           "SPDXRef-bdbd600f-dbdf-49a1-a329-a339f1123ffd",
			Name:             "scanelf",
			Description:      "Scan ELF binaries for stuff",
			DownloadLocation: "https://wiki.gentoo.org/wiki/Hardened/PaX_Utilities",
			LicenseConcluded: "GPL-2.0-only",
			ExternalRefs: []output.ExternalRef{
				{
					ReferenceCategory: "SECURITY",
					ReferenceLocator:  "cpe:2.3:a:scanelf:scanelf:1.3.4-r0:*:*:*:*:*:*:*",
					ReferenceType:     "cpe23Type",
				},
				{
					ReferenceCategory: "PACKAGE_MANAGER",
					ReferenceLocator:  "pkg:apk/alpine/scanelf@1.3.4-r0?arch=x86_64\u0026upstream=pax-utils\u0026distro=alpine",
					ReferenceType:     "purl",
				},
			},
			FilesAnalyzed:   false,
			Homepage:        "",
			LicenseDeclared: "GPL-2.0-only",
			Originator:      "Person: Natanael Copa \u003cncopa@alpinelinux.org\u003e",
			SourceInfo:      "Information parsed from APK DB: lib/apk/db/installed, lib/apk/db/installed",
			VersionInfo:     "1.3.4-r0",
			Copyright:       "NOASSERTION",
		},
		{
			SpdxID:           "SPDXRef-418ee75b-cb1a-4abe-aad6-d757c7a91610",
			Name:             "scanf",
			Description:      "scanf is an implementation of the C function scanf(3).",
			DownloadLocation: "NOASSERTION",
			LicenseConcluded: "NOASSERTION",
			ExternalRefs: []output.ExternalRef{
				{
					ReferenceCategory: "SECURITY",
					ReferenceLocator:  "cpe:2.3:a:scanf:scanf:1.0.0:*:*:*:*:*:*:*",
					ReferenceType:     "cpe23Type",
				},
				{
					ReferenceCategory: "PACKAGE_MANAGER",
					ReferenceLocator:  "pkg:gem/scanf@1.0.0",
					ReferenceType:     "purl",
				},
			},
			FilesAnalyzed:   false,
			Homepage:        "https://github.com/ruby/scanf",
			LicenseDeclared: "NOASSERTION",
			Originator:      "Person: David Alan Black",
			SourceInfo:      "Information parsed from gem metadata: usr/share/gems/specifications/default/scanf-1.0.0.gemspec",
			VersionInfo:     "1.0.0",
			Copyright:       "NOASSERTION",
		},
	}

	_output := spdxJSONPackages(packages)

	for i, spdxPkg := range _output {
		for j, exRef := range _output[i].ExternalRefs {
			if exRef != expected[i].ExternalRefs[j] {
				t.Errorf("Test Failed: Expected output of %v, received: %v", expected[i].ExternalRefs[j], exRef)
			}
		}

		if spdxPkg.SpdxID != expected[i].SpdxID ||
			spdxPkg.Name != expected[i].Name ||
			spdxPkg.Description != expected[i].Description ||
			spdxPkg.DownloadLocation != expected[i].DownloadLocation ||
			spdxPkg.LicenseConcluded != expected[i].LicenseConcluded ||
			spdxPkg.FilesAnalyzed != expected[i].FilesAnalyzed ||
			spdxPkg.Homepage != expected[i].Homepage ||
			spdxPkg.LicenseDeclared != expected[i].LicenseDeclared ||
			spdxPkg.Originator != expected[i].Originator ||
			spdxPkg.SourceInfo != expected[i].SourceInfo ||
			spdxPkg.VersionInfo != expected[i].VersionInfo ||
			spdxPkg.Copyright != expected[i].Copyright {
			t.Errorf("Test Failed: Expected output of %v, received: %v", expected[i], _output[i])
		}

	}

}
