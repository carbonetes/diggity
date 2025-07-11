package archive

import (
	"errors"
	"time"
)

// ArchiveMetadata contains metadata information about an archive file
type ArchiveMetadata struct {
	Path       string     `json:"path"`
	Type       string     `json:"type"`
	TotalFiles int        `json:"total_files"`
	TotalSize  int64      `json:"total_size"`
	Files      []FileInfo `json:"files"`
}

// FileInfo contains information about a file within an archive
type FileInfo struct {
	Name         string    `json:"name"`
	Size         int64     `json:"size"`
	ModTime      time.Time `json:"mod_time"`
	IsCompressed bool      `json:"is_compressed"`
}

// Common archive errors
var (
	// ErrInvalidArchive indicates that the archive file is invalid or corrupted
	ErrInvalidArchive = errors.New("invalid or corrupted archive file")

	// ErrFileNotFound indicates that a specific file was not found in the archive
	ErrFileNotFound = errors.New("file not found in archive")

	// ErrUnsupportedArchiveType indicates that the archive type is not supported
	ErrUnsupportedArchiveType = errors.New("unsupported archive type")

	// ErrArchiveCorrupted indicates that the archive is corrupted
	ErrArchiveCorrupted = errors.New("archive file is corrupted")

	// ErrArchiveTooBig indicates that the archive is too large to process
	ErrArchiveTooBig = errors.New("archive file is too large")

	// ErrNestedArchiveLimit indicates that nested archive depth limit has been reached
	ErrNestedArchiveLimit = errors.New("nested archive depth limit exceeded")
)
