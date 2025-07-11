package archive

import (
	"archive/zip"
	"io"
	"path/filepath"
)

// Extractor handles archive file extraction and metadata operations
type Extractor struct{}

// NewExtractor creates a new Extractor instance
func NewExtractor() *Extractor {
	return &Extractor{}
}

// ExtractFileFromArchive extracts a specific file from an archive
func (e *Extractor) ExtractFileFromArchive(archivePath, targetFile string) ([]byte, error) {
	zipReader, err := zip.OpenReader(archivePath)
	if err != nil {
		return nil, err
	}
	defer zipReader.Close()

	for _, file := range zipReader.File {
		if file.Name == targetFile {
			rc, err := file.Open()
			if err != nil {
				return nil, err
			}
			defer rc.Close()

			return io.ReadAll(rc)
		}
	}

	return nil, ErrFileNotFound
}

// ListArchiveContents returns a list of all files in an archive
func (e *Extractor) ListArchiveContents(archivePath string) ([]string, error) {
	zipReader, err := zip.OpenReader(archivePath)
	if err != nil {
		return nil, err
	}
	defer zipReader.Close()

	var files []string
	for _, file := range zipReader.File {
		if !file.FileInfo().IsDir() {
			files = append(files, file.Name)
		}
	}

	return files, nil
}

// GetArchiveMetadata returns metadata about an archive file
func (e *Extractor) GetArchiveMetadata(archivePath string) (*ArchiveMetadata, error) {
	zipReader, err := zip.OpenReader(archivePath)
	if err != nil {
		return nil, err
	}
	defer zipReader.Close()

	metadata := &ArchiveMetadata{
		Path:       archivePath,
		Type:       filepath.Ext(archivePath),
		TotalFiles: len(zipReader.File),
		Files:      make([]FileInfo, 0, len(zipReader.File)),
	}

	var totalSize int64
	for _, file := range zipReader.File {
		if !file.FileInfo().IsDir() {
			metadata.Files = append(metadata.Files, FileInfo{
				Name:         file.Name,
				Size:         file.FileInfo().Size(),
				ModTime:      file.FileInfo().ModTime(),
				IsCompressed: file.Method != zip.Store,
			})
			totalSize += file.FileInfo().Size()
		}
	}

	metadata.TotalSize = totalSize
	return metadata, nil
}

// ExtractAll extracts all files from an archive to a destination directory
func (e *Extractor) ExtractAll(archivePath, destDir string) error {
	// This would implement full archive extraction
	// For now, we'll return a placeholder
	return ErrArchiveCorrupted // Placeholder - not yet implemented
}

// ValidateArchive checks if an archive file is valid and not corrupted
func (e *Extractor) ValidateArchive(archivePath string) error {
	zipReader, err := zip.OpenReader(archivePath)
	if err != nil {
		return ErrInvalidArchive
	}
	defer zipReader.Close()

	// Basic validation - try to read file headers
	for _, file := range zipReader.File {
		if file.FileInfo().IsDir() {
			continue
		}

		rc, err := file.Open()
		if err != nil {
			return ErrArchiveCorrupted
		}
		rc.Close()
	}

	return nil
}
