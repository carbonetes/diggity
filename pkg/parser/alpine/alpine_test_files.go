package alpine

import (
	"testing"

	"github.com/carbonetes/diggity/pkg/model"
)

type (
	ApkPackageResult struct {
		pkg      *model.Package
		expected *model.Package
	}
	AlpineFilesResult struct {
		content  string
		expected []model.File
	}
)

func TestGetAlpineFiles(t *testing.T) {
	var tests = []AlpineFilesResult{
		{apkContent1, apkFiles},
		{apkContent2, nil},
	}
	for _, test := range tests {
		files := getAlpineFiles(test.content)
		if len(files) != len(test.expected) {
			t.Errorf("Test Failed: Expected Packages of length %+v, Received: %+v.", len(test.expected), len(files))
		}
		CheckAlpineFiles(t, test, files)
	}
}

func CheckAlpineFiles(t *testing.T, test AlpineFilesResult, files []model.File) {
	for i, file := range files {
		if file != test.expected[i] {
			t.Errorf("Test Failed: Expected file result of %v, received: %v", test.expected[i], file)
			continue
		}

		if file.Digest == nil && test.expected[i].Digest != nil {
			t.Errorf("Test Failed: Expected Digest not nil.")
			continue
		}

		if file.Digest.(map[string]string)["algorithm"] != test.expected[i].Digest.(map[string]string)["algorithm"] {
			t.Errorf("Test Failed: Expected digest algorithm of %v, received: %v", test.expected[i], file)
			continue
		}

		if file.Digest.(map[string]string)["value"] != test.expected[i].Digest.(map[string]string)["value"] {
			t.Errorf("Test Failed: Expected digest value of %v, received: %v", test.expected[i], file)
			continue
		}
	}
}
