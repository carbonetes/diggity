package filesystem

import (
	"errors"
	"os"
	"time"
)

// FileInfo contains information about a file or directory
type FileInfo struct {
	Name    string      `json:"name"`
	Path    string      `json:"path"`
	Size    int64       `json:"size"`
	ModTime time.Time   `json:"mod_time"`
	IsDir   bool        `json:"is_dir"`
	Mode    os.FileMode `json:"mode"`
}

// ScanOptions contains options for filesystem scanning
type ScanOptions struct {
	ExcludedDirs      []string `json:"excluded_dirs"`
	ExcludedFiles     []string `json:"excluded_files"`
	MaxFileSize       int64    `json:"max_file_size"`
	FollowSymlinks    bool     `json:"follow_symlinks"`
	IncludeHiddenDirs bool     `json:"include_hidden_dirs"`
}

// DefaultScanOptions returns default scanning options
func DefaultScanOptions() *ScanOptions {
	return &ScanOptions{
		ExcludedDirs:      []string{".git", ".vscode", "node_modules", ".svn", ".hg"},
		ExcludedFiles:     []string{},
		MaxFileSize:       100 * 1024 * 1024, // 100MB
		FollowSymlinks:    false,
		IncludeHiddenDirs: false,
	}
}

// Common filesystem errors
var (
	// ErrFileNotFound indicates that a file was not found
	ErrFileNotFound = errors.New("file not found")

	// ErrDirectoryNotFound indicates that a directory was not found
	ErrDirectoryNotFound = errors.New("directory not found")

	// ErrPermissionDenied indicates that access to a file or directory was denied
	ErrPermissionDenied = errors.New("permission denied")

	// ErrFileTooBig indicates that a file is too large to process
	ErrFileTooBig = errors.New("file is too large")

	// ErrInvalidPath indicates that a file path is invalid
	ErrInvalidPath = errors.New("invalid file path")

	// ErrIsDirectory indicates that the specified path is a directory when a file was expected
	ErrIsDirectory = errors.New("path is a directory")

	// ErrNotDirectory indicates that the specified path is not a directory when a directory was expected
	ErrNotDirectory = errors.New("path is not a directory")
)
