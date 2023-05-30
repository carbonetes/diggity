package alpm

import (
	"io"
	"os"
	"strings"

	"github.com/carbonetes/diggity/pkg/parser/util"
)

func getAlpmFiles(path string) (*[]string, *[]string, error) {
	var files []string
	var backups []string
	path = strings.Replace(path, "desc", "files", -1)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, nil, err
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, nil, err
	}
	contents := string(data)
	contents = strings.TrimSpace(contents)
	attributes := util.SplitContentsByEmptyLine(contents)
	for _, attribute := range attributes {
		if attribute == "" {
			continue
		}
		attribute = strings.TrimSpace(attribute)
		properties := strings.Split(attribute, "\n")
		key := properties[0]
		values := properties[1:]
		switch key {
		case "%FILES%":
			files = values
		case "%BACKUP%":
			backups = values
		}
	}

	return &files, &backups, nil
}
