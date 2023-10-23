package parser

import (
	"github.com/carbonetes/diggity/internal/secret"
	"github.com/carbonetes/diggity/pkg/parser/common"
	"github.com/carbonetes/diggity/pkg/parser/distro"
	"github.com/carbonetes/diggity/pkg/parser/docker"
	"github.com/carbonetes/diggity/pkg/parser/language/cargo"
	"github.com/carbonetes/diggity/pkg/parser/language/composer"
	"github.com/carbonetes/diggity/pkg/parser/language/conan"
	"github.com/carbonetes/diggity/pkg/parser/language/gem"
	"github.com/carbonetes/diggity/pkg/parser/language/golang/gobin"
	"github.com/carbonetes/diggity/pkg/parser/language/golang/gomod"
	"github.com/carbonetes/diggity/pkg/parser/language/hackage"
	"github.com/carbonetes/diggity/pkg/parser/language/hex"
	"github.com/carbonetes/diggity/pkg/parser/language/java/gradle"
	"github.com/carbonetes/diggity/pkg/parser/language/java/maven"
	"github.com/carbonetes/diggity/pkg/parser/language/javascript/npm"
	"github.com/carbonetes/diggity/pkg/parser/language/javascript/pnpm"
	"github.com/carbonetes/diggity/pkg/parser/language/nuget"
	"github.com/carbonetes/diggity/pkg/parser/language/portage"
	"github.com/carbonetes/diggity/pkg/parser/language/pub"
	"github.com/carbonetes/diggity/pkg/parser/language/python"
	"github.com/carbonetes/diggity/pkg/parser/language/swift/cocoapods"
	"github.com/carbonetes/diggity/pkg/parser/language/swift/swiftpackagemanager"
	"github.com/carbonetes/diggity/pkg/parser/os/apk"
	"github.com/carbonetes/diggity/pkg/parser/os/dpkg"
	"github.com/carbonetes/diggity/pkg/parser/os/nix"
	"github.com/carbonetes/diggity/pkg/parser/os/pacman"
	"github.com/carbonetes/diggity/pkg/parser/os/rpm"
)

type ParserFunctions []func(*common.ParserParams)

var (
	// Parsers is a collection of the find content functions of all parsers.
	Parsers = ParserFunctions{
		apk.FindApkPackagesFromContent,
		dpkg.FindDpkgPackagesFromContent,
		maven.FindJavaPackagesFromContent,
		npm.FindNpmPackagesFromContent,
		composer.FindComposerPackagesFromContent,
		python.FindPythonPackagesFromContent,
		gem.FindGemPackagesFromContent,
		rpm.FindRpmPackagesFromContent,
		pub.FindPubPackagesFromContent,
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
		pacman.FindPacmanPackagesFromContent,
		pnpm.FindPnpmPackagesFromContent,
		gradle.FindGradlePackagesFromContent,
		nix.FindNixPackagesFromContent,
	}
)
