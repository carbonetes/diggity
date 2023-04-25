package files

import (
	"os"
	"path/filepath"

	"github.com/carbonetes/diggity/pkg/model"
)

func Exists(filename string) bool {
	if _, err := os.Stat(filename); err == nil {
		return true
	}
	return false
}

func GetFilesFromDir(source string) (*[]model.Location, error) {
	contents := new([]model.Location)
	// recursive
	err := filepath.Walk(source,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() && (info.Name() == ".git" || info.Name() == ".vscode") {
				return filepath.SkipDir
			}
			*contents = append(*contents, model.Location{Path: path})
			return nil
		})
	if err != nil {
		return nil, err
	}

	return contents, nil
}
