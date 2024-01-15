package reader

import (
	"archive/zip"
	"bytes"
	"io"
	"path/filepath"
	"slices"

	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/pkg/scanner"
	"github.com/carbonetes/diggity/pkg/stream"
	"github.com/carbonetes/diggity/pkg/types"
)

var archiveTypes = []string{".jar", ".war", ".ear", ".jpi", ".hpi"}

func processNestedArchive(reader io.ReaderAt, size int64) error {
	r, err := zip.NewReader(reader, size)
	if err != nil {
		return err
	}

	for _, f := range r.File {
		stream.Emit(stream.FileListEvent, f.Name)

		if f.FileInfo().IsDir() {
			continue
		}

		rc, err := f.Open()
		if err != nil {
			log.Error(err)
		}
		defer rc.Close()

		data, err := io.ReadAll(rc)
		if err != nil {
			log.Error(err)
		}

		if slices.Contains(archiveTypes, filepath.Ext(f.Name)) {
			err = processNestedArchive(bytes.NewReader(data), f.FileInfo().Size())
			if err != nil {
				log.Error(err)
			}
			continue
		}
		category, matched, readFlag := scanner.CheckRelatedFiles(f.Name)
		if matched {
			err = handleArchiveFile(f.Name, category, f, readFlag)
			if err != nil {
				log.Error(err)
			}
		}
	}

	return nil
}

func handleArchiveFile(path, categoty string, file *zip.File, readFlag bool) error {
	manifest := types.ManifestFile{
		Path: path,
	}
	if readFlag {
		err := manifest.ReadArchiveFileContent(file)
		if err != nil {
			return err
		}
	}

	stream.Emit(categoty, manifest)
	return nil
}
