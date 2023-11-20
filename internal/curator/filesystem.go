package curator

import (
	"os"
	"path/filepath"

	"github.com/carbonetes/diggity/internal/scanner"
	"github.com/carbonetes/diggity/pkg/stream"
	"github.com/carbonetes/diggity/pkg/types"
)

func FilesystemScanHandler(data interface{}) interface{} {
	fs, ok := data.(string)
	if !ok {
		log.Error("Filesystem Handler received unknown type")
		return data
	}

	err := filepath.Walk(fs, handleFile)
	if err != nil {
		log.Error(err)
	}

	return data
}

func handleFile(path string, info os.FileInfo, err error) error {
	if info.IsDir() && (info.Name() == ".git" || info.Name() == ".vscode") {
		return filepath.SkipDir
	}
	stream.Emit(stream.FilesystemScanEvent, path)
	category, matched := scanner.CheckRelatedFiles(path)
	if matched {
		if category == "rpm" {
			err = handleRpmFile(path, category)
		} else {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			err = handleManifestFile(path, category, file)
			if err != nil {
				return err
			}
		}
		if err != nil {
			return err
		}
	}

	return nil
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

func handleManifestFile(path, category string, file *os.File) error {
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
