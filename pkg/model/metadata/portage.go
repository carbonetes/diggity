package metadata

// PortageMetadata portage metadata
type PortageMetadata struct {
	Size  int           `json:"size,omitempty"`
	Files []PortageFile `json:"files,omitempty"`
}

// PortageFile file metadata
type PortageFile struct {
	Path   string        `json:"path,omitempty"`
	Digest PortageDigest `json:"digest,omitempty"`
}

// PortageDigest file digest metadata
type PortageDigest struct {
	Algorithm string `json:"algorithm,omitempty"`
	Value     string `json:"value,omitempty"`
}
