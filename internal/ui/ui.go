package ui

import (
	"time"

	"github.com/carbonetes/diggity/internal/logger"
	"github.com/schollz/progressbar/v3"
)

var (
	disabled bool = true
	log           = logger.GetLogger()
	pb       *progressbar.ProgressBar
)

func init() {
	pb = progressbar.NewOptions(-1,
		progressbar.OptionSpinnerType(14),
		progressbar.OptionClearOnFinish(),
	)
}

func run() {
	if disabled {
		return
	}
	for {
		err := pb.Add(1)
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(40 * time.Millisecond)
	}
}

func DoneSpinner() {
	if disabled {
		return
	}
	err := pb.Finish()
	if err != nil {
		log.Fatal(err)
	}

	pb.Close()
}

func Enable() {
	disabled = false
}
