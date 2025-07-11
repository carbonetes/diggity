package common

import (
	"os"
	"path/filepath"
	"strings"
)

// FileWalker provides utilities for finding relevant files
type FileWalker struct{}

// NewFileWalker creates a new FileWalker instance
func NewFileWalker() *FileWalker {
	return &FileWalker{}
}

// FindFiles recursively finds files matching the given check function
func (fw *FileWalker) FindFiles(rootPath string, checkFunc func(string) (bool, error)) ([]string, error) {
	var matchedFiles []string

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			if matches, err := checkFunc(path); err == nil && matches {
				matchedFiles = append(matchedFiles, path)
			}
		}
		return nil
	})

	return matchedFiles, err
}

// VersionCleaner provides utilities for cleaning and normalizing version strings
type VersionCleaner struct{}

// NewVersionCleaner creates a new VersionCleaner instance
func NewVersionCleaner() *VersionCleaner {
	return &VersionCleaner{}
}

// CleanVersion removes common prefixes and normalizes version strings
func (vc *VersionCleaner) CleanVersion(version string, prefixes []string) string {
	version = strings.TrimSpace(version)

	// Remove common prefixes
	for _, prefix := range prefixes {
		version = strings.TrimPrefix(version, prefix)
	}

	// Handle version ranges by taking the first part
	if strings.Contains(version, ",") {
		parts := strings.Split(version, ",")
		if len(parts) > 0 {
			version = strings.TrimSpace(parts[0])
		}
	}

	// Clean again after comma split
	for _, prefix := range prefixes {
		version = strings.TrimPrefix(version, prefix)
	}

	return strings.TrimSpace(version)
}

// GetCommonVersionPrefixes returns common version prefixes for different ecosystems
func (vc *VersionCleaner) GetCommonVersionPrefixes() []string {
	return []string{"^", "~", ">=", "<=", ">", "<", "==", "!=", "="}
}

// FileChecker provides utilities for checking file types
type FileChecker struct{}

// NewFileChecker creates a new FileChecker instance
func NewFileChecker() *FileChecker {
	return &FileChecker{}
}

// CheckFileByName checks if a file matches any of the given filenames
func (fc *FileChecker) CheckFileByName(path string, filenames []string) bool {
	filename := filepath.Base(path)
	for _, name := range filenames {
		if filename == name {
			return true
		}
	}
	return false
}

// CheckFileByExtension checks if a file has any of the given extensions
func (fc *FileChecker) CheckFileByExtension(path string, extensions []string) bool {
	ext := filepath.Ext(path)
	for _, extension := range extensions {
		if ext == extension {
			return true
		}
	}
	return false
}
