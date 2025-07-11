package filesystem

import (
	"os"
	"path/filepath"
	"slices"

	"github.com/carbonetes/diggity/cmd/diggity/ui"
	"github.com/carbonetes/diggity/internal/log"
	"github.com/golistic/urn"
)

// Scanner handles directory and file scanning operations
type Scanner struct {
	excludedDirs  []string
	excludedFiles []string
}

// NewScanner creates a new Scanner instance
func NewScanner() *Scanner {
	return &Scanner{
		excludedDirs:  []string{".git", ".vscode", "node_modules", ".svn", ".hg"},
		excludedFiles: []string{},
	}
}

// SetExcludedDirectories sets the list of directories to exclude from scanning
func (s *Scanner) SetExcludedDirectories(dirs []string) {
	s.excludedDirs = dirs
}

// SetExcludedFiles sets the list of files to exclude from scanning
func (s *Scanner) SetExcludedFiles(files []string) {
	s.excludedFiles = files
}

// ScanDirectory scans a directory and all its subdirectories for files
func (s *Scanner) ScanDirectory(target string, addr *urn.URN) error {
	paths, err := s.collectPaths(target)
	if err != nil {
		return err
	}

	for _, path := range paths {
		ui.AddFile(path)
		if err := s.processPath(path, addr); err != nil {
			log.Debug("Error processing path:", path, err)
		}
	}
	return nil
}

// collectPaths walks through a directory and collects all file paths
func (s *Scanner) collectPaths(target string) ([]string, error) {
	var paths []string
	err := filepath.Walk(target, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip excluded directories
		if info.IsDir() && s.isExcludedDirectory(info.Name()) {
			return filepath.SkipDir
		}

		// Skip excluded files
		if !info.IsDir() && s.isExcludedFile(filepath.Base(path)) {
			return nil
		}

		paths = append(paths, filepath.ToSlash(path))
		return nil
	})
	return paths, err
}

// processPath processes a single file or directory path
func (s *Scanner) processPath(path string, addr *urn.URN) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return nil // Directories are handled by collectPaths
	}

	// Process regular file (archive processing removed)
	return s.processRegularFile(path, addr)
}

// processRegularFile processes regular files
func (s *Scanner) processRegularFile(path string, addr *urn.URN) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// Archive processing and grove system removed - basic file processing only
	log.Debugf("Processing file: %s (size: %d)", path, len(data))
	return nil
}

// isExcludedDirectory checks if a directory should be excluded from scanning
func (s *Scanner) isExcludedDirectory(dirName string) bool {
	return slices.Contains(s.excludedDirs, dirName)
}

// isExcludedFile checks if a file should be excluded from scanning
func (s *Scanner) isExcludedFile(fileName string) bool {
	return slices.Contains(s.excludedFiles, fileName)
}
