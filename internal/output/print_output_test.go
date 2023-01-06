package output

import (
	"path/filepath"
	"sort"
	"testing"

	"github.com/carbonetes/diggity/internal/model"
	"github.com/carbonetes/diggity/internal/parser"
)

var (
	printPackage1 = model.Package{
		ID:      "8fe93afb-86f2-4639-a3eb-6c4e787f210b",
		Name:    "lzo",
		Type:    "rpm",
		Version: "2.08",
	}
	printPackage2 = model.Package{
		ID:      "9583e9ec-df1d-484a-b560-8e1415ea92c2",
		Name:    "gitlab.com/yawning/obfs4.git",
		Type:    "go-module",
		Version: "v0.0.0-20220204003609-77af0cba934d",
	}
	printPackage3 = model.Package{
		ID:      "bdbd600f-dbdf-49a1-a329-a339f1123ffd",
		Name:    "scanelf",
		Type:    "apk",
		Version: "1.3.4-r0",
	}
	printPackage4 = model.Package{
		ID:      "418ee75b-cb1a-4abe-aad6-d757c7a91610",
		Name:    "scanf",
		Type:    "gem",
		Version: "1.0.0",
	}
	printDuplicate = model.Package{
		ID:      "418ee75b-cb1a-4abe-aad6-d757c7a91610",
		Name:    "scanf",
		Type:    "gem",
		Version: "1.0.0",
	}
	printNewVersion = model.Package{
		ID:      "519ee75c-cb1d-4abe-bad7-d758c7a91611",
		Name:    "scanf",
		Type:    "gem",
		Version: "2.0.0",
	}

	resultPackage = model.Package{
		ID:      "1891e769-cd46-494a-8b90-88a809a49104",
		Name:    "musl",
		Type:    "apk",
		Version: "apk 1.2.3-r0",
		Path:    filepath.Join("lib", "apk", "db", "installed"),
		Locations: []model.Location{
			{
				Path:      filepath.Join("lib", "apk", "db", "installed"),
				LayerHash: "9b7240956cfbfefddcd91a2195bfb2ed2cd17bdff81f21111849d643dfaf8131",
			},
		},
		Description: "the musl c library (libc) implementation",
		Licenses: []string{
			"MIT",
		},
		CPEs: []string{
			"cpe:2.3:a:musl:musl:1.2.3-r0:*:*:*:*:*:*:*",
		},
		PURL: model.PURL("pkg:alpine/musl@1.2.3-r0?arch=x86_64&upstream=musl&distro="),
		Metadata: parser.AlpineManifest{
			"Architecture":         "x86_64",
			"BuildTimestamp":       "1649396308",
			"GitCommitHashApk":     "ee13d43a53938d8a04ba787b9423f3270a3c14a7",
			"License":              "MIT",
			"Maintainer":           "Timo Teräs \u003ctimo.teras@iki.fi\u003e",
			"PackageDescription":   "the musl c library (libc) implementation",
			"PackageInstalledSize": "622592",
			"PackageName":          "musl",
			"PackageOrigin":        "cmusl",
			"PackageSize":          "383304",
			"PackageURL":           "https://musl.libc.org/",
			"PackageVersion":       "1.2.3-r0",
			"Provides":             "so:libc.musl-x86_64.so.1=1",
			"PullChecksum":         "Q1aCu0LmUDoAFSOX49uHvkYC1WasQ=",
		},
	}
)

func TestFinalizeResults(t *testing.T) {
	parser.Packages = []*model.Package{&printPackage1, &printPackage2, &printPackage3, &printPackage4, &printDuplicate}
	expected := []*model.Package{&printPackage1, &printPackage2, &printPackage3, &printPackage4}

	finalizeResults()

	if len(parser.Packages) != len(expected) {
		t.Errorf("Test Failed: Expected Packages of length %+v, Received: %+v.", len(expected), len(parser.Packages))
	}

	sortPackages()

	sort.Slice(expected, func(i, j int) bool {
		return expected[i].Name < expected[j].Name
	})

	for i, p := range parser.Packages {
		if p.Name != expected[i].Name {
			t.Errorf("Test Failed: Expected output of %v, received: %v", expected[i].Name, p.Name)
		}
	}
}

func TestSortResults(t *testing.T) {
	parser.Packages = []*model.Package{&printPackage1, &printPackage2, &printPackage3, &printPackage4, &printNewVersion}
	expected := []*model.Package{&printPackage2, &printPackage1, &printPackage3, &printPackage4, &printNewVersion}

	sortResults()

	for i, p := range parser.Packages {
		if p.Name != expected[i].Name {
			t.Errorf("Test Failed: Expected output of %v, received: %v", expected[i].Name, p.Name)
		}
	}
}
