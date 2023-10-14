package github

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/carbonetes/diggity/internal/logger"
	"github.com/carbonetes/diggity/internal/output/json"
	"github.com/carbonetes/diggity/internal/output/save"
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

var log = logger.GetLogger()

// PrintGithubJSON Print Packages in Github JSON format
func PrintGithubJSON(args *model.Arguments, outputType *string, results *model.SBOM) {
	githubJSON, err := getGithubJSON(args, results)

	if err != nil {
		log.Error(err)
	}

	if len(*args.OutputFile) > 0 {
		save.ResultToFile(string(githubJSON), outputType, args.OutputFile)
	} else {
		fmt.Printf("%+v\n", string(githubJSON))
	}
}

// getGithubJSON Init Github JSON Output
func getGithubJSON(args *model.Arguments, results *model.SBOM) ([]byte, error) {
	var image string
	if args.Image == nil {
		if *args.Tar != "" {
			image = *args.Tar
		}
		if *args.Dir != "" {
			image = *args.Tar
		}
	} else {
		image = strings.Split(*args.Image, ":")[0]
	}

	result := output.DependencySnapshot{
		Version: 0,
		Detector: output.Detector{
			Name:    appName,
			URL:     githubUrl,
			Version: versionPackage.FromBuild().Version,
		},
		Metadata:  getSnapshotMetadata(results.Distro),
		Manifests: getPackageManifests(image, results.Packages),
		Scanned:   time.Now().Format(time.RFC3339),
	}
	jsonResult, err := json.ToJSON(result)
	if err != nil {
		return nil, err
	}
	return jsonResult, nil
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
func getPackageManifests(image string, pkgs *[]model.Package) output.PackageManifests {
	manifests := make(output.PackageManifests)

	// Iterate through SBOM Packages
	for _, pkg := range *pkgs {
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
