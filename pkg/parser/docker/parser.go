package docker

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/carbonetes/diggity/pkg/model"
)

func parseImageInfo(files []string) (*[]model.DockerManifest, *model.DockerConfig, error) {
	dockerManifests := new([]model.DockerManifest)
	dockerConfig := new(model.DockerConfig)
	for _, f := range files {
		file, err := os.Open(f)
		if err != nil {
			return nil, nil, err
		}
		defer file.Close()
		parser := json.NewDecoder(file)
		if strings.Contains(file.Name(), "manifest") {
			var manifests []model.DockerManifest
			if err := parser.Decode(&manifests); err != nil {
				return nil, nil, err
			}
			*dockerManifests = append(*dockerManifests, manifests...)
			continue
		}
		if err := parser.Decode(&dockerConfig); err != nil {
			return nil, nil, err
		}
	}

	return dockerManifests, dockerConfig, nil
}
