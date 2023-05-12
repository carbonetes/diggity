package alpm

import (
	"github.com/carbonetes/diggity/pkg/model"
	"github.com/google/uuid"
)

func newPackage(metadata Metadata) *model.Package {
	if metadata["Name"] == nil || metadata["Name"] == "" {
		return nil
	}

	return &model.Package{
		ID:          uuid.NewString(),
		Name:        metadata["Name"].(string),
		Type:        Type,
		Version:     metadata["Version"].(string),
		Description: metadata["Description"].(string),
		Metadata:    metadata,
	}
}
