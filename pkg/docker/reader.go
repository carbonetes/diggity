package docker

import (
	"github.com/carbonetes/diggity/pkg/types"
	"github.com/golistic/urn"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
)

// Reader handles Docker image operations including pulling, processing, and tar handling
type Reader struct {
	config         *types.RegistryConfig
	imageHandler   *ImageHandler
	tarballHandler *TarballHandler
	processor      *Processor
}

// NewReader creates a new Docker reader instance
func NewReader(registryConfig *types.RegistryConfig) *Reader {
	processor := NewProcessor()
	return &Reader{
		config:         registryConfig,
		imageHandler:   NewImageHandler(registryConfig),
		tarballHandler: NewTarballHandler(),
		processor:      processor,
	}
}

// ScanImageHandler handles the scanning of Docker images
func (r *Reader) ScanImageHandler(target string, addr *urn.URN) error {
	image, ref, err := r.imageHandler.GetImage(target)
	if err != nil {
		return err
	}

	return r.processor.ProcessImage(image, ref, addr)
}

// ScanTarballHandler handles the scanning of Docker tarball files
func (r *Reader) ScanTarballHandler(target string, addr *urn.URN) error {
	image, err := r.tarballHandler.ReadTarball(target)
	if err != nil {
		return err
	}

	// Create a dummy reference for tarball
	ref, _ := name.ParseReference("tarball:latest")

	return r.processor.ProcessImage(image, &ref, addr)
}

// GetImage retrieves a Docker image (delegates to ImageHandler)
func (r *Reader) GetImage(input string) (*v1.Image, *name.Reference, error) {
	return r.imageHandler.GetImage(input)
}

// ReadTarball reads a tarball (delegates to TarballHandler)
func (r *Reader) ReadTarball(path string) (*v1.Image, error) {
	return r.tarballHandler.ReadTarball(path)
}

// ProcessImage processes a Docker image (delegates to Processor)
func (r *Reader) ProcessImage(image *v1.Image, ref *name.Reference, addr *urn.URN) error {
	return r.processor.ProcessImage(image, ref, addr)
}
