package file

import (
	"archive/tar"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/carbonetes/diggity/internal/model"
)

const (
	gzFile = ".gz" // Invalid zip file
)

// Contents contains the location of files
var Contents = make([]*model.Location, 0)

// UnTar extract all files from source into dst (directory)
func UnTar(dst string, source string, recursive bool) error {

	reader, _ := os.Open(source)
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

		// Skip unsafe files fo extraction
		if strings.Contains(target, "..") {
			continue
		}
		if strings.Contains(filepath.Base(target), gzFile) {
			continue
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, fs.ModePerm); err != nil {
					return err
				}
			}

		case tar.TypeReg:
			if err := processFile(tarReader, target, os.FileMode(header.Mode), recursive); err != nil {
				return err
			}
		}
	}
}

// Get the content path of the file, and its children, if any
func processFile(tarReader *tar.Reader, target string, fileMode fs.FileMode, recursive bool) error {
	f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, fileMode)

	if err != nil {
		return err
	}

	if strings.Contains(f.Name(), "layer.tar") && recursive {
		childDar := strings.Replace(f.Name(), "layer.tar", "", -1)
		os.Mkdir(childDar, fs.ModePerm)
		defer UnTar(childDar, f.Name(), true)
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
	Contents = append(Contents, &model.Location{Path: f.Name(), LayerHash: path})
	if err := f.Close(); err != nil {
		panic(err)
	}

	return nil
}
