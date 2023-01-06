package spdxutils

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"

	log "github.com/carbonetes/diggity/internal/logger"
)

const (
	// LicenseURL is the reference for the SPDX License List
	LicenseURL = "https://spdx.org/licenses/licenses.json"
	filename   = "licenses.go"
)

type (
	spdxLicenseList struct {
		Version     string        `json:"licenseListVersion"`
		LicenseList []spdxLicense `json:"licenses"`
	}

	spdxLicense struct {
		ID            string   `json:"licenseId"`
		Name          string   `json:"name"`
		IsDeprecated  bool     `json:"isDeprecatedLicenseId"`
		IsOsiApproved bool     `json:"isOsiApproved"`
		SeeAlso       []string `json:"seeAlso"`
	}

	licenseTemplate struct {
		LastUpdate string
		Version    string
		Licenses   map[string]string
	}
)

// Update SPDX License List File
func updateSPDXLicenses() {
	latestList := fetchLatestList()

	if LicenseListVersion == latestList.Version {
		return
	}

	file, err := os.Create(filename)
	if err != nil {
		log.GetLogger().Printf("Error occured when creating file: %+v", err)
	}

	licenses := updateLicenseList(latestList.LicenseList)

	fileTemplate := initFileTemplate()

	if err := fileTemplate.Execute(file,
		licenseTemplate{
			LastUpdate: time.Now().String(),
			Version:    latestList.Version,
			Licenses:   licenses,
		}); err != nil {
		log.GetLogger().Printf("Error occured upon writing file: %+v", err)
	}
	log.GetLogger().Print("Successfully updated SPDX License List File.")
}

// Fetch Latest SPDX License List
func fetchLatestList() spdxLicenseList {
	var licenseList spdxLicenseList
	res, err := http.Get(LicenseURL)

	if err != nil {
		log.GetLogger().Printf("Error occured when fetching license list: %+v", err)
	}

	if err = json.NewDecoder(res.Body).Decode(&licenseList); err != nil {
		log.GetLogger().Printf("Error occured when decoding license list: %+v", err)
	}
	defer res.Body.Close()

	return licenseList
}

// Add new licenses to existing list
func updateLicenseList(licenseList []spdxLicense) map[string]string {
	licenses := LicenseList

	for _, license := range licenseList {
		licenseKey := strings.ToLower(license.ID)
		if _, exists := licenses[licenseKey]; !exists {
			licenses[strings.ToLower(license.ID)] = license.ID
			log.GetLogger().Printf("Adding license: %+v", license.ID)
		}
	}

	return licenses
}

// Init spdx license go file template
func initFileTemplate() *template.Template {
	return template.Must(template.New("").Parse(`package spdxutils
	
	// Source URL : https://spdx.org/licenses/licenses.json
	// Helpers for validating SPDX licenses.
	// Needs to be updated regularly, possibly with automation.
	// Last Updated: {{ .LastUpdate }}

	// LicenseListVersion is the current implemented version for SPDX.
	const LicenseListVersion = "{{ .Version }}"

	// LicenseList contains the referece licenses from the source URL.
	var LicenseList = map[string]string{
	{{- range $key, $value := .Licenses }}
		"{{ $key }}": "{{ $value }}",
	{{- end }}
	}
	`))
}
