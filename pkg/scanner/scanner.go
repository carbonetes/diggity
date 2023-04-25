// Package scan provides functionality for scanning Docker images.
package scan

import (
	"github.com/carbonetes/diggity/internal/secret"
	"github.com/carbonetes/diggity/internal/slsa"
	client "github.com/carbonetes/diggity/pkg/docker"
	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/parser/alpine"
	"github.com/carbonetes/diggity/pkg/parser/bom"
	"github.com/carbonetes/diggity/pkg/parser/cargo"
	"github.com/carbonetes/diggity/pkg/parser/composer"
	"github.com/carbonetes/diggity/pkg/parser/conan"
	"github.com/carbonetes/diggity/pkg/parser/dart"
	"github.com/carbonetes/diggity/pkg/parser/debian"
	"github.com/carbonetes/diggity/pkg/parser/distro"
	"github.com/carbonetes/diggity/pkg/parser/docker"
	"github.com/carbonetes/diggity/pkg/parser/gem"
	"github.com/carbonetes/diggity/pkg/parser/golang"
	"github.com/carbonetes/diggity/pkg/parser/hackage"
	"github.com/carbonetes/diggity/pkg/parser/hex"
	"github.com/carbonetes/diggity/pkg/parser/java"
	"github.com/carbonetes/diggity/pkg/parser/npm"
	"github.com/carbonetes/diggity/pkg/parser/nuget"
	"github.com/carbonetes/diggity/pkg/parser/portage"
	"github.com/carbonetes/diggity/pkg/parser/python"
	"github.com/carbonetes/diggity/pkg/parser/rpm"
	"github.com/carbonetes/diggity/pkg/parser/swift"
	"github.com/carbonetes/diggity/pkg/parser/util"
	"github.com/carbonetes/diggity/pkg/provider"
)

// parsers is a slice of functions that find content from the image.
type (
	parsers []func(*bom.ParserRequirements)
)

var (
	// FindFunctions is a collection of the find content functions of all parsers.
	FindFunctions = parsers{
		alpine.FindAlpinePackagesFromContent,
		debian.FindDebianPackagesFromContent,
		java.FindJavaPackagesFromContent,
		npm.FindNpmPackagesFromContent,
		composer.FindComposerPackagesFromContent,
		python.FindPythonPackagesFromContent,
		gem.FindGemPackagesFromContent,
		rpm.FindRpmPackagesFromContent,
		dart.FindDartPackagesFromContent,
		nuget.FindNugetPackagesFromContent,
		golang.FindGoModPackagesFromContent,
		golang.FindGoBinPackagesFromContent,
		hackage.FindHackagePackagesFromContent,
		cargo.FindCargoPackagesFromContent,
		conan.FindConanPackagesFromContent,
		portage.FindPortagePackagesFromContent,
		hex.FindHexPackagesFromContent,
		swift.FindSwiftPackagesFromContent,
		distro.ParseDistro,
		docker.ParseDockerProperties,
		secret.Search,
	}
)

// ScanImage scans the Docker image specified in the given model.Arguments struct and returns a model.Result struct.
func ScanImage(arguments *model.Arguments) *model.Result {
	credential := provider.NewRegistryAuth(arguments)
	imageId := client.GetImageID(arguments.Image, credential)
	contents, dir := client.ExtractImage(imageId)
	requirements := bom.NewParserRequirements(arguments, dir, contents)
	requirements.WG.Add(len(FindFunctions))
	for _, parser := range FindFunctions {
		go parser(requirements)
	}
	requirements.WG.Wait()
	util.CleanUp()
	result := requirements.Result

	if !*arguments.DisableSecretSearch {
		result.Secret = secret.SecretResults
	}

	if *arguments.Provenance != "" {
		result.SLSA = slsa.Provenance()
	}

	return result
}
