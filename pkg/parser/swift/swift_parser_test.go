package swift

import (
	"path/filepath"
	"testing"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/model/metadata"
)

type (
	SwiftPurlResult struct {
		_package *model.Package
		expected model.PURL
	}
)

var (
	swiftPackage1 = model.Package{
		Name:    "SDWebImage",
		Type:    "pod",
		Version: "3.5.4",
		Path:    "SDWebImage",
		Locations: []model.Location{
			{
				Path: filepath.Join("spec", "fixtures", "Podfile.lock"),
			},
		},
		Description: "",
		Licenses:    []string{},
		CPEs: []string{
			"cpe:2.3:a:SDWebImage:SDWebImage:3.5.4:*:*:*:*:*:*:*",
		},
		PURL: model.PURL("pkg:cocoapods/SDWebImage@3.5.4"),
		Metadata: metadata.PodFileLockMetadataCheckSums{
			Checksums: "571295c8acbb699505a0f68e54506147c3de9ea7",
		},
	}
	swiftPackage2 = model.Package{
		Name:    "SBJson",
		Type:    "pod",
		Version: "3.2",
		Path:    "SBJson",
		Locations: []model.Location{
			{
				Path: filepath.Join("spec", "fixtures", "Podfile.lock"),
			},
		},
		Description: "",
		Licenses:    []string{},
		CPEs: []string{
			"cpe:2.3:a:SBJson:SBJson:3.2:*:*:*:*:*:*:*",
		},
		PURL: model.PURL("pkg:cocoapods/SBJson@3.2"),
		Metadata: metadata.PodFileLockMetadataCheckSums{
			Checksums: "72d7da8a5bf6c236e87194abb10ac573a8bccbef",
		},
	}
	swiftPackage3 = model.Package{
		Name:    "SVPullToRefresh",
		Type:    "pod",
		Version: "0.4.1",
		Path:    "SVPullToRefresh",
		Locations: []model.Location{
			{
				Path: filepath.Join("spec", "fixtures", "Podfile.lock"),
			},
		},
		Description: "",
		Licenses:    []string{},
		CPEs: []string{
			"cpe:2.3:a:SVPullToRefresh:SVPullToRefresh:0.4.1:*:*:*:*:*:*:*",
		},
		PURL: model.PURL("pkg:cocoapods/SVPullToRefresh@0.4.1"),
		Metadata: metadata.PodFileLockMetadataCheckSums{
			Checksums: "d5161ebc833a38b465364412e5e307ca80bbb190",
		},
	}
)

func TestParseSwiftPackages(t *testing.T) {
	podfilePath := filepath.Join("..", "..", "..", "docs", "references", "swift", "Podfile.lock")
	testLocation := model.Location{Path: podfilePath}
	err := parseSwiftPackages(&testLocation)
	if err != nil {
		t.Errorf("Test Failed: Error occurred while reading Swift content.")
	}
}

func TestParseSwiftPURL(t *testing.T) {
	tests := []SwiftPurlResult{
		{&swiftPackage1, model.PURL("pkg:cocoapods/SDWebImage@3.5.4")},
		{&swiftPackage2, model.PURL("pkg:cocoapods/SBJson@3.2")},
		{&swiftPackage3, model.PURL("pkg:cocoapods/SVPullToRefresh@0.4.1")},
	}
	for _, test := range tests {
		parseSwiftPURL(test._package)
		if test._package.PURL != test.expected {
			t.Errorf("Test Failed: Expected an output of %v, received: %v", test.expected, test._package.PURL)
		}
	}
}
