package types

import (
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/types"
)

type ImageManifest struct {
	Digest     v1.Hash         `json:"digest,omitempty"`
	MediaType  types.MediaType `json:"mediatype,omitempty"`
	Size       int64           `json:"size,omitempty"`
	Manifest   v1.Manifest     `json:"manifest,omitempty"`
	ConfigFile v1.ConfigFile   `json:"config_file,omitempty"`
	Layers     []Layer         `json:"layers,omitempty"`
}

type Layer struct {
	Digest    v1.Hash         `json:"digest,omitempty"`
	DiffID    v1.Hash         `json:"diff_id,omitempty"`
	Size      int64           `json:"size,omitempty"`
	MediaType types.MediaType `json:"media_type,omitempty"`
}
