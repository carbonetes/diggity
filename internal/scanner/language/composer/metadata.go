package composer

import "encoding/json"

// ComposerMetadata - composer file metadata
type ComposerMetadata struct {
	Readme           []string          `json:"_readme"`
	ContentHash      string            `json:"content-hash"`
	Packages         []ComposerPackage `json:"packages"`
	PackagesDev      []ComposerPackage `json:"packages-dev"`
	Aliases          []string          `json:"aliases"`
	MinimumStability string            `json:"minimum-stability"`
	StabilityFlags   interface{}       `json:"stability-flags"`
	PreferStable     bool              `json:"prefer-stable"`
	PreferLowest     bool              `json:"prefer-lowest"`
	Platform         interface{}       `json:"platform"`
	PlatformDev      interface{}       `json:"platform-dev"`
}

// ComposerPackage - composer packages
type ComposerPackage struct {
	Name            string           `json:"name"`
	Version         string           `json:"version"`
	Source          ComposerObject   `json:"source"`
	Dist            ComposerObject   `json:"dist"`
	Require         ComposerObject   `json:"require"`
	Provide         ComposerObject   `json:"provide"`
	RequireDev      ComposerObject   `json:"require-dev"`
	Suggest         ComposerObject   `json:"suggest"`
	Type            string           `json:"type"`
	Extract         ComposerObject   `json:"extra"`
	Autoload        ComposerObject   `json:"autoload"`
	NotificationURL string           `json:"notification-url"`
	License         []string         `json:"license"`
	Authors         []ComposerObject `json:"authors"`
	Description     string           `json:"description"`
	Homepage        string           `json:"homepage"`
	Keywords        []string         `json:"keywords"`
	Time            string           `json:"time"`
}

// ComposerObject common objects for composer metadata
type ComposerObject map[string]interface{}

func readManifestFile(content []byte) ComposerMetadata {
	var metadata ComposerMetadata
	err := json.Unmarshal(content, &metadata)
	if err != nil {
		log.Error("Failed to unmarshal composer.lock")
	}
	return metadata
}
