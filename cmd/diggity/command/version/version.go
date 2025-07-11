package version

import (
	"encoding/json"
	"fmt"

	"github.com/carbonetes/diggity/pkg/version"
	"github.com/spf13/cobra"
)

// NewVersionCommand creates the version command
func NewVersionCommand() *cobra.Command {
	var (
		outputFormat string
		shortFormat  bool
	)

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Long:  "Display detailed version and build information for diggity",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runVersion(versionOptions{
				outputFormat: outputFormat,
				shortFormat:  shortFormat,
			})
		},
	}

	// Add flags
	cmd.Flags().StringVarP(&outputFormat, "format", "f", "text", "output format (text, json)")
	cmd.Flags().BoolVarP(&shortFormat, "short", "s", false, "show short version only")

	return cmd
}

type versionOptions struct {
	outputFormat string
	shortFormat  bool
}

func runVersion(opts versionOptions) error {
	info := version.GetInfo()

	switch {
	case opts.shortFormat:
		fmt.Println(version.GetVersion())
	case opts.outputFormat == "json":
		jsonData, err := json.MarshalIndent(info, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal version info to JSON: %w", err)
		}
		fmt.Println(string(jsonData))
	default:
		fmt.Println(info.DetailedString())
	}

	return nil
}
