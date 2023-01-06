package parser

import (
	"path/filepath"
	"testing"

	"github.com/carbonetes/diggity/internal/model"
)

type (
	InitPythonPackages struct {
		_package *model.Package
		metadata PythonMetadata
		location model.Location
		expected *model.Package
	}
	PythonPurlResult struct {
		_package *model.Package
		expected model.PURL
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
		PURL: model.PURL("pkg:python/wsgiref@0.1.2"),
		Metadata: PythonMetadata{
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
		PURL: model.PURL("pkg:python/Python@2.7.16"),
		Metadata: PythonMetadata{
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
	pythonMetadata1 = PythonMetadata{
		"Author":           "Phillip J. Eby",
		"Author-email":     "web-sig@python.org",
		"License":          "PSF or ZPL",
		"Metadata-Version": "1.0",
		"Name":             "wsgiref",
		"Platform":         "UNKNOWN",
		"Summary":          "WSGI (PEP 333) Reference Library",
		"Version":          "0.1.2",
	}
	pythonMetadata2 = PythonMetadata{
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
	pythonEggPath := filepath.Join("..", "..", "docs", "references", "python", "argparse.egg-info")
	pythonPackagepath := filepath.Join("..", "..", "docs", "references", "python", "METADATA")
	pythonEggLocation := model.Location{Path: pythonEggPath}
	pythonPackagelocation := model.Location{Path: pythonPackagepath}
	errPythonEgg := readPythonContent(&pythonEggLocation)
	errPythonPackage := readPythonContent(&pythonPackagelocation)
	if errPythonEgg != nil {
		t.Errorf("Test Failed: Error occurred while reading Python egg-info content. %v", errPythonEgg)
	}
	if errPythonPackage != nil {
		t.Errorf("Test Failed: Error occurred while reading Python package content. %v", errPythonPackage)
	}
}

func TestInitPythonPackages(t *testing.T) {
	var _package1, _package2 model.Package
	tests := []InitPythonPackages{
		{&_package1, pythonMetadata1, pythonLocation1, &pythonPackage1},
		{&_package2, pythonMetadata2, pythonLocation2, &pythonPackage2},
	}
	for _, test := range tests {
		output := initPythonPackages(test._package, test.metadata, &test.location)
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
		{&pythonPackage1, model.PURL("pkg:python/wsgiref@0.1.2")},
		{&pythonPackage2, model.PURL("pkg:python/Python@2.7.16")},
	}
	for _, test := range tests {
		parsePythonPackageURL(test._package)
		if test._package.PURL != test.expected {
			t.Errorf("Test Failed: Expected an output of %v, received: %v", test.expected, test._package.PURL)
		}
	}
}
