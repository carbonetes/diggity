package debian

import (
	"path/filepath"
	"strings"

	"github.com/carbonetes/diggity/pkg/model"
)

// Parse files found on metadata
func parseConffiles(files []string) *[]model.Conffile {
	conffiles := new([]model.Conffile)

	for _, f := range files {
		var conffile model.Conffile
		properties := strings.Split(f, " ")
		if len(properties) < 1 {
			continue
		}
		conffile.Path = properties[0]
		conffile.Digest.Algorithm = "md5"
		conffile.Digest.Value = properties[1]
		if filepath.Ext(conffile.Path) == ".conf" {
			conffile.IsConfigFile = true
		}
		*conffiles = append(*conffiles, conffile)
	}
	return conffiles
}
