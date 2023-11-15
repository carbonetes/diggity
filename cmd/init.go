package cmd

import (
	"github.com/carbonetes/diggity/pkg/types"
)

func init() {
	// root flags
	root.Flags().StringP("tar", "t", "", "Read a tarball from a path on disk for archives created from docker save (e.g. 'diggity path/to/image.tar)'")
	root.Flags().StringP("dir", "d", "", "Read directly from a path on disk (any directory) (e.g. 'diggity path/to/dir)'")

	root.Flags().StringP("output", "o", string(types.Table), "Supported output types")
	root.Flags().BoolP("quiet", "q", false, "Disable all output except scan result")
	root.Flags().StringArray("parsers", []string{}, "Allow only selected parsers to run")
	root.Flags().Bool("allow-file-listing", false, "Allow parsers to list files related to the packages")
}
