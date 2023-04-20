package docker

import (
	"context"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/carbonetes/diggity/internal/file"
	"github.com/docker/docker/pkg/ioutils"
	"github.com/google/uuid"
)

// ExtractImage extracts a Docker image to a temporary directory and returns the path to the directory.
func ExtractImage(target *string) *string {
	ids := new([]string)
	*ids = append(*ids, *target)

	// Get a reader for the saved Docker image.
	reader, err := docker.ImageSave(context.Background(), *ids)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	// Create a temporary directory to extract the Docker image to.
	tempDir, err := ioutils.TempDir("", "")
	if err != nil {
		log.Fatal(err)
	}
	tarFileName := "diggity-tmp-" + uuid.NewString() + ".tar"
	tarPath := filepath.Join(tempDir, tarFileName)
	tarFile, err := os.Create(tarPath)
	if err != nil {
		log.Fatal(err)
	}

	// Create a directory to extract the Docker image to.
	extractDir := strings.Replace(tarFile.Name(), ".tar", "", -1)
	err = os.Mkdir(extractDir, fs.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	// Copy the Docker image to the temporary file.
	_, err = io.Copy(tarFile, reader)
	if err != nil {
		log.Fatal(err)
	}

	// Extract the Docker image to the temporary directory.
	if err = file.UnTar(extractDir, tarFile.Name(), true); err != nil {
		log.Fatal(err)
	}

	return &extractDir
}

func ExtractTarFile(tar *string) *string {
	dir, err := ioutils.TempDir("", "")
	if err != nil {
		log.Fatal(err.Error())
	}

	folder := "diggity-tmp-dir" + uuid.NewString()
	target := filepath.Join(dir, folder)
	err = os.Mkdir(target, fs.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	if err := file.UnTar(target, *tar, true); err != nil {
		log.Fatal(err)
	}
	return &target
}
