package spdxutils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	versionPackage "github.com/carbonetes/diggity/internal/version"
	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/model/metadata"
	"github.com/carbonetes/diggity/pkg/parser/os/apk"
	"github.com/carbonetes/diggity/pkg/parser/os/dpkg"
	"github.com/carbonetes/diggity/pkg/parser/language/gem"
	"github.com/carbonetes/diggity/pkg/parser/language/python"

	"github.com/google/uuid"
	spdxcommon "github.com/spdx/tools-golang/spdx/v2/common"
	spdx23 "github.com/spdx/tools-golang/spdx/v2/v2_3"
)

const (
	// Version : current implemented version (2.3)
	Version = "SPDX-2.3"
	// DataLicense : 6.2 Data license field Table 3 https://spdx.github.io/spdx-spec/v2.2.2/document-creation-information/
	DataLicense = "CC0-1.0"
	// Creator : Carbonetes
	Creator = "Carbonetes"
	// Ref : SPDX Ref Prefix
	Ref = "SPDXRef-"
	// Doc : Document Prefix
	Doc = "DOCUMENT"
	// NoAssertion : NO ASSERTION (For licenses)
	NoAssertion = "NOASSERTION"
	// None : NONE
	None = "NONE"

	// Extrnal Ref Types
	cpeType  = spdxcommon.TypeSecurityCPE23Type
	purlType = spdxcommon.TypePackageManagerPURL

	organization     = "Organization"
	tool             = "Tool"
	person           = "Person"
	security         = "SECURITY"
	packageManager   = "PACKAGE_MANAGER"
	licenseSeparator = " AND "
	parsedFrom       = "Information parsed from"
	namespace        = "https://console.carbonetes.com/diggity/image/"
	url              = "https://spdx.org/licenses/licenses.json"
)

// CreateInfo : Contains creator and tool information.
var (
	Tool       = "Tool: " + versionPackage.FromBuild().Version
	CreateInfo = []spdxcommon.Creator{
		{
			Creator:     Creator,
			CreatorType: organization,
		},
		{
			Creator:     versionPackage.FromBuild().Version,
			CreatorType: tool,
		},
	}
)

// ExternalRefs helper
func ExternalRefs(p *model.Package) (refs []*spdx23.PackageExternalReference) {
	// Init CPEs
	for _, cpe := range p.CPEs {
		var cpeRef spdx23.PackageExternalReference
		cpeRef.Category = security
		cpeRef.Locator = cpe
		cpeRef.RefType = cpeType
		refs = append(refs, &cpeRef)
	}

	// Init PURL
	var purlRef spdx23.PackageExternalReference
	purlRef.Category = packageManager
	purlRef.Locator = string(p.PURL)
	purlRef.RefType = purlType
	refs = append(refs, &purlRef)

	return refs
}

// LicensesDeclared helper
func LicensesDeclared(p *model.Package) string {
	// Check if package has licenses
	if len(p.Licenses) <= 0 {
		return None
	}

	var licenses []string

	// Validate Licenses from License List
	for _, license := range p.Licenses {
		if CheckLicense(license) != "" {
			licenses = append(licenses, license)
		}
	}

	if len(licenses) > 0 {
		return strings.Join(licenses, licenseSeparator)
	}

	return NoAssertion

}

// Homepage helper
func Homepage(p *model.Package) string {
	switch m := p.Metadata.(type) {
	case metadata.PackageJSON:
		return m.Homepage
	case gem.Metadata:
		if val, ok := m["homepage"]; ok {
			return val.(string)
		}
	}
	return ""
}

// SourceInfo helper
func SourceInfo(p *model.Package) string {
	var source string
	var locations []string

	switch p.Type {
	case "apk":
		source = "APK DB"
	case "php":
		source = "PHP composer manifest"
	case "pub":
		source = "pubspec manifest"
	case "deb":
		source = "DPKG DB"
	case "gem":
		source = "gem metadata"
	case "go-module":
		source = "go-module information"
	case "java":
		source = "java archive"
	case "npm":
		source = "node module manifest"
	case "dotnet":
		source = "dotnet project assets"
	case "python":
		source = "python package manifest"
	case "rpm":
		source = "RPM DB"
	case "rust-crate":
		source = "rust cargo manifest"
	case "conan":
		source = "conan manifest"
	case "hackage":
		source = "stack or cabal manifest"
	case "pod":
		source = "cocoapods manifest"
	case "hex":
		source = "mix o rebar3 manifest"
	case "portage":
		source = "Portage DB"
	default:
		source = ""
	}

	if len(p.Locations) > 0 {
		for _, loc := range p.Locations {
			locations = append(locations, FormatPath(loc.Path))
		}
	}

	return fmt.Sprintf("%s %s: %s", parsedFrom, source, strings.Join(locations, ", "))
}

// DownloadLocation helper
func DownloadLocation(p *model.Package) string {
	var url string

	switch m := p.Metadata.(type) {
	case apk.Metadata:
		if val, ok := m["PackageURL"]; ok {
			url = val.(string)
		}
	case metadata.PackageJSON:
		switch m.Repository.(type) {
		case map[string]interface{}:
			repo := m.Repository.(map[string]interface{})
			if _, ok := repo["url"]; ok {
				url = repo["url"].(string)
			}
		case string:
			url = m.Repository.(string)
		}
	default:
		return NoAssertion
	}

	if strings.TrimSpace(url) == "" {
		return None
	}

	return url
}

// Originator helper
func Originator(p *model.Package) (string, string) {
	var originator string

	switch m := p.Metadata.(type) {
	// Cases with existing metadata models
	case metadata.RPMMetadata:
		return organization, m.Vendor
	case metadata.PackageJSON:
		switch m.Author.(type) {
		case map[string]interface{}:
			author := m.Author.(map[string]interface{})
			authorDetails := []string{}

			if val, ok := author["name"]; ok {
				authorDetails = append(authorDetails, val.(string))
			}
			if val, ok := author["email"]; ok {
				authorDetails = append(authorDetails, val.(string))
			}
			originator = strings.Join(authorDetails, " ")
		case string:
			author := m.Author.(string)
			originator = FormatAuthor(author)
		}

	// Cases with metadata declared within the parser
	case apk.Metadata:
		if val, ok := m["Maintainer"]; ok {
			originator = val.(string)
		}
	case python.Metadata:
		if val, ok := m["Author"]; ok {
			originator = val.(string)
		}
	case dpkg.Metadata:
		if val, ok := m["Maintainer"]; ok {
			originator = val.(string)
		}
	case gem.Metadata:
		if val, ok := m["authors"]; ok {
			originator = val.([]string)[0]
		}
	}

	if originator != "" {
		return person, originator
	}

	return "", ""
}

// FormatName helper
func FormatName(args *model.Arguments) string {
	// Check if tar or dir was scanned
	if args.Image == nil {
		if *args.Tar != "" {
			tarFile := filepath.Base(*args.Tar)
			return strings.Split(tarFile, ".")[0]
		}
		if *args.Dir != "" {
			return filepath.Base(*args.Dir)
		}
	}
	if strings.Contains(*args.Image, ":") {
		return strings.Split(*args.Image, ":")[0]
	}
	return *args.Image
}

// FormatNamespace helper
func FormatNamespace(imageName string) string {
	return namespace + imageName + "-" + uuid.NewString()
}

// FormatPath helper
func FormatPath(path string) string {
	pathSlice := strings.Split(path, string(os.PathSeparator))
	return strings.Join(pathSlice, "/")
}

// FormatTagID helper
func FormatTagID(p *model.Package) string {
	return fmt.Sprintf("%sPackage-%+v-%s-%s", Ref, p.Type, p.Name, p.ID)
}

// CheckLicense helper
func CheckLicense(id string) string {
	license := LicenseList[strings.ToLower(id)]
	return license
}

// FormatAuthor helper
func FormatAuthor(authorString string) string {
	author := []string{}

	// Check for empty author
	if strings.TrimSpace(authorString) == "" {
		return ""
	}

	authorDetails := strings.Split(authorString, " ")
	if len(authorDetails) == 1 {
		return authorDetails[0]
	}

	for _, detail := range authorDetails {
		if strings.Contains(detail, "http") && strings.Contains(detail, ".") && strings.Contains(detail, "/") {
			continue
		}
		author = append(author, detail)
	}

	return strings.Join(author, " ")
}
