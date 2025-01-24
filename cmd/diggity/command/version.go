package command

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/carbonetes/diggity/cmd/diggity/build"
	"github.com/carbonetes/diggity/internal/log"
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

			info := build.FromBuild()
			switch versionInput {
			case "text":
				// Version
				log.Print("Application:         ", info.AppName)
				log.Print("Version:             ", info.Version)
				log.Print("Build Date:          ", info.BuildDate)
				// Git
				log.Print("Git Commit:          ", info.GitCommit)
				log.Print("Git Description:     ", info.GitDesc)
				// Golang
				log.Print("Go Version:          ", info.GoVersion)
				log.Print("Compiler:            ", info.Compiler)
				log.Print("Platform:            ", info.Platform)
			case "json":

				jsonFormat := json.NewEncoder(os.Stdout)
				jsonFormat.SetEscapeHTML(false)
				jsonFormat.SetIndent("", " ")
				err := jsonFormat.Encode(&struct {
					build.Version
				}{
					Version: info,
				})
				if err != nil {
					return err
				}
			default:
				return errors.New("invalid output format: " + versionInput)
			}
			return nil
		},
	}
)