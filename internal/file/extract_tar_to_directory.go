package file

import (
	"archive/tar"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/carbonetes/diggity/pkg/model"
)

const (
	gzFile           = ".gz"                  // Invalid zip file
	invalidCharRegex = `[,@<>:'"|?*#%&{}$=!]` // Invalid filename characters
)

// Contents contains the location of files
var Contents = make([]*model.Location, 0)

// UnTar extract all files from source into dst (directory)
func UnTar(dst string, source string, recursive bool) error { 
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
				if err := os.MkdirAll(target, fs.ModePerm); err != nil {
					return err
				}
			}

		case tar.TypeReg:
			if strings.Contains(target, "layer.tar") && recursive {
				childDir := strings.TrimSuffix(target, "layer.tar")
				if err := os.Mkdir(childDir, fs.ModePerm); err != nil {
					return err
				}

				if err := processNestedTar(childDir, tarReader); err != nil {
					return err
				}
			} else {
				if err := processFile(target, tarReader, os.FileMode(header.Mode), recursive); err != nil {
					return err
				}
			}
		}
	}
}

// Process nested tar file
func processNestedTar(childDir string, parentReader *tar.Reader) error {
	for {
		header, err := parentReader.Next()
		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			return err
		case header == nil:
			continue
		}

		target := filepath.Join(childDir, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, fs.ModePerm); err != nil {
					return err
				}
			}

		case tar.TypeReg:
			if err := processFile(target, parentReader, os.FileMode(header.Mode), true); err != nil {
				return err
			}
		}
	}
}

// Get the content path of the file, and its children, if any
func processFile(target string, tarReader *tar.Reader, fileMode fs.FileMode, recursive bool) error {
	f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, fileMode)
	if err != nil {
		// Skip incorrect names
		if strings.Contains(err.Error(), "The filename, directory name, or volume label syntax is incorrect.") {
			return nil
		}
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, tarReader)
	if err != nil {
		return err
	}

	paths := strings.Split(target, string(os.PathSeparator))

	// Layer SHA Regex
	regex := regexp.MustCompile(`\b[A-Fa-f0-9]{64}\b`)
	var path string
	for _, _path := range paths {
		path = regex.FindString(_path)
		if len(path) > 0 {
			break
		}
	}

	Contents = append(Contents, &model.Location{Path: target, LayerHash: path})

	return nil
}
