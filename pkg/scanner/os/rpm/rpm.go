package rpm

import (
	"path/filepath"
	"slices"

	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/pkg/types"
)

const Type string = "rpm"

var ManifestFiles = []string{"Packages", "Packages.db", "rpmdb.sqlite"}

func Scan(data interface{}) interface{} {
	rpmdb, ok := data.(types.RpmDB)

	if !ok {
		log.Error("RPM Handler received unknown type")
	}
	readRpmDb(rpmdb)

	return data
}

func CheckRelatedFiles(file string) (string, bool, bool) {
	if slices.Contains(ManifestFiles, filepath.Base(file)) {
		return Type, true, true
	}
	return "", false, false
}
