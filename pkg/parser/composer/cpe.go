package composer

import (
	"strings"

	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/pkg/model"
)

func generateComposerPackageCpes(pkg *model.Package) {
	vendorProduct := strings.Split(pkg.Name, "/")
	if len(vendorProduct) == 0 {
		vendorProduct = []string{
			pkg.Name,
			pkg.Name,
		}
	}
	cpe.NewCPE23(pkg, vendorProduct[0], vendorProduct[1], pkg.Version)
}
