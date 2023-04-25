package python

import (
	"path/filepath"
	"testing"

	"github.com/carbonetes/diggity/pkg/model"
)

type (
	InitPythonPackages struct {
		pkg      *model.Package
		metadata Metadata
		location model.Location
		expected *model.Package
	}
	PythonPurlResult struct {
		pkg      *model.Package
		expected model.PURL
	}

	ParseRequirementsResult struct {
		input   string
		name    string
		version string
	}

	PoetryFileMetadataResult struct {
		input    string
		expected map[string]string
	}
)

var (
	pythonPackage1 = model.Package{
		Name:    "wsgiref",
		Type:    "python",
		Version: "0.1.2",
		Path:    "wsgiref",
		Locations: []model.Location{
			{
				Path:      filepath.Join("usr", "lib", "python2.7", "wsgiref.egg-info"),
				LayerHash: "ebd02e2b3b566df7b16480f9d0796a720cb87ecf4dbe522f5ca01f649eddfe64",
			},
		},

		Description: "",
		Licenses: []string{
			"PSF or ZPL",
		},
		CPEs: []string{
			"cpe:2.3:a:wsgiref:wsgiref:0.1.2:*:*:*:*:*:*:*",
		},
		PURL: model.PURL("pkg:pypi/wsgiref@0.1.2"),
		Metadata: Metadata{
			"Author":           "Phillip J. Eby",
			"Author-email":     "web-sig@python.org",
			"License":          "PSF or ZPL",
			"Metadata-Version": "1.0",
			"Name":             "wsgiref",
			"Platform":         "UNKNOWN",
			"Summary":          "WSGI (PEP 333) Reference Library",
			"Version":          "0.1.2",
		},
	}
	pythonPackage2 = model.Package{
		Name:    "Python",
		Type:    "python",
		Version: "2.7.16",
		Path:    "Python",
		Locations: []model.Location{
			{
				Path:      filepath.Join("usr", "lib", "python2.7", "lib-dynload", "Python-2.7.egg-info"),
				LayerHash: "ebd02e2b3b566df7b16480f9d0796a720cb87ecf4dbe522f5ca01f649eddfe64",
			},
		},
		Description: "",
		Licenses: []string{
			"PSF license",
		},
		CPEs: []string{
			"cpe:2.3:a:Python:Python:2.7.16:*:*:*:*:*:*:*",
		},
		PURL: model.PURL("pkg:pypi/Python@2.7.16"),
		Metadata: Metadata{
			"Author":           "Guido van Rossum and the Python community",
			"Author-email":     "python-dev@python.org",
			"Classifier":       "Topic :: Software Development",
			"Description":      "Python is an interpreted, interactive, object-oriented programming        language. It is often compared to Tcl, Perl, Scheme or Java.        Python combines remarkable power with very clear syntax. It has        modules, classes, exceptions, very high level dynamic data types, and        dynamic typing. There are interfaces to many system calls and        libraries, as well as to various windowing systems (X11, Motif, Tk,        Mac, MFC). New built-in modules are easily written in C or C++. Python        is also usable as an extension language for applications that need a        programmable interface.",
			"Home-page":        "http://www.python.org/2.7",
			"License":          "PSF license",
			"Metadata-Version": "1.1",
			"Name":             "Python",
			"Platform":         "Many",
			"Summary":          "A high-level object-oriented programming language",
			"Version":          "2.7.16",
		},
	}
	Metadata1 = Metadata{
		"Author":           "Phillip J. Eby",
		"Author-email":     "web-sig@python.org",
		"License":          "PSF or ZPL",
		"Metadata-Version": "1.0",
		"Name":             "wsgiref",
		"Platform":         "UNKNOWN",
		"Summary":          "WSGI (PEP 333) Reference Library",
		"Version":          "0.1.2",
	}
	Metadata2 = Metadata{
		"Author":           "Guido van Rossum and the Python community",
		"Author-email":     "python-dev@python.org",
		"Classifier":       "Topic :: Software Development",
		"Description":      "Python is an interpreted, interactive, object-oriented programming        language. It is often compared to Tcl, Perl, Scheme or Java.        Python combines remarkable power with very clear syntax. It has        modules, classes, exceptions, very high level dynamic data types, and        dynamic typing. There are interfaces to many system calls and        libraries, as well as to various windowing systems (X11, Motif, Tk,        Mac, MFC). New built-in modules are easily written in C or C++. Python        is also usable as an extension language for applications that need a        programmable interface.",
		"Home-page":        "http://www.python.org/2.7",
		"License":          "PSF license",
		"Metadata-Version": "1.1",
		"Name":             "Python",
		"Platform":         "Many",
		"Summary":          "A high-level object-oriented programming language",
		"Version":          "2.7.16",
	}
	pythonLocation1 = model.Location{
		Path:      filepath.Join("AppData", "Local", "Temp", "1944417693", "diggity-tmp-2c4e711d-1330-44b1-978a-0ca69359e51d", "ebd02e2b3b566df7b16480f9d0796a720cb87ecf4dbe522f5ca01f649eddfe64", "usr", "lib", "python2.7", "wsgiref.egg-info"),
		LayerHash: "f1a5f5ce6b163fac7f09b47645c56d2ab676bdcdb268eef06a4d9b782a75bfd0",
	}
	pythonLocation2 = model.Location{
		Path:      filepath.Join("AppData", "Local", "Temp", "1944417693", "diggity-tmp-2c4e711d-1330-44b1-978a-0ca69359e51d", "ebd02e2b3b566df7b16480f9d0796a720cb87ecf4dbe522f5ca01f649eddfe64", "usr", "lib", "python2.7", "lib-dynload", "Python-2.7.egg-info"),
		LayerHash: "f1a5f5ce6b163fac7f09b47645c56d2ab676bdcdb268eef06a4d9b782a75bfd0",
	}
)

func TestReadPythonContent(t *testing.T) {
	pythonEggPath := filepath.Join("..", "..", "..", "docs", "references", "python", "argparse.egg-info")
	pythonPackagepath := filepath.Join("..", "..", "..", "docs", "references", "python", "METADATA")
	pythonEggLocation := model.Location{Path: pythonEggPath}
	pythonPackagelocation := model.Location{Path: pythonPackagepath}
	pkgs := new([]model.Package)
	errPythonEgg := readPythonContent(&pythonEggLocation, pkgs)
	errPythonPackage := readPythonContent(&pythonPackagelocation, pkgs)
	if errPythonEgg != nil {
		t.Errorf("Test Failed: Error occurred while reading Python egg-info content. %v", errPythonEgg)
	}
	if errPythonPackage != nil {
		t.Errorf("Test Failed: Error occurred while reading Python package content. %v", errPythonPackage)
	}
}

func TestReadRequirementsContent(t *testing.T) {
	requirementsPath := filepath.Join("..", "..", "..", "docs", "references", "python", "requirements.txt")
	requirementsLocation := model.Location{Path: requirementsPath}
	pkgs := new([]model.Package)
	err := readRequirementsContent(&requirementsLocation, pkgs)
	if err != nil {
		t.Errorf("Test Failed: Error occurred while reading requirements.txt content. %v", err)
	}
}

func TestInitPythonPackages(t *testing.T) {
	var pkg1, pkg2 model.Package
	tests := []InitPythonPackages{
		{&pkg1, Metadata1, pythonLocation1, &pythonPackage1},
		{&pkg2, Metadata2, pythonLocation2, &pythonPackage2},
	}
	for _, test := range tests {
		output := initPythonPackages(test.metadata, &test.location)
		if output.Name != test.expected.Name ||
			output.Version != test.expected.Version ||
			output.Description != test.expected.Description ||
			string(output.PURL) != string(test.expected.PURL) {

			t.Errorf("Test Failed:\n Expected output of %v \n, Received: %v \n", test.expected, output)
		}
		for i := range output.CPEs {
			if output.CPEs[i] != test.expected.CPEs[i] {
				t.Errorf("Test Failed:\n Expected output of %v \n, Received: %v \n", test.expected.CPEs[i], output.CPEs[i])
			}
		}
	}
}

func TestParsePythonPackageURL(t *testing.T) {
	tests := []PythonPurlResult{
		{&pythonPackage1, model.PURL("pkg:pypi/wsgiref@0.1.2")},
		{&pythonPackage2, model.PURL("pkg:pypi/Python@2.7.16")},
	}
	for _, test := range tests {
		parsePythonPackageURL(test.pkg)
		if test.pkg.PURL != test.expected {
			t.Errorf("Test Failed: Expected an output of %v, received: %v", test.expected, test.pkg.PURL)
		}
	}
}

func TestParseRequirements(t *testing.T) {
	tests := []ParseRequirementsResult{
		{"test==1.0.0", "test", "1.0.0"},
		{"test== 1.0.0", "test", "1.0.0"},
		{"test ==1.0.0", "test", "1.0.0"},
		{"test == 1.0.0", "test", "1.0.0"},
		{"test==1.0.0 #Comment", "test", "1.0.0"},
		{"Django==1.11.29", "Django", "1.11.29"},
		{" flask == 4.0.0", "flask", "4.0.0"},
	}
	for _, test := range tests {
		if name, version := parseRequirements(test.input); name != test.name || version != test.version {
			t.Errorf("Test Failed: Expected an output of %+v and %+v, received: %v and %+v", test.name, test.version, name, version)
		}
	}
}

func TestPoetryFileMetadata(t *testing.T) {
	tests := []PoetryFileMetadataResult{
		{`{file = "attrs-21.4.0-py2.py3-none-any.whl", hash = "sha256:2d27e3784d7a565d36ab851fe94887c5eccd6a463168875832a1be79c82828b4"},`,
			map[string]string{
				"file": "attrs-21.4.0-py2.py3-none-any.whl",
				"hash": "sha256:2d27e3784d7a565d36ab851fe94887c5eccd6a463168875832a1be79c82828b4",
			}},
		{`{file = "PyYAML-6.0-cp310-cp310-macosx_10_9_x86_64.whl", hash = "sha256:d4db7c7aef085872ef65a8fd7d6d09a14ae91f691dec3e87ee5ee0539d516f53"},`,
			map[string]string{
				"file": "PyYAML-6.0-cp310-cp310-macosx_10_9_x86_64.whl",
				"hash": "sha256:d4db7c7aef085872ef65a8fd7d6d09a14ae91f691dec3e87ee5ee0539d516f53",
			}},
		{`{file = "urllib3-1.26.9.tar.gz", hash = "sha256:aabaf16477806a5e1dd19aa41f8c2b7950dd3c746362d7e3223dbe6de6ac448e"},`,
			map[string]string{
				"file": "urllib3-1.26.9.tar.gz",
				"hash": "sha256:aabaf16477806a5e1dd19aa41f8c2b7950dd3c746362d7e3223dbe6de6ac448e",
			}},
	}
	for _, test := range tests {
		output := poetryFileMetadata(test.input)
		if output["file"] != test.expected["file"] || output["hash"] != test.expected["hash"] {
			t.Errorf("Test Failed: Expected an output of %v, received: %v", test.expected, output)
		}
	}
}
