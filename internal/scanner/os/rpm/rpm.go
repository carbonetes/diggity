package rpm

import (
	"path/filepath"
	"slices"

	"github.com/carbonetes/diggity/internal/logger"
	"github.com/carbonetes/diggity/pkg/types"
)

const Type string = "rpm"

var (
	ManifestFiles = []string{"Packages", "Packages.db", "rpmdb.sqlite"}
	log           = logger.GetLogger()
)

func Scan(data interface{}) interface{} {
	rpmdb, ok := data.(types.RpmDB)

	if !ok {
		log.Fatal("RPM Handler received unknown type")
	}
	readRpmDb(rpmdb)

	return data
}

func CheckRelatedFiles(file string) (string, bool) {
	if slices.Contains(ManifestFiles, filepath.Base(file)) {
		return Type, true
	}
	return "", false
}
