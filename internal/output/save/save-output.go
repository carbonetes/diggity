package save

import (
	"os"
)

// ResultToFile saves output to a file
func ResultToFile(result string, filename *string) {
	file, _ := os.Create(*filename)
	err := os.WriteFile(file.Name(), []byte(result), 0644)
	if err != nil {
		panic(err)
	}
}
