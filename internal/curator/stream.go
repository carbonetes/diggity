package curator

import (
	"github.com/carbonetes/diggity/pkg/stream"
)

func Init() {
	stream.Attach(stream.ParametersStoreKey, ParametersStoreWatcher)
	stream.Attach(stream.ImageScanEvent, IndexImageFilesystem)
	stream.Attach(stream.TarballScanEvent, IndexTarballFilesystem)
}
