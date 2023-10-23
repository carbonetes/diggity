package pub

import (
	"strings"

	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/pkg/model"
)

func generateCpes(pkg *model.Package, metadata *Metadata) {
	if metadata == nil {
		cpe.NewCPE23(pkg, pkg.Name, pkg.Name, pkg.Version)
		return
	}

	if val, ok := (*metadata)["author"].(string); ok {
		cpe.NewCPE23(pkg, strings.TrimSpace(val), pkg.Name, pkg.Version)
	} else if val, ok := (*metadata)["authors"].(string); ok {
		cpe.NewCPE23(pkg, strings.TrimSpace(val), pkg.Name, pkg.Version)
	} else {
		cpe.NewCPE23(pkg, pkg.Name, pkg.Name, pkg.Version)
	}
}
