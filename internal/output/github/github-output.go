package github

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/carbonetes/diggity/internal/output/save"
	"github.com/carbonetes/diggity/internal/parser/bom"
	"github.com/carbonetes/diggity/internal/parser/distro"
	versionPackage "github.com/carbonetes/diggity/internal/version"
	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/model/output"
)

const (
	appName            = "diggity"
	githubUrl          = "https://github.com/carbonetes/diggity"
	directRelationship = "direct"
	runtimeScope       = "runtime"
)

// DependencyMetadata dependency metadata
type DependencyMetadata map[string]interface{}

// PrintGithubJSON Print Packages in Github JSON format
func PrintGithubJSON() {
	githubJSON, err := getGithubJSON()

	if err != nil {
		panic(err)
	}

	if len(*bom.Arguments.OutputFile) > 0 {
		save.ResultToFile(string(githubJSON))
	} else {
		fmt.Printf("%+v\n", string(githubJSON))
	}
}

// getGithubJSON Init Github JSON Output
func getGithubJSON() ([]byte, error) {
	var image string
	if bom.Arguments.Image == nil {
		if *bom.Arguments.Tar != "" {
			image = *bom.Arguments.Tar
		}
		if *bom.Arguments.Dir != "" {
			image = *bom.Arguments.Tar
		}
	} else {
		image = strings.Split(*bom.Arguments.Image, ":")[0]
	}

	result := output.DependencySnapshot{
		Version: 0,
		Detector: output.Detector{
			Name:    appName,
			URL:     githubUrl,
			Version: versionPackage.FromBuild().Version,
		},
		Metadata:  getSnapshotMetadata(distro.Distro()),
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
	for _, pkg := range bom.Packages {
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
