package scanner

import (
	"github.com/carbonetes/diggity/internal/scanner/distro"
	"github.com/carbonetes/diggity/internal/scanner/os/apk"
	"github.com/carbonetes/diggity/internal/scanner/os/dpkg"
	"github.com/carbonetes/diggity/internal/scanner/os/rpm"
	"github.com/carbonetes/diggity/internal/scanner/secret"
	"github.com/carbonetes/diggity/pkg/stream"
)

type FileChecker func(file string) (string, bool)

var All = []string{
	apk.Type,
	dpkg.Type,
	distro.Type,
	secret.Type,
	rpm.Type,
}

var FileCheckers = []FileChecker{
	apk.CheckRelatedFile,
	dpkg.CheckRelatedFile,
	rpm.CheckRelatedFiles,
	distro.CheckRelatedFile,
	secret.CheckRelatedFile,
}

func init() {
	stream.Attach(apk.Type, apk.Scan)
	stream.Attach(dpkg.Type, dpkg.Scan)
	stream.Attach(rpm.Type, rpm.Scan)
	stream.Attach(distro.Type, distro.Scan)
	stream.Attach(secret.Type, secret.Scan)
}

func CheckRelatedFiles(file string) (string, bool) {
	for _, checker := range FileCheckers {
		category, matched := checker(file)
		if matched {
			return category, true
		}
	}
	return "", false
}
