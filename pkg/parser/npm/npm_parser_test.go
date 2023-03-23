package npm

import (
	"path/filepath"
	"testing"

	"github.com/carbonetes/diggity/pkg/model"
)

type (
	NpmPurlResult struct {
		pkg      *model.Package
		expected model.PURL
	}
)

var (
	NpmPackagePURL1 = model.Package{
		Name:    "shebang-command",
		Type:    "npm",
		Version: "1.2.0",
		Path:    "shebang-command",
		Locations: []model.Location{
			{
				Path:      filepath.Join("usr", "local", "lib", "node_modules", "npm", "node_modules", "minipass-sized", "package-lock.json"),
				LayerHash: "d2d038f3cbb76afe98bb889cd54e1dd33e0119ee8feb3edd5d48968fa7f1fc50",
			},
		},
		Description: "",
		Licenses: []string{
			"",
		},
		CPEs: []string{
			"cpe:2.3:a:shebang-command:shebang-command:1.2.0:*:*:*:*:*:*:*",
			"cpe:2.3:a:shebang-command:shebang_command:1.2.0:*:*:*:*:*:*:*",
			"cpe:2.3:a:shebang_command:shebang_command:1.2.0:*:*:*:*:*:*:*",
			"cpe:2.3:a:shebang_command:shebang-command:1.2.0:*:*:*:*:*:*:*",
		},
		PURL: model.PURL("pkg:npm/shebang-command@1.2.0"),
		Metadata: LockMetadata{
			"version":   "1.2.0",
			"resolved":  "https://registry.npmjs.org/shebang-command/-/shebang-command-1.2.0.tgz",
			"integrity": "sha1-RKrGW2lbAzmJaMOfNj/uXer98eo=",
			"Requires": []string{
				"shebang-regex: ^1.0.0",
			},
		},
	}
	NpmPackagePURL2 = model.Package{
		Name:    "make-error",
		Type:    "npm",
		Version: "1.3.5",
		Path:    "make-error",
		Locations: []model.Location{
			{
				Path:      filepath.Join("usr", "local", "lib", "node_modules", "npm", "node_modules", "minipass-sized", "package-lock.json"),
				LayerHash: "d2d038f3cbb76afe98bb889cd54e1dd33e0119ee8feb3edd5d48968fa7f1fc50",
			},
		},
		Description: "",
		Licenses: []string{
			"",
		},
		CPEs: []string{
			"cpe:2.3:a:make-error:make-error:1.3.5:*:*:*:*:*:*:*",
			"cpe:2.3:a:make-error:make_error:1.3.5:*:*:*:*:*:*:*",
			"cpe:2.3:a:make_error:make_error:1.3.5:*:*:*:*:*:*:*",
			"cpe:2.3:a:make_error:make-error:1.3.5:*:*:*:*:*:*:*",
		},

		PURL: model.PURL("pkg:npm/make-error@1.3.5"),
		Metadata: LockMetadata{
			"version":   "1.3.5",
			"resolved":  "https://registry.npmjs.org/make-error/-/make-error-1.3.5.tgz",
			"integrity": "sha512-c3sIjNUow0+8swNwVpqoH4YCShKNFkMaw6oH1mNS2haDZQqkeZFlHS3dhoeEbKKmJB4vXpJucU6oH75aDYeE9g==",
			"Requires":  "",
		},
	}
	NpmPackagePURL3 = model.Package{
		Name:    "lcov-parse",
		Type:    "npm",
		Version: "1.0.0",
		Path:    "lcov-parse",
		Locations: []model.Location{
			{
				Path:      filepath.Join("usr", "local", "lib", "node_modules", "npm", "node_modules", "npm-normalize-package-bin", "package-lock.json"),
				LayerHash: "d2d038f3cbb76afe98bb889cd54e1dd33e0119ee8feb3edd5d48968fa7f1fc50",
			},
		},
		Description: "",
		Licenses: []string{
			"",
		},
		CPEs: []string{
			"cpe:2.3:a:lcov-parse:lcov-parse:1.0.0:*:*:*:*:*:*:*",
			"cpe:2.3:a:lcov-parse:lcov_parse:1.0.0:*:*:*:*:*:*:*",
			"cpe:2.3:a:lcov_parse:lcov_parse:1.0.0:*:*:*:*:*:*:*",
			"cpe:2.3:a:lcov_parse:lcov-parse:1.0.0:*:*:*:*:*:*:*",
		},
		PURL: model.PURL("pkg:npm/lcov-parse@1.0.0"),
		Metadata: LockMetadata{
			"version":   "1.0.0",
			"resolved":  "https://registry.npmjs.org/lcov-parse/-/lcov-parse-1.0.0.tgz",
			"integrity": "sha1-6w1GtUER68VhrLTECO+TY73I9+A=",
			"Requires":  "",
		},
	}
)

func TestReadNpmContent(t *testing.T) {
	packageLockPath := filepath.Join("..", "..", "..", "docs", "references", "npm", "package.json")
	testLocation := model.Location{Path: packageLockPath}
	err := readNpmContent(&testLocation)
	if err != nil {
		t.Errorf("Test Failed: Error occurred while reading Npm content. %v", err)
	}

}

func TestReadNpmLockContent(t *testing.T) {
	packagePath := filepath.Join("..", "..", "..", "docs", "references", "npm", "package-lock.json")
	testLocation := model.Location{Path: packagePath}
	err := readNpmLockContent(&testLocation)
	if err != nil {
		t.Errorf("Test Failed: Error occurred while reading Npm Lock content. %v", err)
	}
}

func TestReadYarnLockContent(t *testing.T) {
	yarnPath := filepath.Join("..", "..", "..", "docs", "references", "npm", "yarn.lock")
	testLocation := model.Location{Path: yarnPath}
	err := readYarnLockContent(&testLocation)
	if err != nil {
		t.Errorf("Test Failed: Error occurred while reading Yarn Lock content. %v", err)
	}
}

func TestParseNpmPackageURL(t *testing.T) {
	tests := []NpmPurlResult{
		{&NpmPackagePURL1, model.PURL("pkg:npm/shebang-command@1.2.0")},
		{&NpmPackagePURL2, model.PURL("pkg:npm/make-error@1.3.5")},
		{&NpmPackagePURL3, model.PURL("pkg:npm/lcov-parse@1.0.0")},
	}
	for _, test := range tests {
		parseNpmPackageURL(test.pkg)
		if test.pkg.PURL != test.expected {
			t.Errorf("Test Failed: Expected an output of %v, received: %v", test.expected, test.pkg.PURL)
		}
	}
}
