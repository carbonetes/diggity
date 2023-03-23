package conan

import (
	"path/filepath"
	"testing"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/model/metadata"
)

type (
	InitConanPackageResult struct {
		location *model.Location
		metadata interface{}
		expected *model.Package
	}

	ConanPurlResult struct {
		_package *model.Package
		expected model.PURL
	}

	ConanNameVersionResult struct {
		input   string
		name    string
		version string
	}
)

var (
	conanPackage1 = model.Package{
		Name:    "zlib",
		Type:    conan,
		Version: "1.2.11",
		Path:    "zlib",
		Locations: []model.Location{
			{
				Path: conanFile,
			},
		},
		Description: "",
		Licenses:    []string{},
		CPEs: []string{
			"cpe:2.3:a:zlib:zlib:1.2.11:*:*:*:*:*:*:*",
		},
		PURL: model.PURL("pkg:conan/zlib@1.2.11"),
		Metadata: metadata.ConanMetadata{
			Name:    "zlib",
			Version: "1.2.11",
		},
	}
	conanPackage2 = model.Package{
		Name:    "pkgb",
		Type:    conan,
		Version: "0.1",
		Path:    "pkgb",
		Locations: []model.Location{
			{
				Path: conanLock,
			},
		},
		Description: "",
		Licenses:    []string{},
		CPEs: []string{
			"cpe:2.3:a:pkgb:pkgb:0.1:*:*:*:*:*:*:*",
		},
		PURL: model.PURL("pkg:conan/pkgb@0.1"),
		Metadata: metadata.ConanLockNode{
			Ref:      "pkgb/0.1@user/testing",
			Path:     "..\\conanfile.py",
			Context:  "host",
			Requires: []string{"1"},
		},
	}

	conanMetadata1 = "zlib/1.2.11"
	conanMetadata2 = metadata.ConanLockNode{
		Ref:      "pkgb/0.1@user/testing",
		Requires: []string{"1"},
		Path:     "..\\conanfile.py",
		Context:  "host",
	}
)

func TestReadConanFileContent(t *testing.T) {
	conanPath := filepath.Join("..", "..", "..", "docs", "references", "conan", conanFile)
	testLocation := model.Location{Path: conanPath}
	err := readConanFileContent(&testLocation)
	if err != nil {
		t.Error("Test Failed: Error occurred while reading conanfile.txt content.")
	}
}

func TestReadConanLockContent(t *testing.T) {
	conanLockPath := filepath.Join("..", "..", "..", "docs", "references", "conan", conanLock)
	testLocation := model.Location{Path: conanLockPath}
	err := readConanFileContent(&testLocation)
	if err != nil {
		t.Error("Test Failed: Error occurred while reading conan.lock content.")
	}
}

func TestInitConanPackage(t *testing.T) {
	tests := []InitConanPackageResult{
		{&model.Location{Path: conanFile}, conanMetadata1, &conanPackage1},
		{&model.Location{Path: conanLock}, conanMetadata2, &conanPackage2},
	}

	for _, test := range tests {
		output := initConanPackage(test.location, test.metadata)

		if output.Type != test.expected.Type ||
			output.Path != test.expected.Path ||
			output.Name != test.expected.Name ||
			output.Version != test.expected.Version ||
			output.Description != test.expected.Description ||
			len(output.Licenses) != len(test.expected.Licenses) ||
			len(output.Locations) != len(test.expected.Locations) ||
			len(output.CPEs) != len(test.expected.CPEs) ||
			string(output.PURL) != string(test.expected.PURL) {
			t.Errorf("Test Failed:\n Expected output of %v \n, Received: %v \n", test.expected, output)
		}

		for i := range output.Licenses {
			if output.Licenses[i] != test.expected.Licenses[i] {
				t.Errorf("Test Failed:\n Expected output of %v \n, Received: %v \n", test.expected.Licenses[i], output.Licenses[i])
			}
		}
		for i := range output.Locations {
			if output.Locations[i] != test.expected.Locations[i] {
				t.Errorf("Test Failed:\n Expected output of %v \n, Received: %v \n", test.expected.Locations[i], output.Locations[i])
			}
		}
		for i := range output.CPEs {
			if output.CPEs[i] != test.expected.CPEs[i] {
				t.Errorf("Test Failed:\n Expected output of %v \n, Received: %v \n", test.expected.CPEs[i], output.CPEs[i])
			}
		}

		switch test.metadata.(type) {
		case string:
			outputMetadata := output.Metadata.(metadata.ConanMetadata)
			expectedMetadata := test.expected.Metadata.(metadata.ConanMetadata)

			if outputMetadata.Name != expectedMetadata.Name ||
				outputMetadata.Version != expectedMetadata.Version {
				t.Errorf("Test Failed:\n Expected output of %v \n, Received: %v \n", test.expected, output)
			}
		case metadata.ConanLockNode:
			outputMetadata := output.Metadata.(metadata.ConanLockNode)
			expectedMetadata := test.expected.Metadata.(metadata.ConanLockNode)

			if outputMetadata.Ref != expectedMetadata.Ref ||
				outputMetadata.Path != expectedMetadata.Path ||
				outputMetadata.Context != expectedMetadata.Context ||
				len(outputMetadata.Requires) != len(expectedMetadata.Requires) ||
				outputMetadata.PackageID != expectedMetadata.PackageID ||
				outputMetadata.Prev != expectedMetadata.Prev ||
				outputMetadata.BuildRequires != expectedMetadata.PythonRequires {
				t.Errorf("Test Failed:\n Expected output of %v \n, Received: %v \n", test.expected, output)
			}
		}
	}
}

func TestParseConanPackageURL(t *testing.T) {
	_package1 := model.Package{
		Name:    conanPackage1.Name,
		Version: conanPackage1.Version,
	}
	_package2 := model.Package{
		Name:    conanPackage2.Name,
		Version: conanPackage2.Version,
	}

	tests := []ConanPurlResult{
		{&_package1, model.PURL(conanPackage1.PURL)},
		{&_package2, model.PURL(conanPackage2.PURL)},
	}
	for _, test := range tests {
		parseConanPackageURL(test._package)
		if test._package.PURL != test.expected {
			t.Errorf("Test Failed: Expected an output of %v, received: %v", test.expected, test._package.PURL)
		}
	}
}

func TestConanNameVersion(t *testing.T) {
	tests := []ConanNameVersionResult{
		{"name/version", "name", "version"},
		{"name/version@user/testing", "name", "version"},
		{"poco/1.9.4", "poco", "1.9.4"},
		{"zlib/1.2.11", "zlib", "1.2.11"},
		{"pkga/0.1@user/testing", "pkga", "0.1"},
	}

	for _, test := range tests {
		if outputName, outputVersion := conanNameVersion(test.input); outputName != test.name || outputVersion != test.version {
			t.Errorf("Test Failed: Expected output of [%v, %v], received: [%v, %v]", test.name, test.version, outputName, outputVersion)
		}
	}
}
