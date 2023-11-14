package scanner

import (
	"github.com/carbonetes/diggity/internal/scanner/distro"
	"github.com/carbonetes/diggity/internal/scanner/os/apk"
	"github.com/carbonetes/diggity/internal/scanner/os/dpkg"
	"github.com/carbonetes/diggity/internal/scanner/secret"
)

type FileChecker func(file string) (string, bool)

var FileCheckers = []FileChecker{
	apk.CheckRelatedFile,
	dpkg.CheckRelatedFile,
	distro.CheckRelatedFile,
	secret.CheckRelatedFile,
}

func CheckRelatedFiles(file string) (string, bool) {
	for _, checker := range FileCheckers {
		_type, matched := checker(file)
		if matched {
			return _type, true
		}
	}
	return "", false
}
