package sbom

import (
	"errors"
	"strings"

	"github.com/carbonetes/diggity/internal/docker"
	"github.com/carbonetes/diggity/internal/file"
	"github.com/carbonetes/diggity/internal/logger"
	"github.com/carbonetes/diggity/internal/model"
	"github.com/carbonetes/diggity/internal/output"
	"github.com/carbonetes/diggity/internal/parser"
	"github.com/carbonetes/diggity/internal/parser/bom"
	"github.com/carbonetes/diggity/internal/ui"

	"os"

	"github.com/schollz/progressbar/v3"
)

const (
	tarFile      string = "tar"
	image        string = "image"
	dir          string = "dir"
	unknown      string = "Unknown"
	defaultTag   string = "latest"
	tagSeparator string = ":"
)

var (
	log = logger.GetLogger()
)

// Start SBOM extraction
func Start(arguments *model.Arguments) {
	if *arguments.Quiet {
		log = logger.SetQuietMode(log)
	}
	//check image if DIR
	source, spinnerMsg := file.CheckUserInput(arguments)
	if source == image && !strings.Contains(*arguments.Image, tagSeparator) {
		log.Printf("Using default tag:" + defaultTag)
	}

	extractSpinner := ui.InitSpinner(spinnerMsg)
	//Extract Image
	if !*arguments.Quiet {
		// Pull (if needed) and Extract Image
		spinnerMsg = extractImage(source, arguments, extractSpinner)

		// Run Parsers
		parseSpinner := ui.InitSpinner(spinnerMsg)
		go ui.RunSpinner(parseSpinner)
		startParsers(arguments)
		ui.DoneSpinner(parseSpinner)
	} else {
		extractImage(source, arguments, extractSpinner)
		startParsers(arguments)
	}

	//Print Results and Cleanup
	output.PrintResults()
}

// Run parsers
func startParsers(arguments *model.Arguments) {
	bom.InitParsers(*arguments)
	bom.WG.Add(len(parser.FindFunctions))
	for _, parser := range parser.FindFunctions {
		go parser()
	}
	bom.WG.Wait()
	cleanUp()
}

// Clear temp files
func cleanUp() {
	err := os.RemoveAll(docker.Dir())
	if err != nil {
		err = errors.New("clean-up: " + err.Error())
		bom.Errors = append(bom.Errors, &err)
	}
}

// GetResults for event bus
func GetResults() string {
	return output.GetResults()
}

// Extract image
func extractImage(source string, arguments *model.Arguments, spinner *progressbar.ProgressBar) string {
	switch source {
	case tarFile:
		dir := *arguments.Tar
		if file.Exists(dir) {
			docker.ExtractFromDir(arguments.Tar)
			return "Parsing Tar file..."
		}
		log.Printf("%s not found\n", *arguments.Tar)
		os.Exit(0)
	case image:
		if !strings.Contains(*arguments.Image, tagSeparator) {
			modifiedTag := *arguments.Image + tagSeparator + defaultTag
			arguments.Image = &modifiedTag
		}
		docker.ExtractImage(arguments, spinner)
		return "Parsing Image..."
	case dir:
		dir := *arguments.Dir
		if file.Exists(dir) {
			err := file.GetFilesFromDir(dir)
			if err != nil {
				panic(err)
			}

			docker.CreateTempDir()
			return "Parsing Directory..."
		}
		log.Printf("%s not found\n", *arguments.Dir)
		os.Exit(0)
	}
	return ""
}
