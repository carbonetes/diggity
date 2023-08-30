package swiftpackagemanager

import (
	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/model/metadata"
)

func newV1Package(pin *metadata.Pin) *model.Package {
	var p model.Package
	p.Name = pin.Name
	p.Version = pin.State.Version
	p.Type = Type
	p.PURL = setPurl(pin.Name, pin.State.Version)
	p.Metadata = pin
	generateCpes(&p)
	return &p
}

func newV2Package(pin *metadata.Pin) *model.Package {
	var p model.Package
	p.Name = pin.Identity
	p.Version = pin.State.Version
	p.Type = Type
	p.PURL = setPurl(pin.Identity, pin.State.Version)
	p.Metadata = pin
	generateCpes(&p)
	return &p
}

func setPurl(name, version string) model.PURL {
	return model.PURL("pkg:swift" + "/" + name + "@" + version)
}
