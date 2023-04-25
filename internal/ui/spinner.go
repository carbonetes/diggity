package ui

import (
	"github.com/carbonetes/diggity/internal/logger"
	"github.com/schollz/progressbar/v3"
)

var (
	disabled bool
	log      = logger.GetLogger()
)

// InitSpinner generates simple spinner
func InitSpinner(text string) *progressbar.ProgressBar {
	if disabled {
		return nil
	}
	pb := progressbar.NewOptions(-1,
		progressbar.OptionSpinnerType(14),
		progressbar.OptionSetDescription(text),
		progressbar.OptionClearOnFinish(),
	)
	return pb
}

// RunSpinner starts a spinner
func RunSpinner(spinner *progressbar.ProgressBar) {
	if disabled {
		return
	}
	for {
		err := spinner.Add(1)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// DoneSpinner stops and closes a spiner
func DoneSpinner(pb *progressbar.ProgressBar) {
	if disabled {
		return
	}
	err := pb.Finish()
	if err != nil {
		log.Fatal(err)
	}

	pb.Close()
}

// Disable spinner
func Disable() {
	disabled = true
}
