package alpm

import (
	"io"
	"os"
	"strings"
)

func getAlpmFiles(path string) *[]string {

	path = strings.Replace(path, "desc", "files", -1)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}

	file, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil
	}
	contents := string(data)
	contents = strings.TrimSpace(contents)
	files := strings.Split(contents, "\n")[1:]

	return &files
}
