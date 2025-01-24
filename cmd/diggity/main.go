package main

import (
	"github.com/carbonetes/diggity/cmd/diggity/command"
	"github.com/carbonetes/diggity/internal/log"
)

func main() {
	if err := command.Run(); err != nil {
		log.Fatal(err)
	}
}
