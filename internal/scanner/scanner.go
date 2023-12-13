package scanner

import (
	"github.com/carbonetes/diggity/internal/scanner/distro"
	"github.com/carbonetes/diggity/internal/scanner/language/cargo"
	"github.com/carbonetes/diggity/internal/scanner/language/cocoapods"
	"github.com/carbonetes/diggity/internal/scanner/language/composer"
	"github.com/carbonetes/diggity/internal/scanner/language/conan"
	"github.com/carbonetes/diggity/internal/scanner/language/cran"
	"github.com/carbonetes/diggity/internal/scanner/language/golang"
	"github.com/carbonetes/diggity/internal/scanner/language/gradle"
	"github.com/carbonetes/diggity/internal/scanner/language/hackage"
	"github.com/carbonetes/diggity/internal/scanner/language/hex"
	"github.com/carbonetes/diggity/internal/scanner/language/maven"
	"github.com/carbonetes/diggity/internal/scanner/language/npm"
	"github.com/carbonetes/diggity/internal/scanner/language/nuget"
	"github.com/carbonetes/diggity/internal/scanner/language/pub"
	"github.com/carbonetes/diggity/internal/scanner/language/pypi"
	"github.com/carbonetes/diggity/internal/scanner/language/rubygem"
	"github.com/carbonetes/diggity/internal/scanner/language/swift"
	"github.com/carbonetes/diggity/internal/scanner/os/apk"
	"github.com/carbonetes/diggity/internal/scanner/os/dpkg"
	"github.com/carbonetes/diggity/internal/scanner/os/rpm"
	"github.com/carbonetes/diggity/internal/scanner/secret"
	"github.com/carbonetes/diggity/pkg/stream"
)

type FileChecker func(file string) (string, bool, bool)

var All = []string{
	apk.Type,
	cargo.Type,
	cran.Type,
	cocoapods.Type,
	composer.Type,
	conan.Type,
	distro.Type,
	dpkg.Type,
	golang.Type,
	gradle.Type,
	hackage.Type,
	hex.Type,
	maven.Type,
	npm.Type,
	nuget.Type,
	pub.Type,
	pypi.Type,
	rpm.Type,
	rubygem.Type,
	secret.Type,
	swift.Type,
}

var FileCheckers = []FileChecker{
	apk.CheckRelatedFile,
	cargo.CheckRelatedFile,
	cran.CheckRelatedFiles,
	cocoapods.CheckRelatedFile,
	composer.CheckRelatedFile,
	conan.CheckRelatedFile,
	distro.CheckRelatedFile,
	dpkg.CheckRelatedFile,
	golang.CheckRelatedFile,
	gradle.CheckRelatedFile,
	hackage.CheckRelatedFile,
	hex.CheckRelatedFile,
	maven.CheckRelatedFile,
	npm.CheckRelatedFile,
	nuget.CheckRelatedFile,
	pub.CheckRelatedFile,
	pypi.CheckRelatedFile,
	rpm.CheckRelatedFiles,
	rubygem.CheckRelatedFile,
	secret.CheckRelatedFile,
	swift.CheckRelatedFile,
}

func init() {
	stream.Attach(apk.Type, apk.Scan)
	stream.Attach(cargo.Type, cargo.Scan)
	stream.Attach(cran.Type, cran.Scan)
	stream.Attach(cocoapods.Type, cocoapods.Scan)
	stream.Attach(composer.Type, composer.Scan)
	stream.Attach(conan.Type, conan.Scan)
	stream.Attach(distro.Type, distro.Scan)
	stream.Attach(dpkg.Type, dpkg.Scan)
	stream.Attach(golang.Type, golang.Scan)
	stream.Attach(gradle.Type, gradle.Scan)
	stream.Attach(hackage.Type, hackage.Scan)
	stream.Attach(hex.Type, hex.Scan)
	stream.Attach(maven.Type, maven.Scan)
	stream.Attach(npm.Type, npm.Scan)
	stream.Attach(nuget.Type, nuget.Scan)
	stream.Attach(pub.Type, pub.Scan)
	stream.Attach(pypi.Type, pypi.Scan)
	stream.Attach(rpm.Type, rpm.Scan)
	stream.Attach(rubygem.Type, rubygem.Scan)
	stream.Attach(secret.Type, secret.Scan)
	stream.Attach(swift.Type, swift.Scan)
}

func CheckRelatedFiles(file string) (string, bool, bool) {
	for _, checker := range FileCheckers {
		category, matched, readFlag := checker(file)
		if matched {
			return category, matched, readFlag
		}
	}
	return "", false, false
}
