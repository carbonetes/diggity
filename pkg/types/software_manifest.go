package types

type SoftwareManifest struct {
	SBOM       interface{} `json:"sbom"`
	Distro     Distro      `json:"distro,omitempty"`
	ImageInfo  ImageInfo   `json:"image_info,omitempty"`
	Secrets    []Secret    `json:"secrets"`
	Files      []string    `json:"files"`
	Parameters Parameters  `json:"parameters"`
	Duration   float64     `json:"duration"`
}
