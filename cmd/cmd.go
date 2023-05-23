package cmd

import (
	"os"
)

// Execute Diggity
func Execute() {
	err := diggity.Execute()
	if err != nil {
		os.Exit(1)
	}
}
