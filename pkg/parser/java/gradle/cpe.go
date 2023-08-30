package gradle

import (
	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/pkg/model"
)

func generateCpes(pkg *model.Package, vendor string) {
	cpe.NewCPE23(pkg, vendor, pkg.Name, pkg.Version)
}