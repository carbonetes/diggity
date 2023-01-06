package cmd

import (
	"os"
)

func Execute() {
	err := diggity.Execute()
	if err != nil {
		os.Exit(1)
	}
}
