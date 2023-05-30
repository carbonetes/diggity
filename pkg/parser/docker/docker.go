package docker

import (
	"errors"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/parser/bom"
)

const parserErr = "docker-parser: "

// ParseDockerProperties appends docker json files to parser.Result
func ParseDockerProperties(req *bom.ParserRequirements) {
	var imageInfo model.ImageInfo
	tarDirectory, err := getTarDir(*req.Dir, *req.Arguments.Dir)
	if err != nil {
		err = errors.New(parserErr + err.Error())
		*req.Errors = append(*req.Errors, err)
	}

	files, err := getJSONFilesFromDir(tarDirectory.Name())
	if err != nil {
		err = errors.New(parserErr + err.Error())
		*req.Errors = append(*req.Errors, err)
	}

	err = tarDirectory.Close()
	if err != nil {
		err = errors.New(parserErr + err.Error())
		*req.Errors = append(*req.Errors, err)
	}

	manifests, config, err := parseImageInfo(files)
	if err != nil {
		err = errors.New(parserErr + err.Error())
		*req.Errors = append(*req.Errors, err)
	}

	if manifests != nil {
		imageInfo.DockerManifest = append(imageInfo.DockerManifest, *manifests...)
	}

	if config != nil {
		imageInfo.DockerConfig = *config
	}

	req.SBOM.ImageInfo = imageInfo

	defer req.WG.Done()
}
