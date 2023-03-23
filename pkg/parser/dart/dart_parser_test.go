package dart

import (
	"path/filepath"
	"testing"

	"github.com/carbonetes/diggity/pkg/model"
)

type (
	DartPurlResult struct {
		_package *model.Package
		expected model.PURL
	}
)

var (
	dartPackage1 = model.Package{
		Name:    "js_runtime",
		Type:    "pub",
		Version: "2.18.0",
		Path:    "js_runtime",
		Locations: []model.Location{
			{
				Path:      filepath.Join("usr", "lib", "dart", "lib", "_internal", "js_runtime", "pubspec.yaml"),
				LayerHash: "9e0a4214e413f64ed3d287ab0048149f5e8b39c1cc83cb93919b3879ddba874a",
			},
		},
		Description: "",
		Licenses: []string{
			"BSD 3-Clause",
		},
		CPEs: []string{
			"cpe:2.3:a:js_runtime:js_runtime:2.18.0:*:*:*:*:*:*:*",
			"cpe:2.3:a:js_runtime:js-runtime:2.18.0:*:*:*:*:*:*:*",
			"cpe:2.3:a:js-runtime:js-runtime:2.18.0:*:*:*:*:*:*:*",
			"cpe:2.3:a:js-runtime:js_runtime:2.18.0:*:*:*:*:*:*:*",
		},
		PURL: model.PURL("pkg:dart/js_runtime@2.18.0"),
		Metadata: Metadata{
			"environment": []string{
				"sdk: >=2.2.2",
			},
			"name":       "js_runtime",
			"publish_to": "none",
		},
	}
	dartPackage2 = model.Package{
		Name:    "sdk_library_metadata",
		Type:    "pub",
		Version: "2.18.0",
		Path:    "sdk_library_metadata",
		Locations: []model.Location{
			{
				Path:      filepath.Join("usr", "lib", "dart", "lib", "_internal", "sdk_library_metadata", "pubspec.yaml"),
				LayerHash: "9e0a4214e413f64ed3d287ab0048149f5e8b39c1cc83cb93919b3879ddba874a",
			},
		},
		Description: "",
		Licenses: []string{
			"BSD 3-Clause",
		},
		CPEs: []string{
			"cpe:2.3:a:sdk_library_metadata:sdk_library_metadata:2.18.0:*:*:*:*:*:*:*",
			"cpe:2.3:a:sdk_library_metadata:sdk-library_metadata:2.18.0:*:*:*:*:*:*:*",
			"cpe:2.3:a:sdk_library_metadata:sdk-library-metadata:2.18.0:*:*:*:*:*:*:*",
			"cpe:2.3:a:sdk-library_metadata:sdk-library-metadata:2.18.0:*:*:*:*:*:*:*",
			"cpe:2.3:a:sdk-library_metadata:sdk_library-metadata:2.18.0:*:*:*:*:*:*:*",
			"cpe:2.3:a:sdk-library_metadata:sdk_library_metadata:2.18.0:*:*:*:*:*:*:*",
			"cpe:2.3:a:sdk-library-metadata:sdk_library_metadata:2.18.0:*:*:*:*:*:*:*",
			"cpe:2.3:a:sdk-library-metadata:sdk-library_metadata:2.18.0:*:*:*:*:*:*:*",
			"cpe:2.3:a:sdk-library-metadata:sdk-library-metadata:2.18.0:*:*:*:*:*:*:*",
		},
		PURL: model.PURL("pkg:dart/sdk_library_metadata@2.18.0"),
		Metadata: Metadata{
			"environment": []string{
				"sdk: >=2.12.0 <3.0.0",
			},
			"name":       "sdk_library_metadata",
			"publish_to": "none",
		},
	}
)

func TestParseDartPackages(t *testing.T) {
	pubspecPath := filepath.Join("..", "..", "..", "docs", "references", "dart", "pubspec.yaml")
	testLocation := model.Location{Path: pubspecPath}
	err := parseDartPackages(&testLocation)
	if err != nil {
		t.Errorf("Test Failed: Error occurred while reading Dart content.")
	}

}

func TestParseDartPackagesLock(t *testing.T) {
	pubspecLockPath := filepath.Join("..", "..", "..", "docs", "references", "dart", "pubspec.lock")
	testLocation := model.Location{Path: pubspecLockPath}
	err := parseDartPackagesLock(&testLocation)
	if err != nil {
		t.Errorf("Test Failed: Error occurred while reading Dart content.")
	}
}

func TestParseDartPURL(t *testing.T) {
	tests := []DartPurlResult{
		{&dartPackage1, model.PURL("pkg:dart/js_runtime@2.18.0")},
		{&dartPackage2, model.PURL("pkg:dart/sdk_library_metadata@2.18.0")},
	}
	for _, test := range tests {
		parseDartPURL(test._package)
		if test._package.PURL != test.expected {
			t.Errorf("Test Failed: Expected an output of %v, received: %v", test.expected, test._package.PURL)
		}
	}
}
