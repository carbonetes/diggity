package portage

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/carbonetes/diggity/internal/model"
	"github.com/carbonetes/diggity/internal/model/metadata"
)

type (
	PortageNameVersionResult struct {
		input   string
		name    string
		version string
	}

	ParsePortageFileResult struct {
		input    string
		expected metadata.PortageFile
	}

	PortagePurlResult struct {
		_package *model.Package
		expected model.PURL
	}
)

var (
	portagePackage1 = model.Package{
		Name:    filepath.Join("dev-util", "gperf"),
		Type:    portage,
		Version: "3.1-r1",
		Path:    filepath.Join("dev-util", "gperf-3.1-r1"),
		Locations: []model.Location{
			{
				Path: filepath.Join("var", "db", "pkg", "dev-util", "gperf-3.1-r1", portageContent),
			},
		},
		Licenses: []string{
			"GPL-2",
		},
		CPEs: []string{
			"cpe:2.3:a:dev-util\\gperf:dev-util\\gperf:3.1-r1:*:*:*:*:*:*:*",
		},
		PURL: model.PURL("pkg:ebuild/dev-util/gperf@3.1-r1"),
		Metadata: metadata.PortageMetadata{
			Size: 516979,
		},
	}
	portagePackage2 = model.Package{
		Name:    filepath.Join("acct-group", "audio"),
		Type:    portage,
		Version: "0-r1",
		Path:    filepath.Join("acct-group", "audio-0-r1"),
		Locations: []model.Location{
			{
				Path: filepath.Join("var", "db", "pkg", "acct-group", "audio-0-r1", portageContent),
			},
		},
		CPEs: []string{
			"cpe:2.3:a:acct-group\\audio:acct-group\\audio:0-r1:*:*:*:*:*:*:*",
		},
		PURL: model.PURL("pkg:ebuild/acct-group/audio@0-r1"),
		Metadata: metadata.PortageMetadata{
			Size: 11,
			Files: []metadata.PortageFile{
				{
					Path: "/usr/lib/sysusers.d/acct-group-audio.conf",
					Digest: metadata.PortageDigest{
						Algorithm: "md5",
						Value:     "5e63f5b43622b84d87989175ae09a94b",
					},
				},
			},
		},
	}
)

func TestPortageNameVersion(t *testing.T) {
	tests := []PortageNameVersionResult{
		{filepath.Join("var", "db", "pkg", "acct-group", "cdrom-0-r1"), filepath.Join("acct-group", "cdrom"), "0-r1"},
		{filepath.Join("var", "db", "pkg", "dev-libs", "mpc-1.3.1"), filepath.Join("dev-libs", "mpc"), "1.3.1"},
		{filepath.Join("var", "db", "pkg", "net-misc", "curl-7.87.0-r2"), filepath.Join("net-misc", "curl"), "7.87.0-r2"},
		{filepath.Join("var", "db", "pkg", "sys-apps", "attr-2.5.1-r2"), filepath.Join("sys-apps", "attr"), "2.5.1-r2"},
		{filepath.Join("var", "db", "pkg", "virtual", "libc-1-r1"), filepath.Join("virtual", "libc"), "1-r1"},
	}

	for _, test := range tests {
		if outputName, outputVersion := portageNameVersion(test.input); outputName != test.name || outputVersion != test.version {
			t.Errorf("Test Failed: Expected output of [%v, %v], received: [%v, %v]", test.name, test.version, outputName, outputVersion)
		}
	}
}

func TestGetPortageLicenses(t *testing.T) {
	path := filepath.Join("..", "..", "..", "docs", "references", "portage", "var", "db", "pkg", "dev-util", "gperf-3.1-r1", portageContent)
	var _package model.Package
	expected := []string{"GPL-2"}

	if err := getPortageLicenses(&_package, path); err != nil {
		t.Error("Test Failed: Error occurred while reading portage LICENSE file.")
	}

	if len(_package.Licenses) != len(expected) {
		t.Errorf("Test Failed:\n Expected output of %v \n, Received: %v \n", expected, _package.Licenses)
	}

	for i, license := range _package.Licenses {
		if license != expected[i] {
			t.Errorf("Test Failed:\n Expected output of %v \n, Received: %v \n", expected[i], license)
		}
	}

}

func TestGetPortageFiles(t *testing.T) {
	path := filepath.Join("..", "..", "..", "docs", "references", "portage", "var", "db", "pkg", "dev-util", "gperf-3.1-r1", portageContent)
	var md metadata.PortageMetadata

	if err := getPortageFiles(&md, path); err != nil {
		t.Error("Test Failed: Error occurred while reading portage CONTENT file.")
	}
	if len(md.Files) == 0 {
		t.Error("Test Failed: Expected non-empty files.")
	}

}

func TestParsePortageFile(t *testing.T) {
	tests := []ParsePortageFileResult{
		{"obj /usr/lib/sysusers.d/acct-group-audio.conf 5e63f5b43622b84d87989175ae09a94b 1625882852",
			metadata.PortageFile{
				Path: "/usr/lib/sysusers.d/acct-group-audio.conf",
				Digest: metadata.PortageDigest{
					Algorithm: portageAlgorithm,
					Value:     "5e63f5b43622b84d87989175ae09a94b",
				},
			}},
		{"obj /usr/share/doc/gperf-3.1-r1/AUTHORS ddcc95b0e8d8baf1ae14fff3a47c5487 1671388867",
			metadata.PortageFile{
				Path: "/usr/share/doc/gperf-3.1-r1/AUTHORS",
				Digest: metadata.PortageDigest{
					Algorithm: portageAlgorithm,
					Value:     "ddcc95b0e8d8baf1ae14fff3a47c5487",
				},
			}},
	}

	for _, test := range tests {
		if output := parsePortageFile(test.input); !reflect.DeepEqual(output, test.expected) {
			t.Errorf("Test Failed:\n Expected output of %v \n, Received: %v \n", test.expected, output)
		}
	}
}

func TestParsePortagePackageURL(t *testing.T) {
	_package1 := model.Package{
		Name:    portagePackage1.Name,
		Version: portagePackage1.Version,
	}
	_package2 := model.Package{
		Name:    portagePackage2.Name,
		Version: portagePackage2.Version,
	}

	tests := []PortagePurlResult{
		{&_package1, model.PURL(portagePackage1.PURL)},
		{&_package2, model.PURL(portagePackage2.PURL)},
	}
	for _, test := range tests {
		parsePortagePURL(test._package)
		if test._package.PURL != test.expected {
			t.Errorf("Test Failed: Expected an output of %v, received: %v", test.expected, test._package.PURL)
		}
	}
}
