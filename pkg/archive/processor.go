package archive

import (
	"archive/zip"
	"bytes"
	"io"
	"path/filepath"
	"slices"

	"github.com/carbonetes/diggity/cmd/diggity/ui"
	"github.com/carbonetes/diggity/internal/log"
	"github.com/golistic/urn"
)

// Processor handles archive file processing and extraction logic
type Processor struct{}

// NewProcessor creates a new Processor instance
func NewProcessor() *Processor {
	return &Processor{}
}

// IsArchiveFile checks if a file is a supported archive type
func (p *Processor) IsArchiveFile(filename string, supportedTypes []string) bool {
	ext := filepath.Ext(filename)
	return slices.Contains(supportedTypes, ext)
}

// ProcessArchive processes an archive file and extracts its contents for scanning
func (p *Processor) ProcessArchive(reader io.ReaderAt, path string, size int64, addr *urn.URN) error {
	valid, zipReader := p.isValidZip(reader, size)
	if !valid {
		return ErrInvalidArchive
	}

	for _, file := range zipReader.File {
		ui.AddFile(file.Name)

		if file.FileInfo().IsDir() {
			continue
		}

		if err := p.processArchiveFile(file, path, addr); err != nil {
			log.Debug("Error processing archive file:", err)
			continue
		}
	}

	return nil
}

// ProcessArchiveFromPath processes an archive file from a file path
func (p *Processor) ProcessArchiveFromPath(filePath string, addr *urn.URN) error {
	zipReader, err := zip.OpenReader(filePath)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	for _, file := range zipReader.File {
		ui.AddFile(file.Name)

		if file.FileInfo().IsDir() {
			continue
		}

		if err := p.processArchiveFile(file, filePath, addr); err != nil {
			log.Debug("Error processing archive file:", err)
			continue
		}
	}

	return nil
}

// processArchiveFile processes individual files within an archive
func (p *Processor) processArchiveFile(file *zip.File, parentPath string, addr *urn.URN) error {
	rc, err := file.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	data, err := io.ReadAll(rc)
	if err != nil {
		return err
	}

	// Emit file event for scanning
	p.emitFileEvent(data, file.Name, parentPath, addr)

	// Recursively process nested archives
	if p.IsArchiveFile(file.Name, []string{".jar", ".war", ".ear", ".jpi", ".hpi", ".zip"}) {
		return p.processNestedArchive(data, file.Name, addr)
	}

	return nil
}

// processNestedArchive processes archives found within other archives
func (p *Processor) processNestedArchive(data []byte, filePath string, addr *urn.URN) error {
	reader := bytes.NewReader(data)
	return p.ProcessArchive(reader, filePath, int64(len(data)), addr)
}

// emitFileEvent emits file events - grove system removed
func (p *Processor) emitFileEvent(data []byte, filePath, parentPath string, addr *urn.URN) {
	// Grove system removed - this is now a no-op
	log.Debugf("Processing archive file: %s (parent: %s, size: %d)", filePath, parentPath, len(data))
}

// isValidZip checks if the provided reader contains a valid zip archive
func (p *Processor) isValidZip(reader io.ReaderAt, size int64) (bool, *zip.Reader) {
	zipReader, err := zip.NewReader(reader, size)
	if err != nil {
		return false, nil
	}
	return true, zipReader
}
