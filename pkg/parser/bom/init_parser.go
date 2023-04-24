package bom

import (
	"errors"
	"strings"
	"sync"

	"github.com/carbonetes/diggity/internal/docker"
	"github.com/carbonetes/diggity/internal/file"
	"github.com/carbonetes/diggity/internal/logger"
	client "github.com/carbonetes/diggity/pkg/docker"
	"github.com/carbonetes/diggity/pkg/model"
)

var (
	Target = new(string)
	// Arguments - CLI Arguments
	Arguments = model.NewArguments()
	// Packages - common collection of packages found by parsers
	Packages []*model.Package
	// WG - common waitgroup for all the parsers
	WG sync.WaitGroup
	// Errors - common errors encountered by parsers
	Errors []*error
	log    = logger.GetLogger()
)

// InitParsers initialize arguments
func InitParsers(argument model.Arguments) {
	Arguments = &argument
	if len(*Arguments.Image) > 0 {
		if !strings.Contains(*Arguments.Image, ":") {
			modifiedTag := *Arguments.Image + ":latest"
			Arguments.Image = &modifiedTag
		}
		credential := model.NewRegistryAuth(Arguments)
		imageId := client.GetImageID(Arguments.Image, credential)
		Target = client.ExtractImage(imageId)
	} else if len(*Arguments.Dir) > 0 {
		if file.Exists(*Arguments.Dir) {
			err := file.GetFilesFromDir(*Arguments.Dir)
			if err != nil {
				log.Fatal(err)
			}
			docker.CreateTempDir()
		} else {
			log.Fatal(errors.New("Directory not found!"))
		}
	} else if len(*Arguments.Tar) > 0 {
		if file.Exists(*Arguments.Tar) {
			Target = client.ExtractTarFile(Arguments.Tar)
		} else {
			log.Fatal(errors.New("Tar file not found!"))
		}
	} else {
		log.Fatal(errors.New("No valid scanning target provided!"))
	}
}
