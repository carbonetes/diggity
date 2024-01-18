package reader

import (
	"debug/elf"
	"os"
	"path/filepath"

	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/pkg/scanner"
	"github.com/carbonetes/diggity/pkg/stream"
	"github.com/carbonetes/diggity/pkg/types"
)

func FilesystemScanHandler(data interface{}) interface{} {
	input, ok := data.(string)
	if !ok {
		log.Error("Filesystem Handler received unknown type")
		return data
	}
	var paths []string
	// recursive
	err := filepath.Walk(input,
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
		log.Error(err)
	}
	for _, path := range paths {
		stream.Emit(stream.FileListEvent, path)
		category, matched, readFlag := scanner.CheckRelatedFiles(path)
		if matched {
			switch category {
			case "rpm":
				err := handleRpmFile(path, category)
				if err != nil {
					log.Error(err)
				}
			default:
				if !readFlag {
					stream.Emit(category, types.ManifestFile{
						Path: path,
					})
					continue
				}
				file, err := os.Open(path)
				if err != nil {
					log.Error(err)
				}
				err = handleManifestFile(path, category, file, false)
				if err != nil {
					log.Error(err)
				}
			}
		}
	}
	return data
}

func handleRpmFile(path, category string) error {
	rpmDb := types.RpmDB{
		Path: path,
	}
	err := rpmDb.ReadDBFile(path)
	if err != nil {
		return err
	}
	stream.Emit(category, rpmDb)
	return nil
}

func handleManifestFile(path, category string, file *os.File, cleanup bool) error {
	manifest := types.ManifestFile{
		Path: path,
	}
	err := manifest.ReadContent(file)
	if err != nil {
		return err
	}
	stream.Emit(category, manifest)

	return nil
}

func handleGeneric(path, category, name string) {
	f, err := elf.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	generic := types.Generic{
		Path: path,
		Name: name,
		File: f,
	}
	generic.ReadROData()
	if len(generic.ROData) == 0 {
		return
	}
	stream.Emit(category, generic)
}
