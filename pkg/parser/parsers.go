package parser

import (
	"github.com/carbonetes/diggity/internal/output"
	secret "github.com/carbonetes/diggity/internal/secret"
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

type (
	parsers []func()
)

var (
	// FindFunctions - collection of the find content functions of all parsers
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

// Start - Run parsers
func Start(arguments *model.Arguments) {
	bom.InitParsers(*arguments)
	bom.WG.Add(len(FindFunctions))
	for _, parser := range FindFunctions {
		go parser()
	}
	bom.WG.Wait()
	util.CleanUp()
}

// GetResults for event bus
func GetResults() string {
	return output.GetResults()
}
