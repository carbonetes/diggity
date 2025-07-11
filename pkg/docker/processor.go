package docker

import (
	"archive/tar"
	"io"
	"path/filepath"
	"slices"
	"sync"

	"github.com/carbonetes/diggity/cmd/diggity/ui"
	"github.com/carbonetes/diggity/internal/log"
	"github.com/golistic/urn"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
)

// Processor handles Docker image processing and layer analysis
type Processor struct {
	archiveExtensions []string
}

// NewProcessor creates a new Processor instance
func NewProcessor() *Processor {
	return &Processor{
		archiveExtensions: []string{".jar", ".war", ".ear", ".jpi", ".hpi", ".zip"},
	}
}

// ProcessImage processes a Docker image for scanning
func (p *Processor) ProcessImage(image *v1.Image, ref *name.Reference, addr *urn.URN) error {
	layers, err := (*image).Layers()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for _, layer := range layers {
		wg.Add(1)
		go func(layer v1.Layer) {
			defer wg.Done()
			p.processLayer(layer, addr)
		}(layer)
	}
	wg.Wait()

	return nil
}

// processLayer processes a single Docker layer
func (p *Processor) processLayer(layer v1.Layer, addr *urn.URN) {
	rc, err := layer.Uncompressed()
	if err != nil {
		log.Debug(err)
		return
	}
	defer rc.Close()

	tr := tar.NewReader(rc)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Debug(err)
			continue
		}

		p.processDockerFile(tr, header, addr)
	}
}

// processDockerFile processes individual files within Docker layers
func (p *Processor) processDockerFile(tr *tar.Reader, header *tar.Header, addr *urn.URN) {
	if header.Typeflag == tar.TypeDir {
		return
	}

	ui.AddFile(header.Name)

	data, err := io.ReadAll(tr)
	if err != nil {
		log.Debug(err)
		return
	}

	// Process archives within Docker images
	if p.isArchiveFile(header.Name) {
		p.processEmbeddedArchive(data, header.Name, addr)
	}
}

// isArchiveFile checks if a file is an archive type
func (p *Processor) isArchiveFile(filename string) bool {
	ext := filepath.Ext(filename)
	return slices.Contains(p.archiveExtensions, ext)
}

// processEmbeddedArchive processes archive files found within Docker images
func (p *Processor) processEmbeddedArchive(data []byte, path string, addr *urn.URN) {
	// Implementation would be similar to archive processing
	// This would delegate to the archive module
	log.Debug("Processing embedded archive:", path)
}

// GetLayerInfo returns information about a Docker layer
func (p *Processor) GetLayerInfo(layer v1.Layer) (*LayerInfo, error) {
	digest, err := layer.Digest()
	if err != nil {
		return nil, err
	}

	size, err := layer.Size()
	if err != nil {
		return nil, err
	}

	mediaType, err := layer.MediaType()
	if err != nil {
		return nil, err
	}

	return &LayerInfo{
		Digest:    digest.String(),
		Size:      size,
		MediaType: string(mediaType),
	}, nil
}

// AnalyzeImage performs comprehensive analysis of a Docker image
func (p *Processor) AnalyzeImage(image *v1.Image) (*ImageAnalysis, error) {
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
	layerInfos := make([]LayerInfo, 0, len(layers))

	for _, layer := range layers {
		layerInfo, err := p.GetLayerInfo(layer)
		if err != nil {
			continue
		}
		layerInfos = append(layerInfos, *layerInfo)
		totalSize += layerInfo.Size
	}

	return &ImageAnalysis{
		MediaType:    string(manifest.MediaType),
		Architecture: config.Architecture,
		OS:           config.OS,
		Created:      config.Created.Time,
		LayerCount:   len(layers),
		TotalSize:    totalSize,
		Layers:       layerInfos,
		Env:          config.Config.Env,
		Cmd:          config.Config.Cmd,
		Entrypoint:   config.Config.Entrypoint,
	}, nil
}
