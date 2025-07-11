package archive

import (
	"io"

	"github.com/golistic/urn"
)

// Reader handles archive file operations including zip, jar, war, and ear files
type Reader struct {
	supportedTypes []string
	processor      *Processor
	extractor      *Extractor
}

// NewReader creates a new archive reader instance
func NewReader() *Reader {
	return &Reader{
		supportedTypes: []string{".jar", ".war", ".ear", ".jpi", ".hpi", ".zip"},
		processor:      NewProcessor(),
		extractor:      NewExtractor(),
	}
}

// GetSupportedTypes returns the list of supported archive file types
func (r *Reader) GetSupportedTypes() []string {
	return r.supportedTypes
}

// IsArchiveFile checks if a file is a supported archive type
func (r *Reader) IsArchiveFile(filename string) bool {
	return r.processor.IsArchiveFile(filename, r.supportedTypes)
}

// ProcessArchive processes an archive file and extracts its contents for scanning
func (r *Reader) ProcessArchive(reader io.ReaderAt, path string, size int64, addr *urn.URN) error {
	return r.processor.ProcessArchive(reader, path, size, addr)
}

// ProcessArchiveFromPath processes an archive file from a file path
func (r *Reader) ProcessArchiveFromPath(filePath string, addr *urn.URN) error {
	return r.processor.ProcessArchiveFromPath(filePath, addr)
}

// ExtractFileFromArchive extracts a specific file from an archive
func (r *Reader) ExtractFileFromArchive(archivePath, targetFile string) ([]byte, error) {
	return r.extractor.ExtractFileFromArchive(archivePath, targetFile)
}

// ListArchiveContents returns a list of all files in an archive
func (r *Reader) ListArchiveContents(archivePath string) ([]string, error) {
	return r.extractor.ListArchiveContents(archivePath)
}

// GetArchiveMetadata returns metadata about an archive file
func (r *Reader) GetArchiveMetadata(archivePath string) (*ArchiveMetadata, error) {
	return r.extractor.GetArchiveMetadata(archivePath)
}
