package curator

import (
	"time"

	"github.com/carbonetes/diggity/pkg/stream"
	"github.com/carbonetes/diggity/pkg/types"
)

// ParametersStoreWatcher watches for changes in the parameters store and emits events based on the scan type.
func ParametersStoreWatcher(data interface{}) interface{} {
	parameters, ok := data.(types.Parameters)
	if !ok {
		log.Print("ParametersStoreWatcher received unknown type")
	}
	stream.Set(stream.ScanStartStoreKey, time.Now())
	switch parameters.ScanType {
	case 1: // Image Scan Type
		stream.Emit(stream.ImageScanEvent, parameters.Input)
	case 2: // Tarball Scan Type
		stream.Emit(stream.TarballScanEvent, parameters.Input)
	default:
		log.Error("Unknown scan type")
	}

	stream.SetScanElapsed(time.Since(stream.GetScanStart()).Seconds())
	return data
}
