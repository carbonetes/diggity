package reader

import (
	"time"

	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/internal/presenter"
	"github.com/carbonetes/diggity/internal/presenter/status"
	"github.com/carbonetes/diggity/pkg/stream"
	"github.com/carbonetes/diggity/pkg/types"
)

// ParametersStoreWatcher watches for changes in the parameters store and emits events based on the scan type.
func ParametersStoreWatcher(data interface{}) interface{} {
	parameters, ok := data.(types.Parameters)
	if !ok {
		log.Print("ParametersStoreWatcher received unknown type")
	}

	if !parameters.Quiet {
		status.Init()
		status.Run()
	} else {
		if parameters.OutputFormat == types.Table {
			parameters.OutputFormat = types.JSON
		}
	}
	stream.Set(stream.ScanStartStoreKey, time.Now())
	switch parameters.ScanType {
	case 1: // Image Scan Type
		stream.Emit(stream.ImageScanEvent, parameters.Input)
	case 2: // Tarball Scan Type
		stream.Emit(stream.TarballScanEvent, parameters.Input)
	case 3: // Filesystem Scan Type
		stream.Emit(stream.FilesystemScanEvent, parameters.Input)
	default:
		log.Error("Unknown scan type")
	}
	elapsed := time.Since(stream.GetScanStart()).Seconds()
	stream.SetScanElapsed(elapsed)
	presenter.DisplayResults(parameters, elapsed)
	return nil
}
