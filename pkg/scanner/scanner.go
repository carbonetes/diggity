// Package scan provides functionality for scanning Docker images.
package scanner

import (
	"log"

	"github.com/carbonetes/diggity/internal/secret"
	"github.com/carbonetes/diggity/internal/slsa"
	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/parser/alpine"
	"github.com/carbonetes/diggity/pkg/parser/alpm"
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
)

type parsers []func(*bom.ParserRequirements)

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
		alpm.FindAlpmPackagesFromContent,
	}
)

// Diggity scans the Docker images, Tar Files, and Codebases(directories) specified in the given model.Arguments struct and returns a sbom(model.Result) struct.
func Scan(arguments *model.Arguments) (*model.SBOM, *[]error) {
	requirements, err := bom.InitParsers(arguments)
	if err != nil {
		log.Fatal(err)
	}

	requirements.WG.Add(len(FindFunctions))
	for _, parser := range FindFunctions {
		parser(requirements)
	}
	requirements.WG.Wait()
	defer util.CleanUp(*requirements.DockerTemp)

	if *arguments.Provenance != "" {
		requirements.SBOM.SLSA = slsa.Provenance(requirements)
	}

	return requirements.SBOM, requirements.Errors
}
