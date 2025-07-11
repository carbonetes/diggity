package filesystem

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
)

// FileHandler handles individual file operations
type FileHandler struct{}

// NewFileHandler creates a new FileHandler instance
func NewFileHandler() *FileHandler {
	return &FileHandler{}
}

// ReadFile reads a single file and returns its contents
func (h *FileHandler) ReadFile(filePath string) ([]byte, error) {
	return os.ReadFile(filePath)
}

// WriteFile writes data to a file
func (h *FileHandler) WriteFile(filePath string, data []byte, perm os.FileMode) error {
	return os.WriteFile(filePath, data, perm)
}

// FileExists checks if a file exists
func (h *FileHandler) FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

// GetFileInfo returns file information
func (h *FileHandler) GetFileInfo(filePath string) (*FileInfo, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	return &FileInfo{
		Name:    filepath.Base(filePath),
		Path:    filePath,
		Size:    info.Size(),
		ModTime: info.ModTime(),
		IsDir:   info.IsDir(),
		Mode:    info.Mode(),
	}, nil
}

// ListDirectory returns the contents of a directory
func (h *FileHandler) ListDirectory(dirPath string) ([]FileInfo, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	var files []FileInfo
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		files = append(files, FileInfo{
			Name:    entry.Name(),
			Path:    filepath.Join(dirPath, entry.Name()),
			Size:    info.Size(),
			ModTime: info.ModTime(),
			IsDir:   entry.IsDir(),
			Mode:    info.Mode(),
		})
	}

	return files, nil
}

// CreateTempFile creates a temporary file with the given content
func (h *FileHandler) CreateTempFile(pattern string, content []byte) (string, error) {
	tmpFile, err := os.CreateTemp("", pattern)
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	_, err = tmpFile.Write(content)
	if err != nil {
		os.Remove(tmpFile.Name())
		return "", err
	}

	return tmpFile.Name(), nil
}

// RemoveFile removes a file or directory
func (h *FileHandler) RemoveFile(filePath string) error {
	return os.Remove(filePath)
}

// RemoveAll removes a file or directory and all its children
func (h *FileHandler) RemoveAll(filePath string) error {
	return os.RemoveAll(filePath)
}

// GetFileContent reads and returns file content as a string
func (h *FileHandler) GetFileContent(filePath string) (string, error) {
	data, err := h.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// GetFileReader returns an io.Reader for the file
func (h *FileHandler) GetFileReader(filePath string) (io.Reader, error) {
	data, err := h.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(data), nil
}

// GetFileSize returns the size of a file
func (h *FileHandler) GetFileSize(filePath string) (int64, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

// IsDirectory checks if a path is a directory
func (h *FileHandler) IsDirectory(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// CreateDirectory creates a directory and any necessary parent directories
func (h *FileHandler) CreateDirectory(dirPath string, perm os.FileMode) error {
	return os.MkdirAll(dirPath, perm)
}
