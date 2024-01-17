package stream

const (
	ParametersStoreKey    string = "parameters"
	ConfigStoreKey        string = "config"
	ImageInstanceStoreKey string = "image-instance"
	FileListStoreKey      string = "file-list"

	ParameterScanTypeStoreKey         string = "parameter.scan-type"
	ParameterInputStoreKey            string = "parameter.input"
	ParameterOutputFormatStoreKey     string = "parameter.output-format"
	ParameterQuietStoreKey            string = "parameter.quiet"
	ParameterMaxFileSizeStoreKey      string = "parameter.max-file-size"
	ParameterScannersStoreKey         string = "parameter.scanners"
	ParameterAllowFileListingStoreKey string = "parameter.allow-file-listing"

	SBOMStoreKey             string = "software-manifest.sbom"
	ComponentsStoreKey       string = "software-manifest.sbom.components"
	DistroStoreKey           string = "software-manifest.distro"
	ImageInfoStoreKey        string = "software-manifest.image-manifest"
	SecretParametersStoreKey string = "software-manifest.secret.parameters"
	SecretsStoreKey          string = "software-manifest.secrets"
	OSReleasesStoreKey       string = "software-manifest.os-releases"

	ScanStartStoreKey   string = "scan.start"
	ScanElapsedStoreKey string = "scan.elapsed"

	CycloneDXComponentsStoreKey string = "cyclonedx.components"
	CycloneDXBOMStoreKey        string = "cyclonedx.bom"
)
