package cmd

import (
	"github.com/carbonetes/diggity/internal/scanner"
	"github.com/carbonetes/diggity/pkg/types"
)

func init() {
	// root flags
	root.Flags().StringP("tar", "t", "", "Read a tarball from a path on disk for archives created from docker save (e.g. 'diggity path/to/image.tar)'")
	root.Flags().StringP("directory", "d", "", "Read directly from a path on disk (any directory) (e.g. 'diggity -fs path/to/directory)'")

	root.Flags().StringP("output", "o", types.Table.String(), "Supported output types are: "+types.GetAllOutputFormat())
	root.Flags().StringP("file", "f", "", "Save scan result to a file")
	root.Flags().BoolP("quiet", "q", false, "Suppress all output except for errors")
	root.Flags().StringArray("scanners", scanner.All, "Allow only selected scanners to run (e.g. --scanners apk,dpkg)")
	root.Flags().Bool("allow-file-listing", false, "Allow parsers to list files related to the packages")
}
