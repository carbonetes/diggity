package scanner_test

import (
	"path/filepath"
	"testing"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/parser/common"
	"github.com/carbonetes/diggity/pkg/scanner"
	"github.com/stretchr/testify/assert"
)

var (
	testImage     string = "alpine"
	testDirectory string = filepath.Join("..", "..", "docs", "references")
	enable        bool   = true
)

func TestScan(t *testing.T) {
	t.Run("Scanning image", func(t *testing.T) {
		arguments := model.NewArguments()
		arguments.Image = &testImage
		arguments.Quiet = &enable
		req, err := common.NewParams(arguments)
		assert.Empty(t, err)
		sbom, errs := scanner.Scan(req.Arguments)
		if len(*errs) > 0 {
			for _, err := range *errs {
				t.Log(err)
			}
		}
		assert.Empty(t, errs)
		assert.NotEmpty(t, sbom.Packages)
		assert.NotEmpty(t, sbom.Distro)
		assert.NotEmpty(t, sbom.ImageInfo)
		assert.NotEmpty(t, sbom.ImageInfo.DockerConfig)
		assert.NotEmpty(t, sbom.ImageInfo.DockerManifest)
	})

	t.Run("Scanning directories", func(t *testing.T) {
		arguments := model.NewArguments()
		arguments.Dir = &testDirectory
		arguments.Quiet = &enable
		req, err := common.NewParams(arguments)
		assert.Nil(t, err)
		sbom, errs := scanner.Scan(req.Arguments)
		assert.Empty(t, errs)
		assert.NotEmpty(t, sbom.Packages)
	})

}
