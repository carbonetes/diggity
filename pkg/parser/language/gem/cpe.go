package gem

import (
	"strings"

	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/pkg/model"
)

func generateCpes(pkg *model.Package, authors *string) {
	if authors == nil {
		cpe.NewCPE23(pkg, pkg.Name, pkg.Name, pkg.Version)
		return
	}
	cpe.NewCPE23(pkg, strings.TrimSpace(*authors), pkg.Name, pkg.Version)
}
