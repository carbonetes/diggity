package types

import (
	"io"

	v1 "github.com/google/go-containerregistry/pkg/v1"
)

type ImageCurator struct {
	ImageTag string
	Image    v1.Image
	Files    []string
	Layers   []ImageLayer
}

type ImageLayer struct {
	Reader io.ReadCloser
	Hash   string
}
