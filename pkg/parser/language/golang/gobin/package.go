package gobin

import (
	"fmt"
	"regexp"
	"runtime/debug"
	"strings"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/parser/language/golang"
	"github.com/carbonetes/diggity/pkg/parser/util"
	"github.com/google/uuid"
)

// Initialize Go Bin package contents
func newPackage(location *model.Location, buildInfo *debug.BuildInfo, module *debug.Module) *model.Package {
	var pkg model.Package
	pkg.ID = uuid.NewString()
	pkg.Type = golang.Type
	pkg.PackageOrigin = model.ApplicationPackage
	pkg.Language = golang.Language
	pkg.Parser = parser
	pkg.Licenses = []string{}
	pkg.Name = module.Path
	if module.Version == "(devel)" {
		pkg.Version = getVersion(buildInfo.Settings)
	} else {
		pkg.Version = module.Version
	}

	// get and format CPEs
	golang.GenerateCpes(&pkg, golang.SplitPath(pkg.Name))

	// pkg.CPEs = golang.FormatDevelCPEs(&pkg)

	// get purl
	golang.SetPurl(&pkg)

	// parse and fill final metadata
	metadata := parseMetadata(buildInfo, module)
	pkg.Metadata = metadata
	pkg.Path = util.TrimUntilLayer(*location)
	// get locations
	pkg.Locations = append(pkg.Locations, model.Location{
		Path:      pkg.Path,
		LayerHash: location.LayerHash,
	})

	return &pkg
}

// Parse Go Bin Build Settings
func parseBuildSettings(settings []debug.BuildSetting, key string) string {
	for _, setting := range settings {
		if setting.Key == key {
			return setting.Value
		}
	}
	return ""
}

// Get and Evaluate Version
func getVersion(settings []debug.BuildSetting) string {
	// Generate Pseudo Version for devel versions, if applicable
	revision := parseBuildSettings(settings, "vcs.revision")
	time := parseBuildSettings(settings, "vcs.time")
	if revision == "" || time == "" {
		return "(devel)"
	}
	return pseudoVersion(time, revision)

}

// A pseudo-version is a specially formatted pre-release version that encodes information about a specific revision in a version control repository.
// Source: https://go.dev/ref/mod#pseudo-versions
func pseudoVersion(time string, revision string) string {
	// A timestamp (yyyymmddhhmmss), which is the UTC time the revision was created. In Git, this is the commit time, not the author time.
	r := regexp.MustCompile("[0-9]+")
	timestamp := strings.Join(r.FindAllString(time, -1), "")
	// A revision identifier (abcdefabcdef), which is a 12-character prefix of the commit hash, or in Subversion, a zero-padded revision number.
	identifier := revision[0:12]

	return fmt.Sprintf("v0.0.0-%+v-%+v", timestamp, identifier)
}
