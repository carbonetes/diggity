package file

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/carbonetes/diggity/pkg/model"
)

// Exists checks if filename exists
func Exists(filename string) bool {
	if _, err := os.Stat(filename); err == nil {
		return true
	}
	return false
}

// CheckUserInput evaluates specified arguments by user
func CheckUserInput(args *model.Arguments) (string, string) {
	if args.Image == nil && len(*args.Dir) == 0 && len(*args.Tar) > 0 {
		return "tar", "Extracting Image tar File..."
	} else if args.Image == nil && len(*args.Dir) > 0 && len(*args.Tar) == 0 {
		return "dir", "Checking File Directory"
	} else if args.Image != nil && len(*args.Dir) == 0 && len(*args.Tar) == 0 {
		//do double checking if image input is valid
		if strings.HasSuffix(*args.Image, ".tar") {
			args.Tar = args.Image
			return "tar", "Extracting Image tar File..."
		}

		dir := *args.Image
		if Exists(dir) {
			args.Dir = args.Image
			return "dir", "Checking File Directory"
		}

		return "image", "Checking image from local..."
	}

	return "", ""
}

// GetFilesFromDir parses files from a dir
func GetFilesFromDir(source string) error {
	// recursive
	err := filepath.Walk(source,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			//exclude hiddent folder for git and vs code if present
			if info.IsDir() && (info.Name() == ".git" || info.Name() == ".vscode") {
				return filepath.SkipDir
			}
			Contents = append(Contents, &model.Location{Path: path})
			return nil
		})
	if err != nil {
		return err
	}

	return nil
}
