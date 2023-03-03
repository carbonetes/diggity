package save

import (
	"os"

	"github.com/carbonetes/diggity/internal/parser/bom"
)

// ResultToFile saves output to a file
func ResultToFile(result string) {
	file, _ := os.Create(*bom.Arguments.OutputFile)
	err := os.WriteFile(file.Name(), []byte(result), 0644)
	if err != nil {
		panic(err)
	}
}
