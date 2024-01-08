package spdx

import (
	"fmt"
	"strings"

	"github.com/carbonetes/diggity/internal/version"
	"github.com/carbonetes/diggity/pkg/types"
	rpmdb "github.com/knqyf263/go-rpmdb/pkg"
	spdxcommon "github.com/spdx/tools-golang/spdx/v2/common"
	spdx23 "github.com/spdx/tools-golang/spdx/v2/v2_3"
)

const (
	// Version : current implemented version (2.3)
	Version = "SPDX-2.3"
	// DataLicense : 6.2 Data license field Table 3 https://spdx.github.io/spdx-spec/v2.3/document-creation-information/
	DataLicense = "CC0-1.0"
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

var (
	Creators = []spdxcommon.Creator{
		{
			Creator:     "carbonetes",
			CreatorType: "Organization",
		},
		{
			Creator:     "diggity-" + version.FromBuild().Version,
			CreatorType: "Tool",
		},
	}
)

// Add references based on component purl and cpes
func ExternalRefs(c types.Component) (refs []*spdx23.PackageExternalReference) {
	// Init CPEs
	for _, cpe := range c.CPEs {
		var cpeRef spdx23.PackageExternalReference
		cpeRef.Category = security
		cpeRef.Locator = cpe
		cpeRef.RefType = cpeType
		refs = append(refs, &cpeRef)
	}

	// Init PURL
	var purlRef spdx23.PackageExternalReference
	purlRef.Category = packageManager
	purlRef.Locator = string(c.PURL)
	purlRef.RefType = purlType
	refs = append(refs, &purlRef)

	return refs
}

// Check if license is in SPDX License List
func LicensesDeclared(c types.Component) string {
	// Check if package has licenses
	if len(c.Licenses) == 0 {
		return None
	}

	var licenses []string

	// Validate Licenses from License List
	for _, license := range c.Licenses {
		if CheckLicense(license) != "" {
			licenses = append(licenses, license)
		}
	}

	if len(licenses) > 0 {
		return strings.Join(licenses, licenseSeparator)
	}

	return NoAssertion

}

// Get the component's homepage from the metadata if it exists
func Homepage(c types.Component) string {
	switch m := c.Metadata.(type) {
	case map[string]interface{}:
		if val, ok := m["homepage"]; ok {
			return val.(string)
		}
	}
	return ""
}

// Determine where the component was parsed from and return the source information string
func SourceInfo(c types.Component) string {
	var source string

	switch c.Type {
	case "apk":
		source = "APK DB"
	case "composer":
		source = "PHP composer manifest"
	case "pub":
		source = "pubspec manifest"
	case "dpkg":
		source = "DPKG DB"
	case "gem":
		source = "gem metadata"
	case "golang":
		source = "go-module information"
	case "java":
		source = "java archive"
	case "npm":
		source = "node module manifest"
	case "nuget":
		source = "dotnet project assets"
	case "pypi":
		source = "python package manifest"
	case "rpm":
		source = "RPM DB"
	case "cargo":
		source = "rust cargo manifest"
	case "conan":
		source = "conan manifest"
	case "hackage":
		source = "stack or cabal manifest"
	case "cocoapods":
		source = "cocoapods manifest"
	case "hex":
		source = "mix o rebar3 manifest"
	case "portage":
		source = "Portage DB"
	case "cran":
		source = "DESCRIPTION file"
	default:
		source = ""
	}

	return fmt.Sprintf("%s %s: %s", parsedFrom, source, c.Origin)
}

// Get the component's download location from the metadata if it exists
func DownloadLocation(c types.Component) string {
	var url string

	switch m := c.Metadata.(type) {
	case map[string]interface{}:
		if val, ok := m["PackageURL"]; ok {
			url = val.(string)
		}

		if _, ok := m["repository"]; ok {
			repo := m["repository"].(map[string]interface{})
			if _, ok := repo["url"]; ok {
				url = repo["url"].(string)
			}
		}
	default:
		return NoAssertion
	}

	if strings.TrimSpace(url) == "" {
		return None
	}

	return url
}

// Get the component's author or maintainer from the metadata if it exists
func Originator(p types.Component) (string, string) {
	var originator string

	switch m := p.Metadata.(type) {
	// Cases with existing metadata models
	case rpmdb.PackageInfo:
		return organization, m.Vendor
	case map[string]interface{}:
		if val, ok := m["Maintainer"]; ok {
			originator = val.(string)
		}

		if val, ok := m["authors"]; ok {
			originator = val.([]string)[0]
		}

		if val, ok := m["Author"]; ok {
			originator = val.(string)
		}

		if _, ok := m["author"]; ok {
			switch m["author"].(type) {
			case map[string]interface{}:
				author := m["author"].(map[string]interface{})
				authorDetails := []string{}

				if val, ok := author["name"]; ok {
					authorDetails = append(authorDetails, val.(string))
				}
				if val, ok := author["email"]; ok {
					authorDetails = append(authorDetails, val.(string))
				}
				originator = strings.Join(authorDetails, " ")
			case string:
				author := m["author"].(string)
				originator = FormatAuthor(author)
			}
		}
	}

	if originator == "" {
		return person, originator
	}

	return "", ""
}
