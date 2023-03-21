package metadata

type HexMetadata struct {
	Name       string `json:"name"`
	Version    string `json:"version"`
	PkgHash    string `json:"pkgHash,omitempty"`
	PkgHashExt string `json:"pkgHashExt,omitempty"`
}
