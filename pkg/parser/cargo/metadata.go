package cargo

import (
	"strings"

	"github.com/carbonetes/diggity/pkg/parser/util"
)

func parseMetadata(pkg string) *Metadata {
	if len(pkg) == 0 {
		return nil
	}

	metadata := make(Metadata)
	pkg = strings.TrimSpace(pkg)
	attributes := strings.Split(pkg, "\n")
	attributes = attributes[1:]
	// Attributes are encoded in https://github.com/rust-lang/cargo/blob/master/src/cargo/ops/lockfile.rs
	for index, attribute := range attributes {
		if !strings.Contains(attribute, "=") {
			continue
		}
		properties := strings.Split(attribute, "=")
		key := util.ToTitle(properties[0])
		if strings.Contains(attribute, "[") {
			values := attributes[index+1 : len(attributes)-1]
			metadata[key] = values
		} else {
			value := properties[1]
			metadata[key] = value
		}
	}
	return &metadata
}
