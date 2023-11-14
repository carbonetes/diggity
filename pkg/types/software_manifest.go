package types

type SoftwareManifest struct {
	SBOM       interface{} `json:"sbom"`
	Distro     Distro      `json:"distro,omitempty"`
	ImageInfo  ImageInfo   `json:"image_info,omitempty"`
	Secrets    []Secret    `json:"secrets"`
	Parameters Parameters  `json:"parameters"`
}
