package swiftpackagemanager

import (
	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/model/metadata"
	"github.com/google/uuid"
)

func newV1Package(pin *metadata.Pin) *model.Package {
	return &model.Package{
		ID:       uuid.NewString(),
		Name:     pin.Name,
		Version:  pin.State.Version,
		Type:     Type,
		PURL:     setPurl(pin.Name, pin.State.Version),
		Metadata: *pin,
	}
}

func newV2Package(pin *metadata.Pin) *model.Package {
	return &model.Package{
		ID:       uuid.NewString(),
		Name:     pin.Identity,
		Version:  pin.State.Version,
		Type:     Type,
		PURL:     setPurl(pin.Identity, pin.State.Version),
		Metadata: *pin,
	}
}

func setPurl(name, version string) model.PURL {
	return model.PURL("pkg:swift" + "/" + name + "@" + version)
}
