package curator

import (
	"github.com/carbonetes/diggity/pkg/stream"
	"github.com/carbonetes/diggity/pkg/types"
)

func ParametersStoreWatcher(data interface{}) interface{} {
	parameters, ok := data.(types.Parameters)
	if !ok {
		log.Print("ParametersStoreWatcher received unknown type")
	}

	switch parameters.ScanType {
	case 1: // Image Scan Type
		stream.Emit(stream.ImageScanEvent, parameters.Input)
	case 2: // Tarball Scan Type
		stream.Emit(stream.TarballScanEvent, parameters.Input)
	default:
		log.Error("Unknown scan type")
	}
	// result := stream.AggrerateSoftwareManifest()
	stream.Emit(stream.ScanCompleteEvent, true)
	return data
}
