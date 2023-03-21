package parser

import (
	"path/filepath"
	"testing"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/model/metadata"
)

type (
	HackagePackageResult struct {
		location *model.Location
		dep      string
		url      string
		expected *model.Package
	}

	HackagePurlResult struct {
		_package *model.Package
		expected model.PURL
	}

	ParseExtraDepResult struct {
		dep     string
		name    string
		version string
		pkgHash string
		size    string
		rev     string
	}

	FormatCabalPackageResult struct {
		input    string
		expected string
	}
)

var (
	haskellPackage1 = model.Package{
		Name:    "Cabal",
		Type:    hackage,
		Version: "3.8.1.0",
		Path:    "Cabal",
		Locations: []model.Location{
			{
				Path: stackYaml,
			},
		},
		Description: "",
		Licenses:    []string{},
		CPEs: []string{
			"cpe:2.3:a:Cabal:Cabal:3.8.1.0:*:*:*:*:*:*:*",
		},
		PURL: model.PURL("pkg:hackage/Cabal@3.8.1.0"),
		Metadata: metadata.HackageMetadata{
			Name:    "Cabal",
			Version: "3.8.1.0",
			PkgHash: "sha256:155d64beeecbae2b19e5d67844532494af88bc8795d4db4146a0c29296f59967",
			Size:    "12220",
		},
	}
	haskellPackage2 = model.Package{
		Name:    "rio-prettyprint",
		Type:    hackage,
		Version: "0.1.4.0",
		Path:    "rio-prettyprint",
		Locations: []model.Location{
			{
				Path: stackYamlLock,
			},
		},
		Description: "",
		Licenses:    []string{},
		CPEs: []string{
			"cpe:2.3:a:rio-prettyprint:rio-prettyprint:0.1.4.0:*:*:*:*:*:*:*",
			"cpe:2.3:a:rio-prettyprint:rio_prettyprint:0.1.4.0:*:*:*:*:*:*:*",
			"cpe:2.3:a:rio_prettyprint:rio_prettyprint:0.1.4.0:*:*:*:*:*:*:*",
			"cpe:2.3:a:rio_prettyprint:rio-prettyprint:0.1.4.0:*:*:*:*:*:*:*",
			"cpe:2.3:a:rio:rio-prettyprint:0.1.4.0:*:*:*:*:*:*:*",
			"cpe:2.3:a:rio:rio_prettyprint:0.1.4.0:*:*:*:*:*:*:*",
		},
		PURL: model.PURL("pkg:hackage/rio-prettyprint@0.1.4.0"),
		Metadata: metadata.HackageMetadata{
			Name:        "rio-prettyprint",
			Version:     "0.1.4.0",
			PkgHash:     "sha256:1f8eb3ead0ef33d3736d53e1de5e9b2c91a0c207cdca23321bd74c401e85f23a",
			Size:        "1301",
			SnapshotURL: "https://raw.githubusercontent.com/commercialhaskell/stackage-snapshots/master/lts/20/0.yaml",
		},
	}
	haskellPackage3 = model.Package{
		Name:    "BoundedChan",
		Type:    hackage,
		Version: "1.0.3.0",
		Path:    "BoundedChan",
		Locations: []model.Location{
			{
				Path: cabalFreeze,
			},
		},
		Description: "",
		Licenses:    []string{},
		CPEs: []string{
			"cpe:2.3:a:BoundedChan:BoundedChan:1.0.3.0:*:*:*:*:*:*:*",
		},
		PURL: model.PURL("pkg:hackage/BoundedChan@1.0.3.0"),
		Metadata: metadata.HackageMetadata{
			Name:    "BoundedChan",
			Version: "1.0.3.0",
		},
	}
)

func TestReadStackContent(t *testing.T) {
	stackPath := filepath.Join("..", "..", "..", "docs", "references", "hackage", stackYaml)
	testLocation := model.Location{Path: stackPath}
	err := readStackContent(&testLocation)
	if err != nil {
		t.Error("Test Failed: Error occurred while reading stack.yaml content.")
	}
}

func TestReadStackLockContent(t *testing.T) {
	stackLockPath := filepath.Join("..", "..", "..", "docs", "references", "hackage", stackYamlLock)
	testLocation := model.Location{Path: stackLockPath}
	err := readStackLockContent(&testLocation)
	if err != nil {
		t.Error("Test Failed: Error occurred while reading stack.yaml.lock content.")
	}
}

func TestReadCabalFreezeContent(t *testing.T) {
	cabalFreezePath := filepath.Join("..", "..", "..", "docs", "references", "hackage", cabalFreeze)
	testLocation := model.Location{Path: cabalFreezePath}
	err := readCabalFreezeContent(&testLocation)
	if err != nil {
		t.Error("Test Failed: Error occurred while reading cabal.project.freeze content.")
	}
}

func TestInitHackagePackage(t *testing.T) {
	tests := []HackagePackageResult{
		{&model.Location{Path: stackYaml}, "Cabal-3.8.1.0@sha256:155d64beeecbae2b19e5d67844532494af88bc8795d4db4146a0c29296f59967,12220", "", &haskellPackage1},
		{&model.Location{Path: stackYamlLock}, "rio-prettyprint-0.1.4.0@sha256:1f8eb3ead0ef33d3736d53e1de5e9b2c91a0c207cdca23321bd74c401e85f23a,1301", "https://raw.githubusercontent.com/commercialhaskell/stackage-snapshots/master/lts/20/0.yaml", &haskellPackage2},
		{&model.Location{Path: cabalFreeze}, "BoundedChan-1.0.3.0", "", &haskellPackage3},
	}

	for _, test := range tests {
		output := initHackagePackage(test.location, test.dep, test.url)
		outputMetadata := output.Metadata.(metadata.HackageMetadata)
		expectedMetadata := test.expected.Metadata.(metadata.HackageMetadata)

		if output.Type != test.expected.Type ||
			output.Path != test.expected.Path ||
			output.Name != test.expected.Name ||
			output.Version != test.expected.Version ||
			output.Description != test.expected.Description ||
			len(output.Licenses) != len(test.expected.Licenses) ||
			len(output.Locations) != len(test.expected.Locations) ||
			len(output.CPEs) != len(test.expected.CPEs) ||
			string(output.PURL) != string(test.expected.PURL) ||
			outputMetadata.Version != expectedMetadata.Version ||
			outputMetadata.PkgHash != expectedMetadata.PkgHash ||
			outputMetadata.Revision != expectedMetadata.Revision ||
			outputMetadata.SnapshotURL != expectedMetadata.SnapshotURL {
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

func TestParseHackageURL(t *testing.T) {
	_package1 := model.Package{
		Name:    haskellPackage1.Name,
		Version: haskellPackage1.Version,
	}
	_package2 := model.Package{
		Name:    haskellPackage2.Name,
		Version: haskellPackage2.Version,
	}
	_package3 := model.Package{
		Name:    haskellPackage3.Name,
		Version: haskellPackage3.Version,
	}

	tests := []HackagePurlResult{
		{&_package1, model.PURL(haskellPackage1.PURL)},
		{&_package2, model.PURL(haskellPackage2.PURL)},
		{&_package3, model.PURL(haskellPackage3.PURL)},
	}
	for _, test := range tests {
		parseHackageURL(test._package)
		if test._package.PURL != test.expected {
			t.Errorf("Test Failed: Expected an output of %v, received: %v", test.expected, test._package.PURL)
		}
	}
}

func TestParseExtraDep(t *testing.T) {
	tests := []ParseExtraDepResult{
		{"Cabal-3.8.1.0@sha256:155d64beeecbae2b19e5d67844532494af88bc8795d4db4146a0c29296f59967,12220",
			"Cabal", "3.8.1.0", "sha256:155d64beeecbae2b19e5d67844532494af88bc8795d4db4146a0c29296f59967", "12220", ""},
		{"fsnotify-0.4.1.0@sha256:44540beabea36aeeef930aa4d5f28091d431904bc9923b6ac4d358831c651235,2854",
			"fsnotify", "0.4.1.0", "sha256:44540beabea36aeeef930aa4d5f28091d431904bc9923b6ac4d358831c651235", "2854", ""},
		{"rio-prettyprint-0.1.4.0@sha256:1f8eb3ead0ef33d3736d53e1de5e9b2c91a0c207cdca23321bd74c401e85f23a,1301",
			"rio-prettyprint", "0.1.4.0", "sha256:1f8eb3ead0ef33d3736d53e1de5e9b2c91a0c207cdca23321bd74c401e85f23a", "1301", ""},
		{"Cabal-syntax-3.8.1.0@sha256:4936765e9a7a8ecbf8fdbe9067f6d972bc0299220063abb2632a9950af64b966,7619",
			"Cabal-syntax", "3.8.1.0", "sha256:4936765e9a7a8ecbf8fdbe9067f6d972bc0299220063abb2632a9950af64b966", "7619", ""},
		{"acme-missiles-0.3@rev:0", "acme-missiles", "0.3", "", "", "rev:0"},
	}

	for _, test := range tests {
		name, version, pkgHash, size, rev := parseExtraDep(test.dep)
		if name != test.name {
			t.Errorf("Test Failed:\n Expected output of %v \n, Received: %v \n", test.name, name)
		}
		if version != test.version {
			t.Errorf("Test Failed:\n Expected output of %v \n, Received: %v \n", test.version, version)
		}
		if pkgHash != test.pkgHash {
			t.Errorf("Test Failed:\n Expected output of %v \n, Received: %v \n", test.pkgHash, pkgHash)
		}
		if size != test.size {
			t.Errorf("Test Failed:\n Expected output of %v \n, Received: %v \n", test.size, size)
		}
		if rev != test.rev {
			t.Errorf("Test Failed:\n Expected output of %v \n, Received: %v \n", test.rev, rev)
		}
	}
}

func TestFormatCabalPackage(t *testing.T) {
	tests := []FormatCabalPackageResult{
		{"any.test ==1.0.0", "test-1.0.0"},
		{"any.ANum ==0.2.0.2", "ANum-0.2.0.2"},
		{"any.Allure ==0.11.0.0", "Allure-0.11.0.0"},
		{"any.BNFC-meta ==0.6.1", "BNFC-meta-0.6.1"},
		{"any.BoundedChan ==1.0.3.0", "BoundedChan-1.0.3.0"},
	}

	for _, test := range tests {
		if output := formatCabalPackage(test.input); output != test.expected {
			t.Errorf("Test Failed:\n Expected output of %v \n, Received: %v \n", test.expected, output)
		}
	}
}
