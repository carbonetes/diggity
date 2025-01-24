package nuget

import (
	"encoding/json"
	"encoding/xml"

	"github.com/carbonetes/diggity/internal/log"
)

// DotnetDeps - .NET Dependencies
type DotnetDeps struct {
	Libraries map[string]DotnetLibrary `json:"libraries"`
}

// DotnetLibrary - .NET libraries
type DotnetLibrary struct {
	Type     string `json:"type"`
	Path     string `json:"path"`
	Sha512   string `json:"sha512"`
	HashPath string `json:"hashPath"`
}

// ProjectFile represents the structure of a .csproj and .vdproj file.
type ProjectFile struct {
	XMLName       xml.Name        `xml:"Project"`
	PropertyGroup PropertyGroup `xml:"PropertyGroup"`
	ItemGroup     ItemGroup     `xml:"ItemGroup"`
}

// PropertyGroup represents a group of properties in the .csproj file
type PropertyGroup struct {
	TargetFramework          string `xml:"TargetFramework"`
	GenerateDocumentationFile bool   `xml:"GenerateDocumentationFile"`
	PackageId                string `xml:"PackageId"`
	PackageVersion           string `xml:"PackageVersion"`
	Version                  string `xml:"Version"`
	Authors                  string `xml:"Authors"`
	Description              string `xml:"Description"`
	PackageRequireLicenseAcceptance bool `xml:"PackageRequireLicenseAcceptance"`
	PackageReleaseNotes      string `xml:"PackageReleaseNotes"`
	Copyright                string `xml:"Copyright"`
	PackageTags              string `xml:"PackageTags"`
	IsPackable               bool   `xml:"IsPackable"`
	GeneratePackageOnBuild   bool   `xml:"GeneratePackageOnBuild"`
	ProjectUrl               string `xml:"ProjectUrl"`
	RepositoryUrl            string `xml:"RepositoryUrl"`
	PackageIcon              string `xml:"PackageIcon"`
	PackageLicenseExpression string `xml:"PackageLicenseExpression"`
}

// ItemGroup represents a group of items in the .csproj file
type ItemGroup struct {
	PackageReferences []PackageReference `xml:"PackageReference"`
	References        []Reference        `xml:"Reference"`
}

// Reference represents a reference in the .csproj file
type Reference struct {
	Include string `xml:"Include,attr"`
	Version      string `xml:"Version,attr"`
	Culture 	string `xml:"Culture,attr"`
	PublicKeyToken string `xml:"PublicKeyToken,attr"`
	ProcessorArchitecture string `xml:"processorArchitecture,attr"`
	HintPath string `xml:"HintPath,omitempty"`
}

// PackageReference represents a package reference in the .csproj file
type PackageReference struct {
	Include      string `xml:"Include,attr"`
	Version      string `xml:"Version,attr"`
	PrivateAssets string `xml:"PrivateAssets,attr,omitempty"`
	IncludeAssets string `xml:"IncludeAssets,attr,omitempty"`
}

func readManifestFile(content []byte) *DotnetDeps {
	var metadata DotnetDeps
	err := json.Unmarshal(content, &metadata)
	if err != nil {
		log.Debug("Failed to unmarshal project.assets.json")
		return nil
	}
	return &metadata
}

func parseProjectFile(content []byte) (*ProjectFile, error) {
	var proj ProjectFile
	if err := xml.Unmarshal(content, &proj); err != nil {
		return nil, err
	}

	return &proj, nil
}
