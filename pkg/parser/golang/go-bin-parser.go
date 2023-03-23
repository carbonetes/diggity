package golang

import (
	"debug/buildinfo"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime/debug"
	"strings"

	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/internal/file"
	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/model/metadata"
	"github.com/carbonetes/diggity/pkg/parser/bom"
	"github.com/carbonetes/diggity/pkg/parser/util"

	"github.com/google/uuid"
)

const (
	goArch      = "GOARCH"
	goOS        = "GOOS"
	goCompiler  = "-compiler"
	devel       = "devel"
	vcsRevision = "vcs.revision"
	vcsTime     = "vcs.time"
)

// GoBinMetadata GoMetadata  metadata
type GoBinMetadata map[string]interface{}

// FindGoBinPackagesFromContent Find go binaries in the file contents
func FindGoBinPackagesFromContent() {
	// Look for go bin file
	if util.ParserEnabled(goType) {
		for _, content := range file.Contents {
			if !strings.Contains(filepath.Base(content.Path), ".") {
				if err := readGoBinContent(content); err != nil {
					err = errors.New("go-bin-parser: " + err.Error())
					bom.Errors = append(bom.Errors, &err)
				}
			}
		}
	}
	defer bom.WG.Done()
}

// Read go binaries content
func readGoBinContent(location *model.Location) error {
	// Modify file permissions to allow read
	err := os.Chmod(location.Path, 0777)
	if err != nil {
		if strings.Contains(err.Error(), noFileErrWin) || strings.Contains(err.Error(), noFileErrMac) {
			return nil
		}
		return err
	}

	goBinFile, err := os.Open(location.Path)
	if err != nil {
		if strings.Contains(err.Error(), noFileErrWin) || strings.Contains(err.Error(), noFileErrMac) {
			return nil
		}
		return err
	}
	defer goBinFile.Close()

	buildData, err := buildinfo.Read(goBinFile)

	// Check if file is Go bin
	if err != nil {
		// Handle expected errors
		if err.Error() == "unrecognized file format" ||
			err.Error() == "not a Go executable" ||
			err.Error() == "EOF" {
			return nil
		}
		return err
	}

	if buildData == nil {
		return nil
	}

	// Parse mod metadata from bin
	if buildData.Main.Path != "" {
		appendGoBinPackage(location, buildData, &buildData.Main)
	}

	// Parse dependencies
	if len(buildData.Deps) > 0 {
		for _, dep := range buildData.Deps {
			if dep.Replace != nil {
				appendGoBinPackage(location, buildData, dep.Replace)
			}
			appendGoBinPackage(location, buildData, dep)
		}
	} else {
		if buildData.Main.Path != "" {
			appendGoBinPackage(location, buildData, nil)
		}
	}
	return nil
}

// Append Go Bin Package to Packages list
func appendGoBinPackage(location *model.Location, buildData *debug.BuildInfo, dep *debug.Module) {
	_package := new(model.Package)
	initGoBinPackage(_package, location, buildData, dep)
	bom.Packages = append(bom.Packages, _package)
}

// Initialize Go Bin package contents
func initGoBinPackage(p *model.Package, location *model.Location, buildData *debug.BuildInfo, dep *debug.Module) *model.Package {
	p.ID = uuid.NewString()
	p.Type = goModule
	p.Licenses = []string{}

	// get locations
	p.Locations = append(p.Locations, model.Location{
		Path:      util.TrimUntilLayer(*location),
		LayerHash: location.LayerHash,
	})

	if dep != nil {
		p.Name = dep.Path
		p.Version = getVersion(dep.Version, buildData.Settings)
	} else {
		p.Name = buildData.Main.Path
		p.Version = buildData.Main.Version
	}
	p.Path = p.Name

	// get and format CPEs
	cpePaths := splitPath(p.Name)

	// check if cpePaths only contains the product
	if len(cpePaths) > 1 {
		cpe.NewCPE23(p, cpePaths[len(cpePaths)-2], cpePaths[len(cpePaths)-1], formatVersion(p.Version))
	} else {
		cpe.NewCPE23(p, "", cpePaths[0], formatVersion(p.Version))
	}

	p.CPEs = formatDevelCPEs(p)

	// get purl
	parseGoPackageURL(p)

	// parse and fill final metadata
	initGoBinMetadata(p, dep, buildData)

	return p
}

// Initialize Go metadata values from content
func initGoBinMetadata(_package *model.Package, dep *debug.Module, buildData *debug.BuildInfo) {
	var finalMetadata = metadata.GoBinMetadata{}

	if dep != nil {
		if len(buildData.Settings) > 0 {
			finalMetadata.Architecture = parseBuildSettings(buildData.Settings, goArch)
			finalMetadata.Compiler = parseBuildSettings(buildData.Settings, goCompiler)
			finalMetadata.OS = parseBuildSettings(buildData.Settings, goOS)
		}
		finalMetadata.GoCompileRelease = buildData.GoVersion
		finalMetadata.H1Digest = dep.Sum
		finalMetadata.Path = dep.Path
		finalMetadata.Version = dep.Version

		_package.Metadata = finalMetadata
	} else {
		finalMetadata.GoCompileRelease = buildData.GoVersion
		finalMetadata.Path = buildData.Main.Path
		finalMetadata.Version = buildData.Main.Version
	}
	_package.Metadata = finalMetadata
}

// Format CPEs for Devel Versions
func formatDevelCPEs(_package *model.Package) []string {
	var newCPEs []string
	if len(_package.CPEs) > 0 {
		for _, cpe := range _package.CPEs {
			if strings.Contains(cpe, devel) {
				newCPEs = append(newCPEs, strings.Replace(cpe, devel, `\(devel\)`, -1))
			} else {
				newCPEs = append(newCPEs, cpe)
			}
		}
	}

	return newCPEs
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

// Format Version String
func formatVersion(version string) string {
	if strings.Contains(version, "(") && strings.Contains(version, ")") {
		version = strings.Replace(version, "(", "", -1)
		version = strings.Replace(version, ")", "", -1)
	}
	return version
}

// Get and Evaluate Version
func getVersion(version string, settings []debug.BuildSetting) string {
	// Generate Pseudo Version for devel versions, if applicable
	if strings.Contains(version, devel) {
		vRevision := parseBuildSettings(settings, vcsRevision)
		vTime := parseBuildSettings(settings, vcsTime)
		if vRevision == "" || vTime == "" {
			return version
		}
		return pseudoVersion(vTime, vRevision)

	}
	return version
}

// A pseudo-version is a specially formatted pre-release version that encodes information about a specific revision in a version control repository.
// Source: https://go.dev/ref/mod#pseudo-versions
func pseudoVersion(vTime string, vRevision string) string {
	// A timestamp (yyyymmddhhmmss), which is the UTC time the revision was created. In Git, this is the commit time, not the author time.
	r := regexp.MustCompile("[0-9]+")
	timestamp := strings.Join(r.FindAllString(vTime, -1), "")
	// A revision identifier (abcdefabcdef), which is a 12-character prefix of the commit hash, or in Subversion, a zero-padded revision number.
	revisionIdentifier := vRevision[0:12]

	return fmt.Sprintf("v0.0.0-%+v-%+v", timestamp, revisionIdentifier)
}
