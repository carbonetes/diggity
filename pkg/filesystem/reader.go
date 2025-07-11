package filesystem

import (
	"io"
	"os"

	"github.com/golistic/urn"
)

// Reader handles filesystem operations for scanning directories and files
type Reader struct {
	excludedDirs  []string
	excludedFiles []string
	scanner       *Scanner
	fileHandler   *FileHandler
}

// NewReader creates a new filesystem reader instance
func NewReader() *Reader {
	return &Reader{
		excludedDirs:  []string{".git", ".vscode", "node_modules", ".svn", ".hg"},
		excludedFiles: []string{},
		scanner:       NewScanner(),
		fileHandler:   NewFileHandler(),
	}
}

// SetExcludedDirectories sets the list of directories to exclude from scanning
func (r *Reader) SetExcludedDirectories(dirs []string) {
	r.excludedDirs = dirs
	r.scanner.SetExcludedDirectories(dirs)
}

// SetExcludedFiles sets the list of files to exclude from scanning
func (r *Reader) SetExcludedFiles(files []string) {
	r.excludedFiles = files
	r.scanner.SetExcludedFiles(files)
}

// ScanDirectory scans a directory and all its subdirectories for files
func (r *Reader) ScanDirectory(target string, addr *urn.URN) error {
	return r.scanner.ScanDirectory(target, addr)
}

// ReadFile reads a single file and returns its contents
func (r *Reader) ReadFile(filePath string) ([]byte, error) {
	return r.fileHandler.ReadFile(filePath)
}

// WriteFile writes data to a file
func (r *Reader) WriteFile(filePath string, data []byte, perm os.FileMode) error {
	return r.fileHandler.WriteFile(filePath, data, perm)
}

// FileExists checks if a file exists
func (r *Reader) FileExists(filePath string) bool {
	return r.fileHandler.FileExists(filePath)
}

// GetFileInfo returns file information
func (r *Reader) GetFileInfo(filePath string) (*FileInfo, error) {
	return r.fileHandler.GetFileInfo(filePath)
}

// ListDirectory returns the contents of a directory
func (r *Reader) ListDirectory(dirPath string) ([]FileInfo, error) {
	return r.fileHandler.ListDirectory(dirPath)
}

// CreateTempFile creates a temporary file with the given content
func (r *Reader) CreateTempFile(pattern string, content []byte) (string, error) {
	return r.fileHandler.CreateTempFile(pattern, content)
}

// RemoveFile removes a file or directory
func (r *Reader) RemoveFile(filePath string) error {
	return r.fileHandler.RemoveFile(filePath)
}

// RemoveAll removes a file or directory and all its children
func (r *Reader) RemoveAll(filePath string) error {
	return r.fileHandler.RemoveAll(filePath)
}

// GetFileContent reads and returns file content as a string
func (r *Reader) GetFileContent(filePath string) (string, error) {
	return r.fileHandler.GetFileContent(filePath)
}

// GetFileReader returns an io.Reader for the file
func (r *Reader) GetFileReader(filePath string) (io.Reader, error) {
	return r.fileHandler.GetFileReader(filePath)
}
