package docker_test

import (
	"path/filepath"
	"testing"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/parser/bom"
	"github.com/carbonetes/diggity/pkg/parser/docker"
)

var (
	args   = model.NewArguments()
	target = filepath.Join("..", "..", "..", "docs", "references", "docker")
)

func TestDocker(t *testing.T) {
	args.Dir = &target
	req, err := bom.InitParsers(args)
	if err != nil {
		t.Fatal(err)
		t.FailNow()
	}

	req.WG.Add(1)
	docker.ParseDockerProperties(req)
	req.WG.Wait()

	imageInfo := req.SBOM.ImageInfo

	if imageInfo.DockerConfig.Created == "" && imageInfo.DockerConfig.OS == "" && imageInfo.DockerConfig.Architecture == "" {
		t.Error("Docker Configuration is empty or incomplete.")
	}

	if len(imageInfo.DockerManifest) == 0 {
		t.Error("Docker Manifest is empty.")
	}

}
