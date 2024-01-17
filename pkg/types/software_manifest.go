package types

type SoftwareManifest struct {
	SBOM       interface{} `json:"sbom"`
	OS         []OSRelease `json:"os,omitempty"`
	ImageInfo  ImageInfo   `json:"image_info,omitempty"`
	Secrets    []Secret    `json:"secrets,omitempty"`
	Files      []string    `json:"files,omitempty"`
	Parameters Parameters  `json:"parameters"`
	SLSA       SLSA        `json:"slsa,omitempty"`
	Duration   float64     `json:"duration"`
}
