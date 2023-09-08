package settings

import (
	"github.com/carbonetes/diggity/internal/secret"
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
	"github.com/carbonetes/diggity/pkg/parser/golang/gobin"
	"github.com/carbonetes/diggity/pkg/parser/golang/gomod"
	"github.com/carbonetes/diggity/pkg/parser/hackage"
	"github.com/carbonetes/diggity/pkg/parser/hex"
	"github.com/carbonetes/diggity/pkg/parser/java/maven"
	"github.com/carbonetes/diggity/pkg/parser/java/sbt"
	"github.com/carbonetes/diggity/pkg/parser/javascript/npm"
	"github.com/carbonetes/diggity/pkg/parser/javascript/pnpm"
	"github.com/carbonetes/diggity/pkg/parser/nuget"
	"github.com/carbonetes/diggity/pkg/parser/portage"
	"github.com/carbonetes/diggity/pkg/parser/python"
	"github.com/carbonetes/diggity/pkg/parser/rpm"
	"github.com/carbonetes/diggity/pkg/parser/swift/cocoapods"
	"github.com/carbonetes/diggity/pkg/parser/swift/swiftpackagemanager"
)

type parsers []func(*bom.ParserRequirements)

var (
	// FindFunctions is a collection of the find content functions of all parsers.
	All = parsers{
		alpine.FindAlpinePackagesFromContent,
		debian.FindDebianPackagesFromContent,
		maven.FindJavaPackagesFromContent,
		npm.FindNpmPackagesFromContent,
		composer.FindComposerPackagesFromContent,
		python.FindPythonPackagesFromContent,
		gem.FindGemPackagesFromContent,
		rpm.FindRpmPackagesFromContent,
		dart.FindDartPackagesFromContent,
		nuget.FindNugetPackagesFromContent,
		gomod.FindGoModPackagesFromContent,
		gobin.FindGoBinPackagesFromContent,
		hackage.FindHackagePackagesFromContent,
		cargo.FindCargoPackagesFromContent,
		conan.FindConanPackagesFromContent,
		portage.FindPortagePackagesFromContent,
		hex.FindHexPackagesFromContent,
		cocoapods.FindSwiftPackagesFromContent,
		swiftpackagemanager.FindSwiftPackagesFromContent,
		distro.ParseDistro,
		docker.ParseDockerProperties,
		secret.Search,
		alpm.FindAlpmPackagesFromContent,
		pnpm.FindPnpmPackagesFromContent,
		sbt.FindSbtPackagesFromContent,
	}
)
