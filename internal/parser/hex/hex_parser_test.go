package parser

import (
	"path/filepath"
	"testing"

	"github.com/carbonetes/diggity/internal/model"
)

type (
	HexPurlResult struct {
		_package *model.Package
		expected model.PURL
	}
)

var (
	hexPackage1 = model.Package{
		Name:    "bypass",
		Type:    "hex",
		Version: "1.0.0",
		Path:    "bypass",
		Locations: []model.Location{
			{
				Path: filepath.Join("hex_core", "mix.lock"),
			},
		},
		Description: "",
		Licenses:    []string{},
		CPEs: []string{
			"cpe:2.3:a:bypass:bypass:1.0.0:*:*:*:*:*:*:*",
		},
		PURL: model.PURL("pkg:hex/bypass@1.0.0"),
		Metadata: Metadata{
			Name:       "bypass",
			Version:    "1.0.0",
			PkgHash:    "b78b3dcb832a71aca5259c1a704b2e14b55fd4e1327ff942598b4e7d1a7ad83d",
			PkgHashExt: "5a1dc855dfcc86160458c7a70d25f65d498bd8012bd4c06a8d3baa368dda3c45",
		},
	}
	hexPackage2 = model.Package{
		Name:    "idna",
		Type:    "hex",
		Version: "6.1.1",
		Path:    "idna",
		Locations: []model.Location{
			{
				Path: filepath.Join("hex_core", "rebar.lock"),
			},
		},
		Description: "",
		Licenses:    []string{},
		CPEs: []string{
			"cpe:2.3:a:idna:idna:6.1.1:*:*:*:*:*:*:*",
		},
		PURL: model.PURL("pkg:hex/idna@6.1.1"),
		Metadata: Metadata{
			Name:    "idna",
			Version: "6.1.1",
		},
	}
	hexPackage3 = model.Package{
		Name:    "plug_cowboy",
		Type:    "hex",
		Version: "2.1.3",
		Path:    "plug_cowboy",
		Locations: []model.Location{
			{
				Path: filepath.Join("hex_core", "mix.lock"),
			},
		},
		Description: "",
		Licenses:    []string{},
		CPEs: []string{
			"cpe:2.3:a:plug_cowboy:plug_cowboy:2.1.3:*:*:*:*:*:*:*",
			"cpe:2.3:a:plug_cowboy:plug-cowboy:2.1.3:*:*:*:*:*:*:*",
			"cpe:2.3:a:plug-cowboy:plug-cowboy:2.1.3:*:*:*:*:*:*:*",
			"cpe:2.3:a:plug-cowboy:plug_cowboy:2.1.3:*:*:*:*:*:*:*",
		},
		PURL: model.PURL("pkg:hex/plug_cowboy@2.1.3"),
		Metadata: Metadata{
			Name:       "plug_cowboy",
			Version:    "2.1.3",
			PkgHash:    "38999a3e85e39f0e6bdfdf820761abac61edde1632cfebbacc445cdcb6ae1333",
			PkgHashExt: "056f41f814dbb38ea44613e0f613b3b2b2f2c6afce64126e252837669eba84db",
		},
	}
)

func TestParseHexRebarPackages(t *testing.T) {
	rebarLockPath := filepath.Join("..", "..", "..", "docs", "references", "hex", "rebar.lock")
	testLocation := model.Location{Path: rebarLockPath}
	err := parseHexRebarPacakges(&testLocation)
	if err != nil {
		t.Errorf("Test Failed: Error occurred while reading Hex - rebar content.")
	}
}

func TestParseHexMixPackages(t *testing.T) {
	mixLockPath := filepath.Join("..", "..", "..", "docs", "references", "hex", "mix.lock")
	testLocation := model.Location{Path: mixLockPath}
	err := parseHexMixPackages(&testLocation)
	if err != nil {
		t.Errorf("Test Failed: Error occurred while reading Hex - mix content.")
	}
}

func TestParseHexPURL(t *testing.T) {
	tests := []HexPurlResult{
		{&hexPackage1, model.PURL(hexPackage1.PURL)},
		{&hexPackage2, model.PURL(hexPackage2.PURL)},
		{&hexPackage3, model.PURL(hexPackage3.PURL)},
	}
	for _, test := range tests {
		parseHexPURL(test._package)
		if test._package.PURL != test.expected {
			t.Errorf("Test Failed: Expected an output of %v, received: %v", test.expected, test._package.PURL)
		}
	}
}
