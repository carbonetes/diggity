package parser

import (
	alpine "github.com/carbonetes/diggity/internal/parser/alpine"
	cargo "github.com/carbonetes/diggity/internal/parser/cargo"
	composer "github.com/carbonetes/diggity/internal/parser/composer"
	conan "github.com/carbonetes/diggity/internal/parser/conan"
	dart "github.com/carbonetes/diggity/internal/parser/dart"
	debian "github.com/carbonetes/diggity/internal/parser/debian"
	distro "github.com/carbonetes/diggity/internal/parser/distro"
	docker "github.com/carbonetes/diggity/internal/parser/docker"
	gem "github.com/carbonetes/diggity/internal/parser/gem"
	golang "github.com/carbonetes/diggity/internal/parser/golang"
	hackage "github.com/carbonetes/diggity/internal/parser/hackage"
	hex "github.com/carbonetes/diggity/internal/parser/hex"
	java "github.com/carbonetes/diggity/internal/parser/java"
	npm "github.com/carbonetes/diggity/internal/parser/npm"
	nuget "github.com/carbonetes/diggity/internal/parser/nuget"
	portage "github.com/carbonetes/diggity/internal/parser/portage"
	python "github.com/carbonetes/diggity/internal/parser/python"
	rpm "github.com/carbonetes/diggity/internal/parser/rpm"
	swift "github.com/carbonetes/diggity/internal/parser/swift"
	secret "github.com/carbonetes/diggity/internal/secret"
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
