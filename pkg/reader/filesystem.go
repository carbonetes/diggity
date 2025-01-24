package reader

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"slices"

	stream "github.com/carbonetes/diggity/cmd/diggity/grove"
	"github.com/carbonetes/diggity/cmd/diggity/ui"
	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/pkg/scanner"
	"github.com/carbonetes/diggity/pkg/types"
	"github.com/golistic/urn"
)

func FilesystemScanHandler(target string, addr *urn.URN) error {
	paths, err := collectPaths(target)
	if err != nil {
		return err
	}

	for _, path := range paths {
		ui.AddFile(path)
		if err := processPath(path, addr); err != nil {
			log.Debug(err)
		}
	}
	return nil
}

func collectPaths(target string) ([]string, error) {
	var paths []string
	err := filepath.Walk(target, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && (info.Name() == ".git" || info.Name() == ".vscode") {
			return filepath.SkipDir
		}
		paths = append(paths, filepath.ToSlash(path))
		return nil
	})
	return paths, err
}

func processPath(path string, addr *urn.URN) error {
	if slices.Contains(archiveTypes, filepath.Ext(path)) {
		return processArchiveFile(path, addr)
	}

	category, matched, readFlag := scanner.CheckRelatedFiles(path)
	if matched {
		return processMatchedFile(path, category, readFlag, addr)
	}
	return nil
}

func processArchiveFile(path string, addr *urn.URN) error {
	reader, err := os.Open(path)
	if err != nil {
		return err
	}
	defer reader.Close()

	stat, err := reader.Stat()
	if err != nil {
		return err
	}

	b, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	processArchive(bytes.NewReader(b), path, stat.Size(), addr)
	return nil
}

func processMatchedFile(path, category string, readFlag bool, addr *urn.URN) error {
	switch category {
	case "rpm":
		return handleRpmFile(path, category, "", addr)
	default:
		if !readFlag {
			stream.Emit(category, types.Payload{
				Address: addr,
				Body:    path,
			})
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		return handleManifestFile(path, category, "", file, addr)
	}
}

func handleRpmFile(path, category, layer string, addr *urn.URN) error {
	rpmDb := types.RpmDB{
		Path: path,
	}

	err := rpmDb.ReadDBFile(path)
	if err != nil {
		return err
	}
	stream.Emit(category, types.Payload{
		Address: addr,
		Layer:   layer,
		Body:    rpmDb,
	})
	return nil
}

func handleManifestFile(path, category, layer string, file *os.File, addr *urn.URN) error {
	manifest := types.ManifestFile{
		Path: path,
	}
	err := manifest.ReadContent(file)
	if err != nil {
		return err
	}

	stream.Emit(category, types.Payload{
		Address: addr,
		Layer:   layer,
		Body:    manifest,
	})

	return nil
}
