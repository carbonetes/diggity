package model

type SBOM struct {
	Packages  *[]Package     `json:"packages"`
	Secret    *SecretResults `json:"secrets,omitempty"`
	ImageInfo ImageInfo      `json:"imageInfo"`
	Distro    *Distro        `json:"distro"`
	SLSA      *SLSA          `json:"slsa,omitempty"`
}