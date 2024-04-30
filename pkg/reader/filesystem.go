package reader

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"slices"

	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/internal/presenter/status"
	"github.com/carbonetes/diggity/pkg/scanner"
	"github.com/carbonetes/diggity/pkg/stream"
	"github.com/carbonetes/diggity/pkg/types"
	"github.com/golistic/urn"
)

func FilesystemScanHandler(target string, addr *urn.URN) error {
	var paths []string
	// recursive
	err := filepath.Walk(target,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() && (info.Name() == ".git" || info.Name() == ".vscode") {
				return filepath.SkipDir
			}

			paths = append(paths, filepath.ToSlash(path))
			return nil
		})
	if err != nil {
		return err
	}
	for _, path := range paths {
		status.AddFile(path)

		// Check if the file is an archive file (e.g. *.jar, *.war, *.ear, *.jpi, *.hpi)
		if slices.Contains(archiveTypes, filepath.Ext(path)) {
			reader, err := os.Open(path)
			if err != nil {
				continue
			}

			stat, err := reader.Stat()
			if err != nil {
				continue
			}

			b, err := io.ReadAll(reader)
			if err != nil {
				continue
			}
			processArchive(bytes.NewReader(b), stat.Size(), addr)
			continue
		}

		category, matched, readFlag := scanner.CheckRelatedFiles(path)
		if matched {
			switch category {
			case "rpm":
				err := handleRpmFile(path, category, addr)
				if err != nil {
					log.Error(err)
				}
			default:
				if !readFlag {
					stream.Emit(category, types.Payload{
						Address: addr,
						Body:    path,
					})
					continue
				}
				file, err := os.Open(path)
				if err != nil {
					log.Error(err)
				}
				err = handleManifestFile(path, category, file, addr)
				if err != nil {
					log.Error(err)
				}
			}
		}
	}
	return nil
}

func handleRpmFile(path, category string, addr *urn.URN) error {
	rpmDb := types.RpmDB{
		Path: path,
	}

	err := rpmDb.ReadDBFile(path)
	if err != nil {
		return err
	}
	stream.Emit(category, types.Payload{
		Address: addr,
		Body:    rpmDb,
	})
	return nil
}

func handleManifestFile(path, category string, file *os.File, addr *urn.URN) error {
	manifest := types.ManifestFile{
		Path: path,
	}
	err := manifest.ReadContent(file)
	if err != nil {
		return err
	}
	stream.Emit(category, types.Payload{
		Address: addr,
		Body:    manifest,
	})

	return nil
}
