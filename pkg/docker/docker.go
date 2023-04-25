package docker

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/carbonetes/diggity/internal/logger"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/ioutils"
	"golang.org/x/exp/slices"
)

var (
	log     = logger.GetLogger()
	docker  *client.Client
	timeout int = 300 // Timeout in seconds for pulling an image from a registry.
)

// init is called before the package is initialized. It sets up a new Docker client.
func init() {
	client, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Fatal(err)
	}
	docker = client
}

// GetImageID returns a Docker image ID given an image target (name:tag) and authentication info
// (if required). If the image is available locally, it returns the local image ID. Otherwise, it
// attempts to pull the image from a registry using the provided credentials. If credentials are not
// provided, it attempts to pull a public image.
func GetImageID(target *string, credential *types.AuthConfig) *string {
	imageId := FindImageFromLocal(target)
	if imageId != nil {
		return imageId
	}

	if credential != nil {
		return PullImageFromRegistry(target, credential)
	} else {
		return PullPublicImage(target)
	}
}

// FindImageFromLocal searches for a Docker image of the given target on the local system.
// If found, it returns its ID; otherwise it returns nil.
func FindImageFromLocal(target *string) *string {
	images, err := docker.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	for _, image := range images {
		if len(image.RepoTags) > 0 {
			if slices.Contains(image.RepoTags, *target) {
				return &image.ID
			}
		}
	}

	return nil
}

// PullPublicImage attempts to pull a public Docker image of the given target from Docker Hub.
// If successful, it returns the image ID; otherwise it waits for timeout seconds and retries.
func PullPublicImage(target *string) *string {
	reader, err := docker.ImagePull(context.Background(), *target, types.ImagePullOptions{})

	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	timer := 0
	for timer <= timeout {
		if FindImageFromLocal(target) != nil {
			return FindImageFromLocal(target)
		}
		time.Sleep(1 * time.Second)
		timer++
	}

	if timer >= timeout {
		log.Fatalf("Connection Timed Out: Pulling image %+v is taking some time. Consider pulling image first before scanning.", target)
	}

	return FindImageFromLocal(target)
}

// PullImageFromRegistry attempts to pull a Docker image of the given target from a registry,
// using the provided authentication credentials. If successful, it returns the image ID;
// otherwise it waits for timeout seconds and retries.
func PullImageFromRegistry(target *string, credential *types.AuthConfig) *string {
	data, err := json.Marshal(credential)
	if err != nil {
		log.Fatal(err)
	}
	auth := base64.URLEncoding.EncodeToString(data)
	reader, err := docker.ImagePull(context.Background(), *target, types.ImagePullOptions{
		RegistryAuth: auth,
	})
	if err != nil {
		log.Fatal(err)
	}

	defer reader.Close()

	timer := 0
	for timer <= timeout {
		if FindImageFromLocal(target) != nil {
			return FindImageFromLocal(target)
		}
		time.Sleep(1 * time.Second)
		timer++
	}

	if timer >= timeout {
		log.Fatalf("Connection Timed Out: Pulling image %+v is taking some time. Consider pulling image first before scanning.", target)
	}

	return FindImageFromLocal(target)
}

func CreateTempDir() *string {
	tempDir, err := ioutils.TempDir("", "")
	if err != nil {
		log.Fatal(err)
	}
	return &tempDir
}
