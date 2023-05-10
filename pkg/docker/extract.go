package docker

import (
	"archive/tar"
	"context"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/google/uuid"
)

const (
	gzFile           = ".gz"                 // Invalid zip file
	invalidCharRegex = `[,@<>'"|?*#%&{}$=!]` // Invalid filename characters
)

// ExtractImage extracts a Docker image to a temporary directory and returns the path to the directory.
func ExtractImage(target *string, dockerTemp *string) (*[]model.Location, *string) {
	contents := new([]model.Location)
	tarFile := SaveImageToTar(target, dockerTemp)
	defer tarFile.Close()
	// Create a directory to extract the Docker image to.
	extractDir := strings.Replace(tarFile.Name(), ".tar", "", -1)
	err := os.Mkdir(extractDir, fs.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	// Extract the Docker image to the temporary directory.
	if err = UnTar(extractDir, tarFile.Name(), true, contents); err != nil {
		log.Fatal(err)
	}

	return contents, &extractDir
}

func SaveImageToTar(image *string, dockerTemp *string) *os.File {
	ids := new([]string)
	*ids = append(*ids, *image)

	// Get a reader for the saved Docker image.
	reader, err := docker.ImageSave(context.Background(), *ids)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	// Create a temporary directory to extract the Docker image to.
	tarFileName := "diggity-tmp-" + uuid.NewString() + ".tar"
	tarPath := filepath.Join(*dockerTemp, tarFileName)
	tarFile, err := os.Create(tarPath)
	if err != nil {
		log.Fatal(err)
	}
	defer tarFile.Close()

	// Copy the Docker image to the temporary file.
	_, err = io.Copy(tarFile, reader)
	if err != nil {
		log.Fatal(err)
	}

	return tarFile
}

func UnTar(dst string, source string, recursive bool, contents *[]model.Location) error {
	r := regexp.MustCompile(invalidCharRegex)
	reader, err := os.Open(source)
	if err != nil {
		return err
	}
	defer reader.Close()
	tarReader := tar.NewReader(reader)

	for {
		header, err := tarReader.Next()
		switch {

		case err == io.EOF:
			return nil

		case err != nil:
			return err

		case header == nil:
			continue
		}

		target := filepath.Join(dst, header.Name)

		// Skip unsafe files for extraction
		if strings.Contains(filepath.Base(target), gzFile) ||
			r.MatchString(filepath.Base(target)) ||
			strings.Contains(target, "..") {
			continue
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				base := filepath.Base(target)
				if strings.Contains(base, ":") {
					newBase := strings.ReplaceAll(base, ":", "-")
					target = strings.ReplaceAll(target, base, newBase)
				}
				if err := os.MkdirAll(target, fs.ModePerm); err != nil {
					return err
				}
			}

		case tar.TypeReg:
			if err := processFile(tarReader, target, os.FileMode(header.Mode), recursive, contents); err != nil {
				return err
			}
		}
	}
}

func processFile(tarReader *tar.Reader, target string, fileMode fs.FileMode, recursive bool, contents *[]model.Location) error {
	f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, fileMode)
	if err != nil {
		//skipping when file cannot be opened
		return nil // replace and use log to debug
	}

	defer f.Close()

	if strings.Contains(f.Name(), "layer.tar") && recursive {
		childDar := strings.Replace(f.Name(), "layer.tar", "", -1)
		_ = os.Mkdir(childDar, fs.ModePerm)

		defer func() {
			_ = UnTar(childDar, f.Name(), true, contents)
		}()
	}
	_, err = io.Copy(f, tarReader)
	if err != nil {
		return err
	}

	paths := strings.Split(f.Name(), string(os.PathSeparator))

	// Layer SHA Regex
	regex := regexp.MustCompile(`\b[A-Fa-f0-9]{64}\b`)
	var path string
	for _, _path := range paths {
		path = regex.FindString(_path)
		if len(path) > 0 {
			break
		}
	}

	*contents = append(*contents, model.Location{Path: f.Name(), LayerHash: path})
	if err := f.Close(); err != nil {
		return err
	}

	return nil
}

func ExtractTarFile(tar *string, dockerTemp *string) (*[]model.Location, *string) {
	contents := new([]model.Location)

	folder := "diggity-tmp-dir" + uuid.NewString()
	target := filepath.Join(*dockerTemp, folder)
	err := os.Mkdir(target, fs.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	if err := UnTar(target, *tar, true, contents); err != nil {
		log.Fatal(err)
	}
	return contents, &target
}
