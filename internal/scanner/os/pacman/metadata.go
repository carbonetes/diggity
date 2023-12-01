package pacman

import (
	"strings"
)

func parseMetadata(attributes []string) map[string]interface{} {
	var metadata = make(map[string]interface{})
	// Attributes based on https://gitlab.archlinux.org/pacman/pacman/-/blob/master/lib/libalpm/be_local.c
	for _, attribute := range attributes {
		if attribute == "" {
			continue
		}
		attribute = strings.TrimSpace(attribute)
		properties := strings.Split(attribute, "\n")
		key := strings.ReplaceAll(properties[0], "%", "")
		values := properties[1:]
		if len(values) > 1 {
			metadata[key] = values
		} else {
			metadata[key] = values[0]
		}
	}
	return metadata
}
