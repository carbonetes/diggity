// Package scan provides functionality for scanning Docker images.
package scanner

import (
	secret "github.com/carbonetes/diggity/internal/secret"
	"github.com/carbonetes/diggity/internal/slsa"
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

// Diggity scans the Docker images, Tar Files, and Codebases(directories) specified in the given model.Arguments struct and returns a sbom(model.Result) struct.
func Scan(arguments *model.Arguments) *model.Result {
    // Initialize parsers with the given arguments.
    bom.InitParsers(*arguments)
    
    // Add the number of FindFunctions to the WaitGroup.
    bom.WG.Add(len(FindFunctions))
    
    // For each parser in FindFunctions, run it concurrently.
    for _, parser := range FindFunctions {
        go parser()
    }
    
    // Wait for all parsers to finish.
    bom.WG.Wait()
    
    // Clean up any temporary files created during parsing.
    util.CleanUp()
    
    // Create a new Result struct with the Distro and Packages fields set.
    output := model.Result{
        Distro:   distro.Distro(),
        Packages: bom.Packages,
    }

    // If secret search is not disabled, add the SecretResults field to the output.
    if !*bom.Arguments.DisableSecretSearch {
        output.Secret = secret.SecretResults
    }

    // If provenance is specified, add the SLSA field to the output.
    if *bom.Arguments.Provenance != "" {
        output.SLSA = slsa.Provenance()
    }

    // Return a pointer to the output struct.
    return &output
}
