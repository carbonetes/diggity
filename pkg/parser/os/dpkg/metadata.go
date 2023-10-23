package dpkg

import (
	"strings"
)

// Metadata debian metadata
type Metadata map[string]interface{}

func parseMetadata(pkg string, listFiles bool) *Metadata {
	metadata := make(Metadata)
	attributes := strings.Split(pkg, "\n")
	var descriptions []string
	var conffiles []string
	var key string
	for _, attribute := range attributes {
		if attribute == " ." {
			continue
		}

		if strings.HasPrefix(attribute, " ") {
			attribute = strings.TrimSpace(attribute)
			switch key {
			case "description":
				descriptions = append(descriptions, attribute)
			case "conffiles":
				conffiles = append(conffiles, attribute)
			}
			continue
		}
		attribute = strings.TrimSpace(attribute)
		parts := strings.SplitN(attribute, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key = strings.ToLower(strings.TrimSpace(parts[0]))
		value := strings.TrimSpace(parts[1])
		metadata[key] = value
	}

	if len(descriptions) > 0 {
		metadata["description"] = strings.Join(descriptions, " ")
	}

	if len(conffiles) > 0 && listFiles {
		files := parseConffiles(conffiles)
		metadata["conffiles"] = files
	}

	return &metadata
}
