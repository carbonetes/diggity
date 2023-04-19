// Package scan provides functionality for scanning Docker images.
package scan

import (
	secret "github.com/carbonetes/diggity/internal/secret"
	"github.com/carbonetes/diggity/internal/slsa"
	client "github.com/carbonetes/diggity/pkg/docker"
	"github.com/carbonetes/diggity/pkg/model"
	alpine "github.com/carbonetes/diggity/pkg/parser/alpine"
	"github.com/carbonetes/diggity/pkg/parser/bom"
	cargo "github.com/carbonetes/diggity/pkg/parser/cargo"
	composer "github.com/carbonetes/diggity/pkg/parser/composer"
	conan "github.com/carbonetes/diggity/pkg/parser/conan"
	dart "github.com/carbonetes/diggity/pkg/parser/dart"
	debian "github.com/carbonetes/diggity/pkg/parser/debian"
	distro "github.com/carbonetes/diggity/pkg/parser/distro"
	docker "github.com/carbonetes/diggity/pkg/parser/docker"
	gem "github.com/carbonetes/diggity/pkg/parser/gem"
	golang "github.com/carbonetes/diggity/pkg/parser/golang"
	hackage "github.com/carbonetes/diggity/pkg/parser/hackage"
	hex "github.com/carbonetes/diggity/pkg/parser/hex"
	java "github.com/carbonetes/diggity/pkg/parser/java"
	npm "github.com/carbonetes/diggity/pkg/parser/npm"
	nuget "github.com/carbonetes/diggity/pkg/parser/nuget"
	portage "github.com/carbonetes/diggity/pkg/parser/portage"
	python "github.com/carbonetes/diggity/pkg/parser/python"
	rpm "github.com/carbonetes/diggity/pkg/parser/rpm"
	swift "github.com/carbonetes/diggity/pkg/parser/swift"
	"github.com/carbonetes/diggity/pkg/parser/util"
	"github.com/carbonetes/diggity/pkg/provider"
)

// parsers is a slice of functions that find content from the image.
type (
	parsers []func()
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
	target := client.ExtractImage(imageId)
	bom.Target = target
	bom.InitParsers(*arguments)
	bom.WG.Add(len(FindFunctions))
	for _, parser := range FindFunctions {
		go parser()
	}
	bom.WG.Wait()
	util.CleanUp()
	output := model.Result{
		Distro:   distro.Distro(),
		Packages: bom.Packages,
	}

	if !*bom.Arguments.DisableSecretSearch {
		output.Secret = secret.SecretResults
	}

	if *bom.Arguments.Provenance != "" {
		output.SLSA = slsa.Provenance()
	}

	return &output
}
