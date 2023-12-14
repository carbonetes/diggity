package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/carbonetes/diggity/internal/version"
	"github.com/spf13/cobra"
)

var (
	versionInput string
	versionCmd   = &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Args:  cobra.MaximumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(versionInput) == 0 {
				versionInput = "text"
			}

			info := version.FromBuild()
			switch versionInput {
			case "text":
				// Version
				fmt.Println("Application:         ", info.AppName)
				fmt.Println("Version:             ", info.Version)
				fmt.Println("Build Date:          ", info.BuildDate)
				// Git
				fmt.Println("Git Commit:          ", info.GitCommit)
				fmt.Println("Git Description:     ", info.GitDesc)
				// Golang
				fmt.Println("Go Version:          ", info.GoVersion)
				fmt.Println("Compiler:            ", info.Compiler)
				fmt.Println("Platform:            ", info.Platform)
			case "json":

				jsonFormat := json.NewEncoder(os.Stdout)
				jsonFormat.SetEscapeHTML(false)
				jsonFormat.SetIndent("", " ")
				err := jsonFormat.Encode(&struct {
					version.Version
				}{
					Version: info,
				})
				if err != nil {
					return fmt.Errorf("show version information error: %+v", err)
				}
			default:
				return fmt.Errorf("unrecognize output format: %s", versionInput)
			}
			return nil
		},
	}
)
