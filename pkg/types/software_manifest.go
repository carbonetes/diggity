package types

type SoftwareManifest struct {
	SBOM          interface{}   `json:"sbom"`
	Distro        Distro        `json:"distro,omitempty"`
	ImageManifest ImageManifest `json:"image_manifest,omitempty"`
	Secret        SecretResult  `json:"secret,omitempty"`
	Parameters    Parameters    `json:"parameters"`
}
