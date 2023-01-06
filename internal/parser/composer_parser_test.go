package parser

import (
	"path/filepath"
	"testing"

	"github.com/carbonetes/diggity/internal/model"
)

type (
	PhpPurlResult struct {
		_package *model.Package
		expected model.PURL
	}
)

var (
	phpPackage1 = model.Package{
		Name:    "league/tactician",
		Type:    phpType,
		Version: "v1.1.0",
		Path:    "league/tactician",
		Locations: []model.Location{
			{
				Path:      filepath.Join("opt", "phpdoc", "composer.lock"),
				LayerHash: "2a3251e94a5184b3c5f4efbc0c8df91cf8479af3745941c9d9102298d258b83",
			},
		},
		Description: "A small, flexible command bus. Handy for building service layers.",
		Licenses: []string{
			"MIT",
		},
		CPEs: []string{
			"cpe:2.3:a:league:tactician:v1.1.0:*:*:*:*:*:*:*",
			"cpe:2.3:a:tactician:tactician:v1.1.0:*:*:*:*:*:*:*",
		},
		PURL: model.PURL("pkg:composer/league/tactician@v1.1.0"),
	}

	phpPackage2 = model.Package{
		Name:    "phpdocumentor/reflection",
		Type:    phpType,
		Version: "5.2.0",
		Path:    "phpdocumentor/reflection",
		Locations: []model.Location{
			{
				Path:      filepath.Join("opt", "phpdoc", "composer.lock"),
				LayerHash: "12a3251e94a5184b3c5f4efbc0c8df91cf8479af3745941c9d9102298d258b83",
			},
		},
		Description: "Reflection library to do Static Analysis for PHP Projects",
		Licenses: []string{
			"MIT",
		},
		CPEs: []string{
			"cpe:2.3:a:phpdocumentor:reflection:5.2.0:*:*:*:*:*:*:*",
			"cpe:2.3:a:reflection:reflection:5.2.0:*:*:*:*:*:*:*",
		},
		PURL: model.PURL("pkg:composer/phpdocumentor/reflection@5.2.0"),
	}

	phpPackage3 = model.Package{
		Name:    "symfony/framework-bundle",
		Type:    phpType,
		Version: "v5.4.10",
		Path:    "symfony/framework-bundle",
		Locations: []model.Location{
			{
				Path:      filepath.Join("opt", "phpdoc", "composer.lock"),
				LayerHash: "12a3251e94a5184b3c5f4efbc0c8df91cf8479af3745941c9d9102298d258b83",
			},
		},
		Description: "Provides a tight integration between Symfony components and the Symfony full-stack framework",
		Licenses: []string{
			"MIT",
		},
		CPEs: []string{
			"cpe:2.3:a:symfony:framework-bundle:v5.4.10:*:*:*:*:*:*:*",
			"cpe:2.3:a:symfony:framework_bundle:v5.4.10:*:*:*:*:*:*:*",
			"cpe:2.3:a:framework-bundle:framework-bundle:v5.4.10:*:*:*:*:*:*:*",
		},
		PURL: model.PURL("pkg:composer/symfony/framework-bundle@v5.4.10"),
	}
)

func TestParseComposerPackages(t *testing.T) {
	composerPath := filepath.Join("..", "..", "docs", "references", "composer", "composer.lock")
	testLocation := model.Location{Path: composerPath}
	err := parseComposerPackages(&testLocation)
	if err != nil {
		t.Error("Test Failed: Error occurred while reading Composer content.")
	}
}

func TestParseComposerPURL(t *testing.T) {
	tests := []PhpPurlResult{
		{&phpPackage1, model.PURL("pkg:composer/league/tactician@v1.1.0")},
		{&phpPackage2, model.PURL("pkg:composer/phpdocumentor/reflection@5.2.0")},
		{&phpPackage3, model.PURL("pkg:composer/symfony/framework-bundle@v5.4.10")},
	}

	for _, test := range tests {
		parseComposerPURL(test._package)
		if test._package.PURL != test.expected {
			t.Errorf("Test Failed: Expected an output of %v, received: %v", test.expected, test._package.PURL)
		}
	}
}
