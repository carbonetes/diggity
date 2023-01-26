package output

import (
	"os"

	"github.com/carbonetes/diggity/internal/parser"
)

func saveResultToFile(result string) {
	file, _ := os.Create(*parser.Arguments.OutputFile)
	err := os.WriteFile(file.Name(), []byte(result), 0644)
	if err != nil {
		panic(err)
	}
}
