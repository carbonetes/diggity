package output

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/carbonetes/diggity/internal/model"
	"github.com/carbonetes/diggity/internal/model/output"
	"github.com/carbonetes/diggity/internal/parser"
	versionPackage "github.com/carbonetes/diggity/internal/version"
)

const (
	appName            = "diggity"
	githubUrl          = "https://github.com/carbonetes/diggity"
	directRelationship = "direct"
	runtimeScope       = "runtime"
)

// DependencyMetadata dependency metadata
type DependencyMetadata map[string]interface{}

// printGithubJSON Print Packages in Github JSON format
func printGithubJSON() {
	githubJSON, err := getGithubJSON()

	if err != nil {
		panic(err)
	}

	if len(*parser.Arguments.OutputFile) > 0 {
		saveResultToFile(string(githubJSON))
	} else {
		fmt.Printf("%+v", string(githubJSON))
	}
}

// getGithubJSON Init Github JSON Output
func getGithubJSON() ([]byte, error) {
	var image string
	if parser.Arguments.Image == nil {
		if *parser.Arguments.Tar != "" {
			image = *parser.Arguments.Tar
		}
		if *parser.Arguments.Dir != "" {
			image = *parser.Arguments.Tar
		}
	} else {
		image = strings.Split(*parser.Arguments.Image, ":")[0]
	}

	result := output.DependencySnapshot{
		Version: 0,
		Detector: output.Detector{
			Name:    appName,
			URL:     githubUrl,
			Version: versionPackage.FromBuild().Version,
		},
		Metadata:  getSnapshotMetadata(parser.Distro()),
		Manifests: getPackageManifests(image),
		Scanned:   time.Now().Format(time.RFC3339),
	}
	return json.MarshalIndent(result, "", " ")
}

// getSnapshotMetadata returns distro metadata
func getSnapshotMetadata(distro *model.Distro) DependencyMetadata {
	if distro.ID != "" || distro.VersionID != "" {
		return DependencyMetadata{
			"diggity:distro": fmt.Sprintf("pkg:generic/%+v@%+v", distro.ID, distro.VersionID),
		}
	}
	return nil
}

// getPackageManifests returns the manifests metadata from the sbom packages discovered
func getPackageManifests(image string) output.PackageManifests {
	manifests := make(output.PackageManifests)

	// Iterate through SBOM Packages
	for _, pkg := range parser.Packages {
		paths := pkg.Locations
		for _, p := range paths {
			locPath := strings.Replace(p.Path, string(os.PathSeparator), "/", -1)
			path := fmt.Sprintf("%+v:/%+v", image, locPath)
			manifest, exists := manifests[path]

			// New manifest
			if !exists {
				manifest = output.PackageManifest{
					Name: path,
					File: output.FileInfo{
						SourceLocation: path,
					},
				}
				if p.LayerHash != "" {
					manifest.Metadata = DependencyMetadata{
						"diggity:filesystem": p.LayerHash,
					}
				}

				// Init Dependency Graph
				manifest.Resolved = map[string]output.DependencyNode{}
				manifests[path] = manifest
			}

			// Fill Dependency Graph
			purl := purlName(pkg.PURL)
			manifest.Resolved[purl] = output.DependencyNode{
				PURL:         string(pkg.PURL),
				Relationship: directRelationship,
				Scope:        runtimeScope,
			}
		}
	}

	return manifests
}

// purlName returns the manifest name from a package's PURL
func purlName(purl model.PURL) string {
	return strings.Split(string(purl), "?")[0]
}
