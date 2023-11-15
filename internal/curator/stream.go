package curator

import (
	"github.com/carbonetes/diggity/internal/presenter"
	"github.com/carbonetes/diggity/internal/presenter/status"
	"github.com/carbonetes/diggity/pkg/stream"
)

func Init() {
	stream.Attach(stream.ParametersStoreKey, ParametersStoreWatcher)
	stream.Attach(stream.ImageScanEvent, IndexImageFilesystem)
	stream.Attach(stream.TarballScanEvent, IndexTarballFilesystem)
	stream.Attach(stream.FilesystemCheckEvent, status.ScanFile)
	stream.Watch(stream.ScanElapsedStoreKey, status.ScanCompleteStatus)
	stream.Watch(stream.ScanElapsedStoreKey, presenter.DisplayResults)
}
