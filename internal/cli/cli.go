package cli

import (
	"time"

	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/internal/presenter"
	"github.com/carbonetes/diggity/internal/presenter/status"
	"github.com/carbonetes/diggity/pkg/config"
	"github.com/carbonetes/diggity/pkg/reader"
	"github.com/carbonetes/diggity/pkg/types"
)

func Start(parameters types.Parameters) {
	start := time.Now()
	if !parameters.Quiet {
		status.Run()
	} else {
		if parameters.OutputFormat == types.Table {
			parameters.OutputFormat = types.JSON
		}
	}

	switch parameters.ScanType {
	case 1: // Image Scan Type
		image, err := reader.GetImage(parameters.Input, &config.Config.Registry)
		if err != nil {
			log.Error(err)
		}
		err = reader.ReadFiles(image)
		if err != nil {
			log.Error(err)
		}
	case 2: // Tarball Scan Type
		image, err := reader.ReadTarballAsImage(parameters.Input)
		if err != nil {
			log.Error(err)
		}
		err = reader.ReadFiles(image)
		if err != nil {
			log.Error(err)
		}
	case 3: // Filesystem Scan Type
		reader.FilesystemScanHandler(parameters.Input)
	default:
		log.Error("Unknown scan type")
	}

	presenter.DisplayResults(parameters, time.Since(start).Seconds())

}