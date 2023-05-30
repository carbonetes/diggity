package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	versionPackage "github.com/carbonetes/diggity/internal/version"
	"github.com/carbonetes/diggity/pkg/model"
	"github.com/spf13/cobra"
)

var (
	version = &cobra.Command{
		Use:   "version",
		Short: "Display Build Version Info Diggity",
		Long:  "Display Build Version Info Diggity",
		Args:  cobra.MaximumNArgs(0),
		RunE: func(_ *cobra.Command, _ []string) error {

			versionInfo := versionPackage.FromBuild()
			switch versionOutputFormat {
			case "text":
				// Version
				fmt.Println("Application:         ", versionInfo.AppName)
				fmt.Println("Version:             ", versionInfo.Version)
				fmt.Println("Build Date:          ", versionInfo.BuildDate)
				// Git
				fmt.Println("Git Commit:          ", versionInfo.GitCommit)
				fmt.Println("Git Description:     ", versionInfo.GitDesc)
				// Golang
				fmt.Println("Go Version:          ", versionInfo.GoVersion)
				fmt.Println("Compiler:            ", versionInfo.Compiler)
				fmt.Println("Platform:            ", versionInfo.Platform)
			case "json":

				jsonFormat := json.NewEncoder(os.Stdout)
				jsonFormat.SetEscapeHTML(false)
				jsonFormat.SetIndent("", " ")
				err := jsonFormat.Encode(&struct {
					model.Version
				}{
					Version: versionInfo,
				})
				if err != nil {
					return fmt.Errorf("show version information error: %+v", err)
				}
			default:
				return fmt.Errorf("unrecognize output format: %s", versionOutputFormat)
			}
			return nil
		},
	}
)
