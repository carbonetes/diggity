package cargo

import (
	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/parser/util"
	"github.com/google/uuid"
)

func newPackage(location *model.Location, metadata Metadata) *model.Package {
	var pkg model.Package
	pkg.ID = uuid.NewString()
	if metadata["Name"] == nil {
		return nil
	}
	pkg.Name = metadata["Name"].(string)
	pkg.Version = metadata["Version"].(string)
	pkg.Path = util.TrimUntilLayer(*location)
	pkg.Type = Type
	pkg.Locations = append(pkg.Locations, model.Location{
		Path:      util.TrimUntilLayer(*location),
		LayerHash: location.LayerHash,
	})
	pkg.Metadata = metadata
	pkg.PURL = model.PURL("pkg:cargo/" + pkg.Name + "@" + pkg.Version)

	generateCargoCpes(&pkg)

	return &pkg
}
