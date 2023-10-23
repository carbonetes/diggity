package hex

import "github.com/carbonetes/diggity/pkg/model"

// Parse PURL
func parseHexPURL(pkg *model.Package) {
	pkg.PURL = model.PURL("pkg" + ":" + "hex" + "/" + pkg.Name + "@" + pkg.Version)
}
