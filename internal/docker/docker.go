package docker

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/carbonetes/diggity/internal/file"
	"github.com/carbonetes/diggity/internal/logger"
	"github.com/carbonetes/diggity/internal/ui"
	"github.com/carbonetes/diggity/pkg/model"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/ioutils"
	"github.com/google/uuid"
	"github.com/schollz/progressbar/v3"
)

var (
	dockerClient *client.Client
	arguments    *model.Arguments
	imageIDS     []string
	tempDir      string
	extractDir   string
	log          = logger.GetLogger()
)

const (
	timeout       int = 60 //timeout in seconds
	imageIDLength int = 12 //short 12-digit hex string of the Image ID
)

// Dir returns generated temporary directory string
func Dir() string {
	return tempDir
}

// ExtractedDir returns the extracted directory of tar file
func ExtractedDir() string {
	return extractDir
}

// ExtractImage extracts docker image contents
func ExtractImage(_arguments *model.Arguments, spinner *progressbar.ProgressBar) {
	arguments = _arguments

	if err := testConnection(); err != nil {
		log.Fatal("\nCannot connect to the Docker daemon at unix:///var/run/docker.sock. Is the docker daemon running?\n")
		os.Exit(1)
	}

	// Get Image ID; Pull from server if needed
	imageIDS = append(imageIDS, getImageID(_arguments))

	// Run Spinner
	if !*arguments.Quiet {
		go ui.RunSpinner(spinner)
		defer ui.DoneSpinner(spinner)
	}

	// Extract Image
	reader, err := dockerClient.ImageSave(context.Background(), imageIDS)
	if err != nil {
		panic(err)
	}
	defer reader.Close()

	tempDir, err = ioutils.TempDir("", "")
	if err != nil {
		panic(err)
	}
	tarFileName := "diggity-tmp-" + uuid.NewString() + ".tar"
	tarPath := filepath.Join(tempDir, tarFileName)
	tarFile, err := os.Create(tarPath)
	if err != nil {
		log.Fatal(err.Error())
	}

	extractDir = strings.Replace(tarFile.Name(), ".tar", "", -1)
	err = os.Mkdir(extractDir, fs.ModePerm)
	if err != nil {
		panic(err)
	}

	_, err = io.Copy(tarFile, reader)
	if err != nil {
		panic(err)
	}

	if err = file.UnTar(extractDir, tarFile.Name(), true); err != nil {
		log.Fatal(err.Error())
	}
}

func init() {
	dockerClient, _ = client.NewClientWithOpts(client.FromEnv)
}

// Test docker client connection
func testConnection() error {
	_, err := dockerClient.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		return err
	}
	return nil
}

// Get image ID
func getImageID(_arguments *model.Arguments) string {
	// Check Image From Local
	imageID := findImageID()

	// If image not in local, pull from server
	if len(imageID) == 0 {
		imageID = pullImageFromServer(_arguments)
	}

	return imageID
}

// Find docker image from local images
func findImageID() (imageID string) {
	images, err := dockerClient.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		log.Fatalf("\nImage %s not found\n", *arguments.Image)
	}

	for _, image := range images {
		// Check if arg is Image ID
		idArg := strings.Split(*arguments.Image, ":")[0]
		if len(idArg) == imageIDLength && strings.Contains(image.ID, idArg) {
			return image.ID
		}

		// Check by Repo Tags
		if len(image.RepoTags) > 1 {
			for _, repoTag := range image.RepoTags {
				if repoTag == *arguments.Image {
					imageID = image.ID
				}
			}
		} else if len(image.RepoTags) == 1 {
			if image.RepoTags[0] == *arguments.Image {
				return image.ID
			}
		}
	}

	return imageID
}

// Pull docker image from server
func pullImageFromServer(_arguments *model.Arguments) string {

	pullingSpinner := ui.InitSpinner("Pulling image from server")

	// Run Spinner
	if !*arguments.Quiet {
		go ui.RunSpinner(pullingSpinner)
		defer ui.DoneSpinner(pullingSpinner)
	}

	// Verify if image is found on server
	reader, err := pullImage(_arguments)
	if err != nil {
		log.Fatalf("\nImage %s not found\n", *arguments.Image)
	}
	defer reader.Close()

	// Wait til image is pulled
	timer := 0
	for timer <= timeout {
		// Check if image is already in image list
		if len(findImageID()) != 0 {
			ui.DoneSpinner(pullingSpinner)
			return findImageID()
		}
		time.Sleep(1 * time.Second)
		timer++
	}
	// Timeout Error
	if timer >= timeout {
		log.Fatalf("Connection Timed Out: Pulling image %+v is taking some time. Consider pulling image first before scanning.", *arguments.Image)
	}

	return ""
}

// Pull image from registry
func pullImage(_arguments *model.Arguments) (io.ReadCloser, error) {
	if hasUserNameAndPassword(_arguments) {
		loginRegistry(_arguments)
		c := &credentials{
			Username: *_arguments.RegistryUsername,
			Password: *_arguments.RegistryPassword,
		}
		encodedJSON, _ := json.Marshal(c)
		authStr := base64.URLEncoding.EncodeToString(encodedJSON)
		return dockerClient.ImagePull(context.Background(), *arguments.Image, types.ImagePullOptions{
			RegistryAuth: authStr,
		})
	}

	return dockerClient.ImagePull(context.Background(), *arguments.Image, types.ImagePullOptions{})
}

// Validate username and password
func hasUserNameAndPassword(_arguments *model.Arguments) bool {
	return (*_arguments.RegistryUsername != "" && *_arguments.RegistryPassword != "") || (*_arguments.RegistryUsername != "" && *_arguments.RegistryToken != "")
}

// NewDockerClient (unimplemented)
func NewDockerClient() {
	log.Fatal("unimplemented")
}

// ExtractFromDir unTars file from arguments input
func ExtractFromDir(source *string) {
	tempDir, _ = ioutils.TempDir("", "")
	tarFileFolderName := "diggity-tmp-dir" + uuid.NewString()
	extractDir = filepath.Join(tempDir, tarFileFolderName)
	err := os.Mkdir(extractDir, fs.ModePerm)
	if err != nil {
		panic(err)
	}

	if err := file.UnTar(extractDir, *source, true); err != nil {
		log.Fatal(err.Error())
	}

}

// CreateTempDir creates temp DIR for java parser when arguments is directory
func CreateTempDir() {
	tempDir, _ = ioutils.TempDir("", "")
}
