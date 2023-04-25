package golang

import (
	"path/filepath"
	"runtime/debug"
	"testing"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/model/metadata"
)

type (
	PseudoVersionResult struct {
		vTime     string
		vRevision string
		expected  string
	}
	GetVersionResult struct {
		version  string
		settings []debug.BuildSetting
		expected string
	}

	FormatVersionResult struct {
		input    string
		expected string
	}

	ParseBuildSettingsResult struct {
		settings []debug.BuildSetting
		key      string
		expected string
	}

	FormatDevelCPEsResult struct {
		pkg      model.Package
		expected []string
	}
	GoBinMetadataResult struct {
		pkg       *model.Package
		dep       *debug.Module
		buildData *debug.BuildInfo
		expected  metadata.GoBinMetadata
	}
	InitGoBinPackageResult struct {
		pkg       *model.Package
		location  *model.Location
		buildData *debug.BuildInfo
		dep       *debug.Module
		expected  *model.Package
	}
)

var (
	goBinLocation1 = model.Location{
		Path:      filepath.Join("AppData", "Local", "Temp", "921108149", "diggity-tmp-cb5342d2-f2dd-4eb3-b6c0-0e2c9f023279", "0cd4836a36e094e1870a2e6c2578a7ad9d9cb42a7313944a6d05ab72892fc3c3", "bin", "gost"),
		LayerHash: "0cd4836a36e094e1870a2e6c2578a7ad9d9cb42a7313944a6d05ab72892fc3c3",
	}
	goBinLocation2 = model.Location{
		Path:      filepath.Join("AppData", "Local", "Temp", "4207199802", "diggity-tmp-c25a6d61-6bb0-4d23-90db-8aee8fe0516c", "1ea8aec45877fad7de4c11ccdf09146ce8ac4be9fe84c8ad036564f5d10b441b", "app", "livego"),
		LayerHash: "1ea8aec45877fad7de4c11ccdf09146ce8ac4be9fe84c8ad036564f5d10b441b",
	}
	goBinLocation3 = model.Location{
		Path:      filepath.Join("AppData", "Local", "Temp", "2199280373", "diggity-tmp-4f2825ae-68df-40d1-88e6-bfe2ffb9132f", "c2537d2244794582fe4891c11ca825226e4460825655b04bf68578aded4924b6", "bin", "golangci-lint"),
		LayerHash: "c2537d2244794582fe4891c11ca825226e4460825655b04bf68578aded4924b6",
	}

	goBinDep1 = debug.Module{
		Path:    "gitlab.com/yawning/obfs4.git",
		Version: "v0.0.0-20220204003609-77af0cba934d",
		Sum:     "h1:tJ8F7ABaQ3p3wjxwXiWSktVDgjZEXkvaRawd2rIq5ws=",
	}

	goBinDep2 = debug.Module{
		Path:    "github.com/kr/pretty",
		Version: "v0.1.0",
		Sum:     "h1:L/CwN0zerZDmRFUapSPitk6f+Q3+0za1rQkzVuMiMFI=",
	}

	goBinDep3 = debug.Module{
		Path:    "github.com/tomarrell/wrapcheck",
		Version: "v1.0.0",
		Sum:     "h1:e/6yv/rH08TZFvkYpaAMrgGbaQHVFdzaPPv4a5EIu+o=",
	}

	goBinPackage1 = model.Package{
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

	goBinPackage2 = model.Package{
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
		Metadata: metadata.GoBinMetadata{
			Architecture:     "amd64",
			Compiler:         "gc",
			OS:               "linux",
			GoCompileRelease: "go1.19.1",
			H1Digest:         "h1:L/CwN0zerZDmRFUapSPitk6f+Q3+0za1rQkzVuMiMFI=",
			Path:             "github.com/kr/pretty",
			Version:          "v0.1.0",
		},
	}

	goBinPackage3 = model.Package{
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
		Metadata: metadata.GoBinMetadata{
			Architecture:     "amd64",
			Compiler:         "gc",
			OS:               "linux",
			GoCompileRelease: "go1.16.2",
			H1Digest:         "h1:e/6yv/rH08TZFvkYpaAMrgGbaQHVFdzaPPv4a5EIu+o=",
			Path:             "github.com/tomarrell/wrapcheck",
			Version:          "v1.0.0",
		},
	}

	packageDevelUnformatted = model.Package{
		CPEs: []string{
			`cpe:2.3:a:jbariel:example-go-calculator:(devel):*:*:*:*:*:*:*`,
			`cpe:2.3:a:jbariel:example_go-calculator:(devel):*:*:*:*:*:*:*`,
			`cpe:2.3:a:jbariel:example_go_calculator:(devel):*:*:*:*:*:*:*`,
			`cpe:2.3:a:example-go-calculator:example-go-calculator:(devel):*:*:*:*:*:*:*`,
		},
	}

	packageDevelFormatted = model.Package{
		CPEs: []string{
			`cpe:2.3:a:jbariel:example-go-calculator:(\(devel\)):*:*:*:*:*:*:*`,
			`cpe:2.3:a:jbariel:example_go-calculator:(\(devel\)):*:*:*:*:*:*:*`,
			`cpe:2.3:a:jbariel:example_go_calculator:(\(devel\)):*:*:*:*:*:*:*`,
			`cpe:2.3:a:example-go-calculator:example-go-calculator:(\(devel\)):*:*:*:*:*:*:*`,
		},
	}

	goBuildSettings1 = []debug.BuildSetting{
		{Key: "-compiler", Value: "gc"},
		{Key: "CGO_ENABLED", Value: "1"},
		{Key: "CGO_CFLAGS", Value: ""},
		{Key: "CGO_CPPFLAGS", Value: ""},
		{Key: "CGO_CXXFLAGS", Value: ""},
		{Key: "CGO_LDFLAGS", Value: ""},
		{Key: "GOARCH", Value: "amd64"},
		{Key: "GOOS", Value: "linux"},
		{Key: "GOAMD64", Value: "v1"},
		{Key: "vcs.revision", Value: "6ed9a7928f649c59cdcc78b5884e4340b0a44a64"},
		{Key: "vcs.time", Value: "2022-12-02T16:40:50Z"},
	}
	goBuildSettings2 = []debug.BuildSetting{
		{Key: "-compiler", Value: "test_compiler"},
		{Key: "CGO_ENABLED", Value: "0"},
		{Key: "GOARCH", Value: "test_arch"},
		{Key: "GOOS", Value: "test_os"},
		{Key: "GOAMD64", Value: "v1"},
		{Key: "vcs.revision", Value: "ed00243a0ce2a0aee75311b06e32d33b44729689"},
		{Key: "vcs.time", Value: "2022-08-18T16:44:29Z"},
	}
	goBuildSettings3 = []debug.BuildSetting{}

	testBuildData1 = debug.BuildInfo{
		GoVersion: "go1.19.1",
		Settings:  goBuildSettings1,
	}

	testBuildData2 = debug.BuildInfo{
		GoVersion: "go1.16.2",
		Settings:  goBuildSettings1,
	}
)

func TestReadGoBinContent(t *testing.T) {
	goBinPath := filepath.Join("..", "..", "..", "docs", "references", "go", "gobin")
	testLocation := model.Location{Path: goBinPath}
	pkgs := new([]model.Package)
	err := readGoBinContent(&testLocation, pkgs)
	if err != nil {
		t.Error("Test Failed: Error occurred while reading go bin content.")
	}
}

func TestInitGoBinPackage(t *testing.T) {
	var pkg1, pkg2, pkg3 model.Package

	tests := []InitGoBinPackageResult{
		{&pkg1, &goBinLocation1, &testBuildData1, &goBinDep1, &goBinPackage1},
		{&pkg2, &goBinLocation2, &testBuildData1, &goBinDep2, &goBinPackage2},
		{&pkg3, &goBinLocation3, &testBuildData2, &goBinDep3, &goBinPackage3},
	}

	for _, test := range tests {
		output := initGoBinPackage(test.pkg, test.location, test.buildData, test.dep)
		outputMetadata := output.Metadata.(metadata.GoBinMetadata)
		expectedMetadata := test.expected.Metadata.(metadata.GoBinMetadata)

		if output.Type != test.expected.Type ||
			output.Path != test.expected.Path ||
			output.Name != test.expected.Name ||
			output.Version != test.expected.Version ||
			output.Description != test.expected.Description ||
			len(output.Licenses) != len(test.expected.Licenses) ||
			len(output.Locations) != len(test.expected.Locations) ||
			len(output.CPEs) != len(test.expected.CPEs) ||
			string(output.PURL) != string(test.expected.PURL) ||
			outputMetadata.Architecture != expectedMetadata.Architecture ||
			outputMetadata.Compiler != expectedMetadata.Compiler ||
			outputMetadata.OS != expectedMetadata.OS ||
			outputMetadata.GoCompileRelease != expectedMetadata.GoCompileRelease ||
			outputMetadata.H1Digest != expectedMetadata.H1Digest ||
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

func TestInitGoBinMetadata(t *testing.T) {
	var pkg1, pkg2, pkg3 model.Package

	tests := []GoBinMetadataResult{
		{&pkg1, &goBinDep1, &testBuildData1, metadata.GoBinMetadata{
			Architecture:     "amd64",
			Compiler:         "gc",
			OS:               "linux",
			GoCompileRelease: "go1.19.1",
			H1Digest:         "h1:tJ8F7ABaQ3p3wjxwXiWSktVDgjZEXkvaRawd2rIq5ws=",
			Path:             "gitlab.com/yawning/obfs4.git",
			Version:          "v0.0.0-20220204003609-77af0cba934d",
		}},
		{&pkg2, &goBinDep2, &testBuildData1, metadata.GoBinMetadata{
			Architecture:     "amd64",
			Compiler:         "gc",
			OS:               "linux",
			GoCompileRelease: "go1.19.1",
			H1Digest:         "h1:L/CwN0zerZDmRFUapSPitk6f+Q3+0za1rQkzVuMiMFI=",
			Path:             "github.com/kr/pretty",
			Version:          "v0.1.0",
		}},
		{&pkg3, &goBinDep3, &testBuildData2, metadata.GoBinMetadata{
			Architecture:     "amd64",
			Compiler:         "gc",
			OS:               "linux",
			GoCompileRelease: "go1.16.2",
			H1Digest:         "h1:e/6yv/rH08TZFvkYpaAMrgGbaQHVFdzaPPv4a5EIu+o=",
			Path:             "github.com/tomarrell/wrapcheck",
			Version:          "v1.0.0",
		}},
	}

	for _, test := range tests {
		initGoBinMetadata(test.pkg, test.dep, test.buildData)

		outputMetadata := test.pkg.Metadata.(metadata.GoBinMetadata)
		expectedMetadata := test.expected

		if outputMetadata != expectedMetadata {
			t.Errorf("Test Failed: Expected output of %v, received: %v", test.expected, test.pkg.Metadata)
		}
	}
}
func TestFormatDevelCPEs(t *testing.T) {
	result1 := packageDevelFormatted.CPEs
	tests := []FormatDevelCPEsResult{
		{packageDevelUnformatted, result1},
	}

	for _, test := range tests {
		output := formatDevelCPEs(&test.pkg)
		if len(output) > 0 {
			for i := range output {
				if string(output[i]) != test.expected[i] {
					t.Errorf("Test Failed: \nIndex: %v of input %v must be equal to %v, received: %v", i, test.pkg.CPEs[i], test.expected[i], output[i])
				}
			}
		}
	}
}

func TestParseBuildSettings(t *testing.T) {
	tests := []ParseBuildSettingsResult{
		{goBuildSettings1, goArch, "amd64"},
		{goBuildSettings1, goOS, "linux"},
		{goBuildSettings1, goCompiler, "gc"},
		{goBuildSettings2, goArch, "test_arch"},
		{goBuildSettings2, goOS, "test_os"},
		{goBuildSettings2, goCompiler, "test_compiler"},
		{goBuildSettings3, goArch, ""},
		{goBuildSettings3, goOS, ""},
		{goBuildSettings3, goCompiler, ""},
	}

	for _, test := range tests {
		if output := parseBuildSettings(test.settings, test.key); output != test.expected {
			t.Errorf("Test Failed: Key %v must have output of %v, received: %v", test.key, test.expected, output)
		}
	}
}

func TestFormatVersion(t *testing.T) {
	tests := []FormatVersionResult{
		{"(devel)", "devel"},
		{"(test)", "test"},
		{"v0.3.8", "v0.3.8"},
		{"v0.0.0-20220516162934-403b01795ae8", "v0.0.0-20220516162934-403b01795ae8"},
		{"", ""},
	}

	for _, test := range tests {
		if output := formatVersion(test.input); output != test.expected {
			t.Errorf("Test Failed: Input %v must have output of %v, received: %v", test.input, test.expected, output)
		}
	}
}

func TestGetVersion(t *testing.T) {
	tests := []GetVersionResult{
		{"v1.0", goBuildSettings1, "v1.0"},
		{"v2.0", goBuildSettings2, "v2.0"},
		{"(devel)", goBuildSettings1, "v0.0.0-20221202164050-6ed9a7928f64"},
		{"(devel)", goBuildSettings2, "v0.0.0-20220818164429-ed00243a0ce2"},
		{"(devel)", goBuildSettings3, "(devel)"},
	}

	for _, test := range tests {
		if output := getVersion(test.version, test.settings); output != test.expected {
			t.Errorf("Test Failed: Expected output of %v, received: %v", test.expected, output)
		}
	}
}

func TestPseudoVersion(t *testing.T) {
	tests := []PseudoVersionResult{
		{"xyz2022test", "abcdefghijklmnopqrsuvwxyz", "v0.0.0-2022-abcdefghijkl"},
		{"2022-12-02T16:40:50Z", "6ed9a7928f649c59cdcc78b5884e4340b0a44a64", "v0.0.0-20221202164050-6ed9a7928f64"},
		{"2022-08-18T16:44:29Z", "ed00243a0ce2a0aee75311b06e32d33b44729689", "v0.0.0-20220818164429-ed00243a0ce2"},
	}

	for _, test := range tests {
		if output := pseudoVersion(test.vTime, test.vRevision); output != test.expected {
			t.Errorf("Test Failed: Expected output of %v, received: %v", test.expected, output)
		}
	}

}
