package curator

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/carbonetes/diggity/internal/scanner"
	"github.com/carbonetes/diggity/pkg/stream"
	"github.com/carbonetes/diggity/pkg/types"
)

func FilesystemScanHandler(data interface{}) interface{} {
	input, ok := data.(string)
	if !ok {
		log.Fatal("Filesystem Handler received unknown type")
		return data
	}

	err := filepath.WalkDir(input, handler)
	if err != nil {
		log.Error(err)
	}

	return data
}

func handler(path string, di fs.DirEntry, err error) error {
	if err != nil {
		log.Fatal(err)
	}
	// if di.IsDir() && (di.Name() == ".git" || di.Name() == ".vscode") {
	// 	return nil
	// }
	// stream.Emit(stream.FilesystemScanEvent, path)
	category, matched := scanner.CheckRelatedFiles(path)
	if matched {
		if category == "rpm" {
			err = handleRpmFile(path, category)
		} else {
			file, err := os.Open(path)
			if err != nil {
				log.Fatal(err)
			}
			err = handleManifestFile(path, category, file)
			if err != nil {
				log.Fatal(err)
			}
		}
		if err != nil {
			log.Fatal(err)
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
