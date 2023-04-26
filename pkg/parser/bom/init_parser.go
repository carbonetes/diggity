package bom

import (
	"errors"
	"log"
	"strings"
	"sync"

	"github.com/carbonetes/diggity/internal/docker"
	client "github.com/carbonetes/diggity/pkg/docker"
	"github.com/carbonetes/diggity/pkg/files"
	"github.com/carbonetes/diggity/pkg/model"
)

type ParserRequirements struct {
	Arguments *model.Arguments
	Dir       *string
	Contents  *[]model.Location
	Result    *model.Result
	WG        sync.WaitGroup
	Errors    *[]error
}

// InitParsers initialize arguments
func InitParsers(args *model.Arguments) (*ParserRequirements, error) {
	if len(*args.Image) > 0 {
		if !strings.Contains(*args.Image, ":") {
			modifiedTag := *args.Image + ":latest"
			args.Image = &modifiedTag
		}
		credential := model.NewRegistryAuth(args)
		imageId := client.GetImageID(args.Image, credential)
		contents, dir := client.ExtractImage(imageId)
		docker.CreateTempDir()
		return &ParserRequirements{
			Arguments: args,
			Dir:       dir,
			Contents:  contents,
			Errors:    new([]error),
			Result: &model.Result{
				Packages: new([]model.Package),
				Secret:   new(model.SecretResults),
				Distro:   new(model.Distro),
				SLSA:     new(model.SLSA),
			},
		}, nil
	} else if len(*args.Dir) > 0 {
		if files.Exists(*args.Dir) {
			contents, err := files.GetFilesFromDir(*args.Dir)
			if err != nil {
				log.Fatal(err)
			}
			docker.CreateTempDir()
			return &ParserRequirements{
				Arguments: args,
				Dir:       args.Dir,
				Contents:  contents,
				Errors:    new([]error),
				Result: &model.Result{
					Packages: new([]model.Package),
					Secret:   new(model.SecretResults),
					Distro:   new(model.Distro),
					SLSA:     new(model.SLSA),
				},
			}, nil
		} else {
			log.Fatal(errors.New("Directory not found!"))
		}
	} else if len(*args.Tar) > 0 {
		if files.Exists(*args.Tar) {
			contents, dir := client.ExtractTarFile(args.Tar)
			docker.CreateTempDir()
			return &ParserRequirements{
				Arguments: args,
				Dir:       dir,
				Contents:  contents,
				Errors:    new([]error),
				Result: &model.Result{
					Packages: new([]model.Package),
					Secret:   new(model.SecretResults),
					Distro:   new(model.Distro),
					SLSA:     new(model.SLSA),
				},
			}, nil
		} else {
			return nil, errors.New("Tar file not found!")
		}
	} else {
		return nil, errors.New("No valid scanning target provided!")
	}
	return nil, nil
}
