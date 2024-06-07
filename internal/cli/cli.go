package cli

import (
	"time"

	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/internal/presenter"
	"github.com/carbonetes/diggity/internal/presenter/status"
	"github.com/carbonetes/diggity/pkg/cdx"
	"github.com/carbonetes/diggity/pkg/cdx/dependency"
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

	// Generate unique address for the scan
	addr, err := types.NewAddress(parameters.Input)
	if err != nil {
		log.Fatal(err)
	}

	cdx.New(addr)
	dependency.NewDependencyNodes(addr)
	switch parameters.ScanType {
	case 1: // Image Scan Type
		target := parameters.Input
		image, ref, err := reader.GetImage(target, &config.Config.Registry)
		if err != nil {
			log.Fatal(err)
		}

		cdx.SetMetadataComponent(addr, cdx.SetImageMetadata(*image, *ref, target))

		err = reader.ReadFiles(image, addr)
		if err != nil {
			log.Fatal(err)
		}
	case 2: // Tarball Scan Type
		image, err := reader.ReadTarballAsImage(parameters.Input)
		if err != nil {
			log.Fatal(err)
		}
		err = reader.ReadFiles(image, addr)
		if err != nil {
			log.Fatal(err)
		}
	case 3: // Filesystem Scan Type
		err := reader.FilesystemScanHandler(parameters.Input, addr)
		if err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatal("Unknown scan type")
	}

	presenter.DisplayResults(parameters, time.Since(start).Seconds(), addr)

}
