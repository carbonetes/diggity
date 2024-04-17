package reader

import (
	"archive/zip"
	"bytes"
	"io"
	"path/filepath"
	"slices"

	"github.com/carbonetes/diggity/internal/presenter/status"
	"github.com/carbonetes/diggity/pkg/scanner"
	"github.com/carbonetes/diggity/pkg/stream"
	"github.com/carbonetes/diggity/pkg/types"
)

// List of archive file types
var archiveTypes = []string{".jar", ".war", ".ear", ".jpi", ".hpi"}

// Process an archive file and check for manifest and related files
func processArchive(reader io.ReaderAt, size int64, addr types.Address) {
	// Check if the file is a valid zip
	// If it is, emit a FileListEvent for each file in the zip
	// If not valid, return and skip the file
	valid, r := isValidZip(reader, size)
	if !valid {
		return
	}

	// Loop through each file in the zip and emit a FileListEvent for each file
	for _, f := range r.File {
		status.AddFile(f.Name)

		// If the file is a directory, skip it
		if f.FileInfo().IsDir() {
			continue
		}

		rc, err := f.Open()
		if err != nil {
			continue
		}
		defer rc.Close()

		data, err := io.ReadAll(rc)
		if err != nil {
			continue
		}

		category, matched, readFlag := scanner.CheckRelatedFiles(f.Name)
		if matched {
			err = handleArchiveFile(f.Name, category, f, readFlag, addr)
			if err != nil {
				continue
			}
		}

		//	if the file is a valid zip, process it as a nested archive
		if slices.Contains(archiveTypes, filepath.Ext(f.Name)) {
			processArchive(bytes.NewReader(data), f.FileInfo().Size(), addr)
		}
	}
}

// handleArchiveFile processes a file in the archive and emits a manifest file event
func handleArchiveFile(path, categoty string, file *zip.File, readFlag bool, addr types.Address) error {
	payload := types.Payload{
		Address: addr,
	}

	manifest := types.ManifestFile{
		Path: path,
	}

	if readFlag {
		err := manifest.ReadArchiveFileContent(file)
		if err != nil {
			return err
		}
	}

	payload.Body = manifest

	stream.Emit(categoty, payload)
	return nil
}

// isValidZip checks if a file is a valid zip
func isValidZip(reader io.ReaderAt, size int64) (bool, *zip.Reader) {
	r, err := zip.NewReader(reader, size)
	if err != nil {
		return false, nil
	}
	return true, r
}
