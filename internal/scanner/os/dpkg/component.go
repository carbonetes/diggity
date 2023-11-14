package dpkg

import (
	"github.com/carbonetes/diggity/pkg/types"
	"github.com/google/uuid"
)

func newComponent(metadata Metadata) types.Component {
	if metadata["package"] == nil || metadata["version"] == nil {
		return types.Component{}
	}

	var desc string
	if val, ok := metadata["description"].(string); ok {
		desc = val
	}

	return types.Component{
		ID:          uuid.NewString(),
		Name:        metadata["package"].(string),
		Version:     metadata["version"].(string),
		Type:        Type,
		PURL:        setPurl(metadata["package"].(string), metadata["version"].(string), metadata["architecture"].(string)),
		Description: desc,
		Metadata:    metadata,
	}
}

// Parse PURL
func setPurl(name, version, architecture string) string {
	return "pkg" + ":" + "deb" + "/" + name + "@" + version + "?arch=" + architecture
}
