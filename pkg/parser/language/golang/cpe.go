package golang

import (
	"strings"

	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/pkg/model"
)

func GenerateCpes(pkg *model.Package, paths []string) {
	// check if cpePaths only contains the product
	if len(paths) > 1 {
		cpe.NewCPE23(pkg, paths[len(paths)-2], paths[len(paths)-1], FormatVersion(pkg.Version))
	} else {
		cpe.NewCPE23(pkg, "", paths[0], FormatVersion(pkg.Version))
	}
}

// Format Version String
func FormatVersion(version string) string {
	if strings.Contains(version, "(") && strings.Contains(version, ")") {
		version = strings.Replace(version, "(", "", -1)
		version = strings.Replace(version, ")", "", -1)
	}
	return version
}