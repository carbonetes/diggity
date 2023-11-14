package apk

import (
	"strings"

	"github.com/carbonetes/diggity/pkg/types"
	"github.com/google/uuid"
)

func newComponent(content string) types.Component {
	content = strings.TrimSpace(content)
	attributes := strings.Split(content, "\n")

	metadata := parseMetadata(attributes)
	if metadata["Name"] == nil {
		return types.Component{}
	}

	var licenses []string
	for _, license := range strings.Split(metadata["License"].(string), " ") {
		if !strings.Contains(strings.ToLower(license), "and") {
			licenses = append(licenses, license)
		}
	}

	return types.Component{
		ID:          uuid.NewString(),
		Name:        metadata["Name"].(string),
		Version:     metadata["Version"].(string),
		Type:        Type,
		Description: metadata["Description"].(string),
		Licenses:    licenses,
		Metadata:    metadata,
	}
}
