package cmd

import "os"

func Execute() {
	err := root.Execute()
	if err != nil {
		os.Exit(1)
	}
}
