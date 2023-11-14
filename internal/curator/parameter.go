package curator

import (
	"fmt"

	"github.com/carbonetes/diggity/internal/helper"
	"github.com/carbonetes/diggity/pkg/stream"
	"github.com/carbonetes/diggity/pkg/types"
)

func ParametersStoreWatcher(data interface{}) interface{} {
	parameters, ok := data.(types.Parameters)
	fmt.Println(parameters)
	if !ok {
		log.Print("ParametersStoreWatcher received unknown type")
	}

	switch parameters.ScanType {
	case 1: // Image Scan Type
		stream.Emit(stream.ImageScanEvent, parameters.Input)
	case 2: // Tarball Scan Type
		stream.Emit(stream.TarballScanEvent, parameters.Input)
	}
	result := stream.AggrerateSoftwareManifest()
	json, err := helper.ToJSON(result)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(json))
	return data
}
