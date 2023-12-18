package portage

// PortageMetadata portage metadata
type Metadata struct {
	Size  int    `json:"size,omitempty"`
	Files []File `json:"files,omitempty"`
}

// PortageFile file metadata
type File struct {
	Path   string `json:"path,omitempty"`
	Digest Digest `json:"digest,omitempty"`
}

// PortageDigest file digest metadata
type Digest struct {
	Algorithm string `json:"algorithm,omitempty"`
	Value     string `json:"value,omitempty"`
}
