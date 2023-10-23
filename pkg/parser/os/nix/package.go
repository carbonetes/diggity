package nix

import (
	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/model/metadata"
	"github.com/google/uuid"
)

func newPackage(metadata *metadata.NixMetadata) *model.Package {
	return &model.Package{
		ID:            uuid.NewString(),
		Name:          metadata.Name,
		Version:       metadata.Version,
		PackageOrigin: model.OSPackage,
		Parser:        Type,
		Type:          Type,
		PURL:          setPURL(metadata.Name, metadata.Version),
		Metadata:      metadata,
	}
}

func setPURL(name, version string) model.PURL {
	return model.PURL("pkg" + ":" + "hex" + "/" + name + "@" + version)
}
