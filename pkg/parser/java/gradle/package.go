package gradle

import (
	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/model/metadata"
	"github.com/google/uuid"
)

func newPackage(metadata metadata.GradleMetadata) *model.Package {
	var p model.Package

	p.ID = uuid.NewString()
	p.Name = metadata.Name
	p.Version = metadata.Version
	p.PURL = setPURL(p.Name, p.Version)
	p.Type = Type
	generateCpes(&p, metadata.Vendor)
	p.Metadata = metadata
	return &p
}

func setPURL(name, version string) model.PURL {
	return model.PURL("pkg:java" + "/" + name + "@" + version)
}