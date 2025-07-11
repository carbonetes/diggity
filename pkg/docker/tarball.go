package docker

import (
	"fmt"
	"os"

	"github.com/carbonetes/diggity/internal/log"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
)

// TarballHandler handles Docker tarball operations
type TarballHandler struct{}

// NewTarballHandler creates a new TarballHandler instance
func NewTarballHandler() *TarballHandler {
	return &TarballHandler{}
}

// ReadTarball reads a tarball from the given path and returns it as a v1.Image
func (h *TarballHandler) ReadTarball(path string) (*v1.Image, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("file not found: %s", path)
	}

	image, err := tarball.ImageFromPath(path, nil)
	if err != nil {
		log.Debug(err)
		return nil, err
	}

	return &image, nil
}

// ValidateTarball checks if a tarball file is a valid Docker image tarball
func (h *TarballHandler) ValidateTarball(path string) error {
	_, err := h.ReadTarball(path)
	return err
}

// GetTarballInfo returns information about a Docker tarball
func (h *TarballHandler) GetTarballInfo(path string) (*TarballInfo, error) {
	image, err := h.ReadTarball(path)
	if err != nil {
		return nil, err
	}

	manifest, err := (*image).Manifest()
	if err != nil {
		return nil, err
	}

	config, err := (*image).ConfigFile()
	if err != nil {
		return nil, err
	}

	layers, err := (*image).Layers()
	if err != nil {
		return nil, err
	}

	var totalSize int64
	for _, layer := range layers {
		size, err := layer.Size()
		if err != nil {
			continue
		}
		totalSize += size
	}

	return &TarballInfo{
		Path:         path,
		MediaType:    string(manifest.MediaType),
		LayerCount:   len(layers),
		TotalSize:    totalSize,
		Created:      config.Created.Time,
		Architecture: config.Architecture,
		OS:           config.OS,
	}, nil
}

// ExtractTarball extracts a Docker tarball to a directory
func (h *TarballHandler) ExtractTarball(tarballPath, extractPath string) error {
	// This would implement tarball extraction logic
	// For now, we'll return a placeholder
	log.Debug("Extracting tarball:", tarballPath, "to:", extractPath)
	return fmt.Errorf("tarball extraction not yet implemented")
}
