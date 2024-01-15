package reader

import (
	"github.com/carbonetes/diggity/pkg/stream"
)

// Init initializes the stream by watching for changes in the parameters store key and scan elapsed store key.
// It also attaches handlers for image scan event, tarball scan event, and filesystem check event.
func Init() {
	stream.Watch(stream.ParametersStoreKey, ParametersStoreWatcher)
	stream.Attach(stream.ImageScanEvent, ImageScanHandler)
	stream.Attach(stream.TarballScanEvent, TarballScanHandler)
	stream.Attach(stream.FilesystemScanEvent, FilesystemScanHandler)
}
