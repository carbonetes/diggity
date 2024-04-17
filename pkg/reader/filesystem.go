package reader

import (
	"os"
	"path/filepath"

	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/internal/presenter/status"
	"github.com/carbonetes/diggity/pkg/scanner"
	"github.com/carbonetes/diggity/pkg/stream"
	"github.com/carbonetes/diggity/pkg/types"
)

func FilesystemScanHandler(target string, addr types.Address) error {
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

func handleRpmFile(path, category string, addr types.Address) error {
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

func handleManifestFile(path, category string, file *os.File, addr types.Address) error {
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
