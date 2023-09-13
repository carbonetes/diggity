package rpm

import (
	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/pkg/model"
)

func generateCpes(pkg *model.Package, name, vendor, version string) {
	cpe.NewCPE23(pkg, vendor, name, version)
}
