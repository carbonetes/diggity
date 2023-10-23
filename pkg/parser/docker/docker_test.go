package docker_test

import (
	"path/filepath"
	"testing"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/parser/common"
	"github.com/carbonetes/diggity/pkg/parser/docker"
)

const dockerReference string = "diggity-tmp-385abb3c-df38-44dd-b30f-467ba364ee3a"

var (
	args   = model.NewArguments()
	target = filepath.Join("..", "..", "..", "docs", "references", "docker", dockerReference)
)

func TestDocker(t *testing.T) {
	args.Dir = &target
	req, err := common.NewParams(args)
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
