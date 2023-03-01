package parser

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/carbonetes/diggity/internal/model"
	"github.com/carbonetes/diggity/internal/model/metadata"

	"golang.org/x/mod/modfile"
	"golang.org/x/mod/module"
)

type (
	SplitPathResult struct {
		input    string
		expected []string
	}
	GoModPurlResult struct {
		_package *model.Package
		expected model.PURL
	}

	GoModMetadataResult struct {
		_package *model.Package
		modPkg   interface{}
		expected metadata.GoModMetadata
	}

	InitGoModPackageResult struct {
		_package *model.Package
		location *model.Location
		modPkg   interface{}
		expected *model.Package
	}
)

var (
	goModLocation1 = model.Location{
		Path:      filepath.Join("AppData", "Local", "Temp", "921108149", "diggity-tmp-cb5342d2-f2dd-4eb3-b6c0-0e2c9f023279", "0cd4836a36e094e1870a2e6c2578a7ad9d9cb42a7313944a6d05ab72892fc3c3", "bin", "gost"),
		LayerHash: "0cd4836a36e094e1870a2e6c2578a7ad9d9cb42a7313944a6d05ab72892fc3c3",
	}
	goModLocation2 = model.Location{
		Path:      filepath.Join("AppData", "Local", "Temp", "4207199802", "diggity-tmp-c25a6d61-6bb0-4d23-90db-8aee8fe0516c", "1ea8aec45877fad7de4c11ccdf09146ce8ac4be9fe84c8ad036564f5d10b441b", "app", "livego"),
		LayerHash: "1ea8aec45877fad7de4c11ccdf09146ce8ac4be9fe84c8ad036564f5d10b441b",
	}
	goModLocation3 = model.Location{
		Path:      filepath.Join("AppData", "Local", "Temp", "2199280373", "diggity-tmp-4f2825ae-68df-40d1-88e6-bfe2ffb9132f", "c2537d2244794582fe4891c11ca825226e4460825655b04bf68578aded4924b6", "bin", "golangci-lint"),
		LayerHash: "c2537d2244794582fe4891c11ca825226e4460825655b04bf68578aded4924b6",
	}

	modPkg1 = modfile.Require{
		Mod: module.Version{
			Path:    "gitlab.com/yawning/obfs4.git",
			Version: "v0.0.0-20220204003609-77af0cba934d",
		},
	}

	modPkg2 = modfile.Replace{
		New: module.Version{
			Path:    "github.com/kr/pretty",
			Version: "v0.1.0",
		},
	}

	modPkg3 = modfile.Require{
		Mod: module.Version{
			Path:    "github.com/tomarrell/wrapcheck",
			Version: "v1.0.0",
		},
	}

	modFile = modfile.File{
		Exclude: []*modfile.Exclude{
			{Mod: module.Version{
				Path: "github.com/test/to-exclude01",
			}},
			{Mod: module.Version{
				Path: "github.com/test/to-exclude02",
			}},
		},
	}

	goModPackage1 = model.Package{
		ID:      "39bbf6e0-cc22-493e-86fb-1f3b4b5d1cf9",
		Name:    "gitlab.com/yawning/obfs4.git",
		Type:    goModule,
		Version: "v0.0.0-20220204003609-77af0cba934d",
		Path:    "gitlab.com/yawning/obfs4.git",
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
		PURL: model.PURL("pkg:golang/gitlab.com/yawning/obfs4.git@v0.0.0-20220204003609-77af0cba934d"),
		Metadata: metadata.GoModMetadata{
			Path:    "gitlab.com/yawning/obfs4.git",
			Version: "v0.0.0-20220204003609-77af0cba934d",
		},
	}

	goModPackage2 = model.Package{
		ID:      "885d17ab-295f-4e30-9a22-0ee9c62a85fe",
		Name:    "github.com/kr/pretty",
		Type:    goModule,
		Version: "v0.1.0",
		Path:    "github.com/kr/pretty",
		Locations: []model.Location{
			{
				Path:      filepath.Join("app", "livego"),
				LayerHash: "1ea8aec45877fad7de4c11ccdf09146ce8ac4be9fe84c8ad036564f5d10b441b",
			},
		},
		Description: "",
		Licenses:    []string{},
		CPEs: []string{
			"cpe:2.3:a:kr:pretty:v0.1.0:*:*:*:*:*:*:*",
			"cpe:2.3:a:pretty:pretty:v0.1.0:*:*:*:*:*:*:*",
		},
		PURL: model.PURL("pkg:golang/github.com/kr/pretty@v0.1.0"),
		Metadata: metadata.GoModMetadata{
			Path:    "github.com/kr/pretty",
			Version: "v0.1.0",
		},
	}

	goModPackage3 = model.Package{
		ID:      "073e8805-bb2d-4265-b6dd-e9b9b49a7358",
		Name:    "github.com/tomarrell/wrapcheck",
		Type:    goModule,
		Version: "v1.0.0",
		Path:    "github.com/tomarrell/wrapcheck",
		Locations: []model.Location{
			{
				Path:      filepath.Join("bin", "golangci-lint"),
				LayerHash: "c2537d2244794582fe4891c11ca825226e4460825655b04bf68578aded4924b6",
			},
		},
		Description: "",
		Licenses:    []string{},
		CPEs: []string{
			"cpe:2.3:a:tomarrell:wrapcheck:v1.0.0:*:*:*:*:*:*:*",
			"cpe:2.3:a:wrapcheck:wrapcheck:v1.0.0:*:*:*:*:*:*:*",
		},
		PURL: model.PURL("pkg:golang/github.com/tomarrell/wrapcheck@v1.0.0"),
		Metadata: metadata.GoModMetadata{
			Path:    "github.com/tomarrell/wrapcheck",
			Version: "v1.0.0",
		},
	}
)

func TestReadGoModContent(t *testing.T) {
	goModPath := filepath.Join("..", "..", "docs", "references", "go", "go.mod")
	testLocation := model.Location{Path: goModPath}
	err := readGoModContent(&testLocation)
	if err != nil {
		t.Error("Test Failed: Error occurred while reading go mod content.")
	}
}

func TestInitGoModPackage(t *testing.T) {
	var _package1, _package2, _package3 model.Package

	tests := []InitGoModPackageResult{
		{&_package1, &goModLocation1, &modPkg1, &goModPackage1},
		{&_package2, &goModLocation2, &modPkg2, &goModPackage2},
		{&_package3, &goModLocation3, &modPkg3, &goModPackage3},
	}

	for _, test := range tests {
		output := initGoModPackage(test._package, test.location, test.modPkg)
		outputMetadata := output.Metadata.(metadata.GoModMetadata)
		expectedMetadata := test.expected.Metadata.(metadata.GoModMetadata)

		if output.Type != test.expected.Type ||
			output.Path != test.expected.Path ||
			output.Name != test.expected.Name ||
			output.Version != test.expected.Version ||
			output.Description != test.expected.Description ||
			len(output.Licenses) != len(test.expected.Licenses) ||
			len(output.Locations) != len(test.expected.Locations) ||
			len(output.CPEs) != len(test.expected.CPEs) ||
			string(output.PURL) != string(test.expected.PURL) ||
			outputMetadata.Path != expectedMetadata.Path ||
			outputMetadata.Version != expectedMetadata.Version {
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
	}
}

func TestInitGoModMetadata(t *testing.T) {
	var _package1, _package2, _package3 model.Package

	tests := []GoModMetadataResult{
		{&_package1, &modPkg1, metadata.GoModMetadata{
			Path:    "gitlab.com/yawning/obfs4.git",
			Version: "v0.0.0-20220204003609-77af0cba934d",
		}},
		{&_package2, &modPkg2, metadata.GoModMetadata{
			Path:    "github.com/kr/pretty",
			Version: "v0.1.0",
		}},
		{&_package3, &modPkg3, metadata.GoModMetadata{
			Path:    "github.com/tomarrell/wrapcheck",
			Version: "v1.0.0",
		}},
	}

	for _, test := range tests {
		initGoModMetadata(test._package, test.modPkg)

		outputMetadata := test._package.Metadata.(metadata.GoModMetadata)
		expectedMetadata := test.expected

		if outputMetadata.Path != expectedMetadata.Path ||
			outputMetadata.Version != expectedMetadata.Version {
			t.Errorf("Test Failed: Expected output of %v, received: %v", test.expected, test._package.Metadata)
		}
	}
}
func TestGoModPackageURL(t *testing.T) {
	_package1 := model.Package{
		Name:    goModPackage1.Name,
		Version: goModPackage1.Version,
	}
	_package2 := model.Package{
		Name:    goModPackage2.Name,
		Version: goModPackage2.Version,
	}
	_package3 := model.Package{
		Name:    goModPackage3.Name,
		Version: goModPackage3.Version,
	}

	tests := []GoModPurlResult{
		{&_package1, model.PURL("pkg:golang/gitlab.com/yawning/obfs4.git@v0.0.0-20220204003609-77af0cba934d")},
		{&_package2, model.PURL("pkg:golang/github.com/kr/pretty@v0.1.0")},
		{&_package3, model.PURL("pkg:golang/github.com/tomarrell/wrapcheck@v1.0.0")},
	}

	for _, test := range tests {
		parseGoPackageURL(test._package)
		if test._package.PURL != test.expected {
			t.Errorf("Test Failed: Expected an output of %v, received: %v", test.expected, test._package.PURL)
		}
	}
}

func TestCleanExcluded(t *testing.T) {
	p1 := new(model.Package)
	p2 := new(model.Package)
	p1.Name = "github.com/test/to-exclude01"
	p2.Name = "github.com/test/to-exclude02"
	packages := []*model.Package{&goModPackage1, &goModPackage2, &goModPackage3, p1, p2}
	expected := []*model.Package{&goModPackage1, &goModPackage2, &goModPackage3}

	packages = cleanExcluded(packages, &modFile)
	if len(packages) != len(expected) {
		t.Errorf("Test Failed: Slice length must be equal with the expected result. Expected: %v, Received: %v", len(expected), len(packages))
	}

	for _, pkg := range packages {
		if strings.Contains(pkg.Name, "github.com/test/to-exclude") {
			t.Errorf("Test Failed: Excluded Package must be removed.")
		}
	}
}

func TestSplitPath(t *testing.T) {
	tests := []SplitPathResult{
		{"golang.org/x/crypto", []string{"golang.org", "x", "crypto"}},
		{"golang.org/x/net", []string{"golang.org", "x", "net"}},
		{"golang.org/x/sys", []string{"golang.org", "x", "sys"}},
		{"golang.org/x/text", []string{"golang.org", "x", "text"}},
	}

	for _, test := range tests {
		output := splitPath(test.input)
		if len(output) > 0 {
			for i := range output {
				if output[i] != test.expected[i] {
					t.Errorf("Test Failed: Index %v of input %v must be equal to %v, received: %v", i, test.input, test.expected[i], output[i])
				}
			}
		}
	}
}
