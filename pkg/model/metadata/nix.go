package metadata

type NixMetadata struct {
	Hash       string `json:"hash"`
	Name       string `json:"name"`
	Version    string `json:"version"`
	Prerelease string `json:"prerelease"`
}
