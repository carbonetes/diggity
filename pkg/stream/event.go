package stream

const (
	ParameterAcceptEvent string = "parameter.accept"
	ImageScanEvent       string = "image.scan"
	TarballScanEvent     string = "tarball.scan"
	FilesystemCheckEvent string = "filesystem.check"
	ScanManifestEvent    string = "scan.manifest."
	ComponentFoundEvent  string = "component.found"
	DistroFoundEvent     string = "distro.found"
	ErrorOccurredEvent   string = "error.occurred"

	ComponentScanEvent string = "scan.component."
	DistroScanEvent    string = "scan.distro"
	SecretScanEvent    string = "scan.secret"

)
