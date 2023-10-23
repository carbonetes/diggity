package apk

import (
	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/pkg/model"
)

func generateAlpineCpes(pkg *model.Package) {
	if len(pkg.Name) == 0 {
		return
	} 
	if len(pkg.Version) == 0 {
		return
	}
	cpe.NewCPE23(pkg, pkg.Name, pkg.Name, pkg.Version)
}
