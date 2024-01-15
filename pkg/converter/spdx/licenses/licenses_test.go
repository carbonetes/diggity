package licenses

import (
	"encoding/json"
	"net/http"
	"testing"
)

const (
	licenseURL = "https://spdx.org/licenses/licenses.json"
)

type (
	LicenseListResult struct {
		input    string
		expected string
	}

	LicensesResult struct {
		Version string `json:"licenseListVersion"`
	}
)

func TestLicenseListVersion(t *testing.T) {
	var latestList LicensesResult
	res, err := http.Get(licenseURL)
	if err != nil {
		t.Fatalf("Error occurred when fetching licenses list: %+v", err)
	}

	if err = json.NewDecoder(res.Body).Decode(&latestList); err != nil {
		t.Fatalf("Error occurred when decoding license list: %+v", err)
	}

	if latestList.Version != ListVersion {
		t.Errorf("Test Failed: Version Mismatch. The licenses may not be up to date. Running automated update. \nFetched Latest Version: %+v, Current Implemented Version: %+v",
			latestList.Version, ListVersion)

		// Auto-Update License List
		updateSPDXLicenses()
		t.Error("Updated License List. Kindly check licenses.go and rerun unit test to validate. Apply necessary changes as needed.")
	}
}

func TestLiceseList(t *testing.T) {
	tests := []LicenseListResult{
		{"0bsd", "0BSD"},
		{"bsd-2.0-clause", "BSD-2-Clause"},
		{"gpl-3", "GPL-3.0-only"},
		{"mit", "MIT"},
		{"zlib", "Zlib"},
	}

	for _, test := range tests {
		if List[test.input] != test.expected {
			t.Errorf("Test Failed: Expected output of %v, received: %v", test.expected, List[test.input])
		}
	}
}