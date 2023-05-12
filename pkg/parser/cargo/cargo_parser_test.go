package cargo

// import (
// 	"path/filepath"
// 	"testing"

// 	"github.com/carbonetes/diggity/pkg/model"
// 	"github.com/carbonetes/diggity/pkg/model/metadata"
// )

// type (
// 	RustPackageResult struct {
// 		metadata map[string]interface{}
// 		expected *model.Package
// 	}

// 	CargoMetadataResult struct {
// 		pkg      *model.Package
// 		metadata map[string]interface{}
// 		expected metadata.CargoMetadata
// 	}

// 	RustPurlResult struct {
// 		pkg      *model.Package
// 		expected model.PURL
// 	}

// 	FormatDependenciesResult struct {
// 		input    string
// 		expected []string
// 	}
// )

// var (
// 	rustMetadata1 = map[string]interface{}{
// 		"name":         "addr2line",
// 		"version":      "0.17.0",
// 		"source":       "registry+https://github.com/rust-lang/crates.io-index",
// 		"checksum":     "b9ecd88a8c8378ca913a680cd98f0f13ac67383d35993f86c90a70e3f137816b",
// 		"dependencies": `[ "compiler_builtins", "gimli", "rustc-std-workspace-alloc", "rustc-std-workspace-core", ]`,
// 	}
// 	rustMetadata2 = map[string]interface{}{
// 		"name":         "getrandom",
// 		"version":      "0.2.8",
// 		"source":       "registry+https://github.com/rust-lang/crates.io-index",
// 		"checksum":     "c05aeb6a22b8f62540c194aac980f2115af067bfe15a0734d7277a768d396b31",
// 		"dependencies": `[ "cfg-if", "js-sys", "libc", "wasi", "wasm-bindgen", ]`,
// 	}
// 	rustMetadata3 = map[string]interface{}{
// 		"name":         "zeroize",
// 		"version":      "1.5.7",
// 		"source":       "registry+https://github.com/rust-lang/crates.io-index",
// 		"checksum":     "c394b5bd0c6f669e7275d9c20aa90ae064cb22e75a1cad54e1b34088034b149f",
// 		"dependencies": `[ "proc-macro2", "quote", "syn", "synstructure", ]`,
// 	}

// 	rustPackage1 = model.Package{
// 		Name:    "addr2line",
// 		Type:    rustCrate,
// 		Version: "0.17.0",
// 		Path:    "addr2line",
// 		Locations: []model.Location{
// 			{
// 				Path: "Cargo.lock",
// 			},
// 		},
// 		Description: "",
// 		Licenses:    []string{},
// 		CPEs: []string{
// 			"cpe:2.3:a:addr2line:addr2line:0.17.0:*:*:*:*:*:*:*",
// 		},
// 		PURL: model.PURL("pkg:cargo/addr2line@0.17.0"),
// 		Metadata: metadata.CargoMetadata{
// 			Name:     "addr2line",
// 			Version:  "0.17.0",
// 			Source:   "registry+https://github.com/rust-lang/crates.io-index",
// 			Checksum: "b9ecd88a8c8378ca913a680cd98f0f13ac67383d35993f86c90a70e3f137816b",
// 			Dependencies: []string{
// 				"compiler_builtins",
// 				"gimli",
// 				"rustc-std-workspace-alloc",
// 				"rustc-std-workspace-core",
// 			},
// 		},
// 	}

// 	rustPackage2 = model.Package{
// 		Name:    "getrandom",
// 		Type:    rustCrate,
// 		Version: "0.2.8",
// 		Path:    "getrandom",
// 		Locations: []model.Location{
// 			{
// 				Path: "Cargo.lock",
// 			},
// 		},
// 		Description: "",
// 		Licenses:    []string{},
// 		CPEs: []string{
// 			"cpe:2.3:a:getrandom:getrandom:0.2.8:*:*:*:*:*:*:*",
// 		},
// 		PURL: model.PURL("pkg:cargo/getrandom@0.2.8"),
// 		Metadata: metadata.CargoMetadata{
// 			Name:     "getrandom",
// 			Version:  "0.2.8",
// 			Source:   "registry+https://github.com/rust-lang/crates.io-index",
// 			Checksum: "c05aeb6a22b8f62540c194aac980f2115af067bfe15a0734d7277a768d396b31",
// 			Dependencies: []string{
// 				"cfg-if",
// 				"js-sys",
// 				"libc",
// 				"wasi",
// 				"wasm-bindgen",
// 			},
// 		},
// 	}

// 	rustPackage3 = model.Package{
// 		Name:    "zeroize",
// 		Type:    rustCrate,
// 		Version: "1.5.7",
// 		Path:    "zeroize",
// 		Locations: []model.Location{
// 			{
// 				Path: "Cargo.lock",
// 			},
// 		},
// 		Description: "",
// 		Licenses:    []string{},
// 		CPEs: []string{
// 			"cpe:2.3:a:zeroize:zeroize:1.5.7:*:*:*:*:*:*:*",
// 		},
// 		PURL: model.PURL("pkg:cargo/zeroize@1.5.7"),
// 		Metadata: metadata.CargoMetadata{
// 			Name:     "zeroize",
// 			Version:  "1.5.7",
// 			Source:   "registry+https://github.com/rust-lang/crates.io-index",
// 			Checksum: "c394b5bd0c6f669e7275d9c20aa90ae064cb22e75a1cad54e1b34088034b149f",
// 			Dependencies: []string{
// 				"proc-macro2",
// 				"quote",
// 				"syn",
// 				"synstructure",
// 			},
// 		},
// 	}
// )

// func TestReadCargoContent(t *testing.T) {
// 	cargoPath := filepath.Join("..", "..", "..", "docs", "references", "rust", cargoLock)
// 	testLocation := model.Location{Path: cargoPath}
// 	pkgs := new([]model.Package)
// 	err := readCargoContent(&testLocation, pkgs)
// 	if err != nil {
// 		t.Error("Test Failed: Error occurred while reading Cargo.lock content.")
// 	}
// }

// func TestInitRustPackage(t *testing.T) {
// 	tests := []RustPackageResult{
// 		{rustMetadata1, &rustPackage1},
// 		{rustMetadata2, &rustPackage2},
// 		{rustMetadata3, &rustPackage3},
// 	}

// 	for _, test := range tests {
// 		output := initRustPackage(&model.Location{Path: cargoLock}, test.metadata)
// 		outputMetadata := output.Metadata.(metadata.CargoMetadata)
// 		expectedMetadata := test.expected.Metadata.(metadata.CargoMetadata)

// 		if output.Type != test.expected.Type ||
// 			output.Path != test.expected.Path ||
// 			output.Name != test.expected.Name ||
// 			output.Version != test.expected.Version ||
// 			output.Description != test.expected.Description ||
// 			len(output.Licenses) != len(test.expected.Licenses) ||
// 			len(output.Locations) != len(test.expected.Locations) ||
// 			len(output.CPEs) != len(test.expected.CPEs) ||
// 			string(output.PURL) != string(test.expected.PURL) ||
// 			outputMetadata.Name != expectedMetadata.Name ||
// 			outputMetadata.Version != expectedMetadata.Version ||
// 			outputMetadata.Source != expectedMetadata.Source ||
// 			outputMetadata.Checksum != expectedMetadata.Checksum ||
// 			len(outputMetadata.Dependencies) != len(expectedMetadata.Dependencies) {
// 			t.Errorf("Test Failed:\n Expected output of %v \n, Received: %v \n", test.expected, output)
// 		}

// 		for i := range output.Licenses {
// 			if output.Licenses[i] != test.expected.Licenses[i] {
// 				t.Errorf("Test Failed:\n Expected output of %v \n, Received: %v \n", test.expected.Licenses[i], output.Licenses[i])
// 			}
// 		}
// 		for i := range output.Locations {
// 			if output.Locations[i] != test.expected.Locations[i] {
// 				t.Errorf("Test Failed:\n Expected output of %v \n, Received: %v \n", test.expected.Locations[i], output.Locations[i])
// 			}
// 		}
// 		for i := range output.CPEs {
// 			if output.CPEs[i] != test.expected.CPEs[i] {
// 				t.Errorf("Test Failed:\n Expected output of %v \n, Received: %v \n", test.expected.CPEs[i], output.CPEs[i])
// 			}
// 		}
// 	}
// }

// func TestInitCargoMetadata(t *testing.T) {
// 	var pkg1, pkg2, pkg3 model.Package

// 	tests := []CargoMetadataResult{
// 		{&pkg1, rustMetadata1, metadata.CargoMetadata{
// 			Name:     "addr2line",
// 			Version:  "0.17.0",
// 			Source:   "registry+https://github.com/rust-lang/crates.io-index",
// 			Checksum: "b9ecd88a8c8378ca913a680cd98f0f13ac67383d35993f86c90a70e3f137816b",
// 			Dependencies: []string{
// 				"compiler_builtins",
// 				"gimli",
// 				"rustc-std-workspace-alloc",
// 				"rustc-std-workspace-core",
// 			},
// 		}},
// 		{&pkg2, rustMetadata2, metadata.CargoMetadata{
// 			Name:     "getrandom",
// 			Version:  "0.2.8",
// 			Source:   "registry+https://github.com/rust-lang/crates.io-index",
// 			Checksum: "c05aeb6a22b8f62540c194aac980f2115af067bfe15a0734d7277a768d396b31",
// 			Dependencies: []string{
// 				"cfg-if",
// 				"js-sys",
// 				"libc",
// 				"wasi",
// 				"wasm-bindgen",
// 			},
// 		}},
// 		{&pkg3, rustMetadata3, metadata.CargoMetadata{
// 			Name:     "zeroize",
// 			Version:  "1.5.7",
// 			Source:   "registry+https://github.com/rust-lang/crates.io-index",
// 			Checksum: "c394b5bd0c6f669e7275d9c20aa90ae064cb22e75a1cad54e1b34088034b149f",
// 			Dependencies: []string{
// 				"proc-macro2",
// 				"quote",
// 				"syn",
// 				"synstructure",
// 			},
// 		}},
// 	}

// 	for _, test := range tests {
// 		initCargoMetadata(test.pkg, test.metadata)
// 		outputMetadata := test.pkg.Metadata.(metadata.CargoMetadata)
// 		expectedMetadata := test.expected
// 		if outputMetadata.Name != expectedMetadata.Name ||
// 			outputMetadata.Version != expectedMetadata.Version ||
// 			outputMetadata.Source != expectedMetadata.Source ||
// 			outputMetadata.Checksum != expectedMetadata.Checksum ||
// 			len(outputMetadata.Dependencies) != len(expectedMetadata.Dependencies) {
// 			t.Errorf("Test Failed: Expected output of %v, received: %v", test.expected, test.pkg.Metadata)
// 		}
// 	}
// }

// func TestParseRustPackageURL(t *testing.T) {
// 	pkg1 := model.Package{
// 		Name:    rustPackage1.Name,
// 		Version: "0.17.0",
// 	}
// 	pkg2 := model.Package{
// 		Name:    rustPackage2.Name,
// 		Version: "0.2.8",
// 	}
// 	pkg3 := model.Package{
// 		Name:    rustPackage3.Name,
// 		Version: "1.5.7",
// 	}

// 	tests := []RustPurlResult{
// 		{&pkg1, model.PURL("pkg:cargo/addr2line@0.17.0")},
// 		{&pkg2, model.PURL("pkg:cargo/getrandom@0.2.8")},
// 		{&pkg3, model.PURL("pkg:cargo/zeroize@1.5.7")},
// 	}

// 	for _, test := range tests {
// 		parseRustPackageURL(test.pkg)
// 		if test.pkg.PURL != test.expected {
// 			t.Errorf("Test Failed: Expected an output of %v, received: %v", test.expected, test.pkg.PURL)
// 		}
// 	}
// }

// func TestFormatDependencies(t *testing.T) {
// 	tests := []FormatDependenciesResult{
// 		{`[ "test",]`, []string{"test"}},
// 		{`[ "linked-hash-map",][[package]]`, []string{"linked-hash-map"}},
// 		{`[ "yoke", "zerofrom", "zerovec-derive",]`, []string{"yoke", "zerofrom", "zerovec-derive"}},
// 		{`[ "proc-macro2", "quote", "syn", "synstructure",]`, []string{"proc-macro2", "quote", "syn", "synstructure"}},
// 		{`[ "proc-macro2", "quote", "syn", "synstructure",][[package]]`, []string{"proc-macro2", "quote", "syn", "synstructure"}},
// 	}

// 	for _, test := range tests {
// 		output := formatDependencies(test.input)

// 		if len(output) != len(test.expected) {
// 			t.Errorf("Test Failed: Slice length must be equal with the expected result. Expected: %v, Received: %v", len(test.expected), len(output))
// 		}

// 		for i, d := range output {
// 			if test.expected[i] != d {
// 				t.Errorf("Test Failed: Expected output of %v, received: %v", test.expected[i], d)
// 			}
// 		}
// 	}
// }
