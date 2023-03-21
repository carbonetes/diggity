package github

import (
	"path/filepath"
	"testing"

	"github.com/carbonetes/diggity/internal/parser/bom"
	"github.com/carbonetes/diggity/pkg/model"
)

var (
	depPackage1 = model.Package{
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
	}
	depPackage2 = model.Package{
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
	}
)

type (
	GetSnapshotMetadataResult struct {
		input    *model.Distro
		expected map[string]interface{}
	}

	GetPackageManifestsResult struct {
		_package *model.Package
		image    string
		expected string
	}

	PurlNameResult struct {
		input    model.PURL
		expected string
	}
)

func TestGetSnapshotMetadata(t *testing.T) {
	distro1 := model.Distro{
		ID:        "alpine",
		VersionID: "3.17.1",
	}
	distro2 := model.Distro{}

	tests := []GetSnapshotMetadataResult{
		{&distro1, map[string]interface{}{"diggity:distro": "pkg:generic/alpine@3.17.1"}},
		{&distro2, nil},
	}

	for _, test := range tests {
		if output := getSnapshotMetadata(test.input); output["diggity:distro"] != test.expected["diggity:distro"] {
			t.Errorf("Test Failed: Expected output of %v, received: %v", test.expected, output)
		}
	}
}

func TestGetPackageManifests(t *testing.T) {
	tests := []GetPackageManifestsResult{
		{&depPackage1, "gost/go", "gost/go:/bin/gost"},
		{&depPackage2, "alpine", "alpine:/lib/apk/db/installed"},
	}

	for _, test := range tests {
		bom.Packages = []*model.Package{test._package}
		output := getPackageManifests(test.image)
		if _, exists := output[test.expected]; !exists {
			t.Errorf("Test Failed: Expected output of %v, received none", output)
		}

		outputManifest := output[test.expected]
		if outputManifest.Name != test.expected {
			t.Errorf("Test Failed: Expected output of %v, received: %v", test.expected, outputManifest.Name)
		}
		if outputManifest.File.SourceLocation != test.expected {
			t.Errorf("Test Failed: Expected output of %v, received: %v", test.expected, outputManifest.File.SourceLocation)
		}
	}
}

func TestPurlName(t *testing.T) {
	tests := []PurlNameResult{
		{depPackage1.PURL, "pkg:go/gitlab.com/yawning/obfs4.git@v0.0.0-20220204003609-77af0cba934d"},
		{depPackage2.PURL, "pkg:apk/alpine/scanelf@1.3.4-r0"},
	}

	for _, test := range tests {
		if output := purlName(test.input); output != test.expected {
			t.Errorf("Test Failed: Expected output of %v, received: %v", test.expected, output)
		}
	}
}
