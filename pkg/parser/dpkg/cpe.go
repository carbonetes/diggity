package dpkg

import (
	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/pkg/model"
)

func generateCpes(pkg *model.Package) {
	cpe.NewCPE23(pkg, "", pkg.Name, pkg.Version)
}