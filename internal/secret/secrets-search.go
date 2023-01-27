package secret

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/carbonetes/diggity/internal/file"
	"github.com/carbonetes/diggity/internal/model"
	"github.com/carbonetes/diggity/internal/parser"

	"golang.org/x/tools/godoc/util"
)

var (
	// Secrets collected secretes
	secrets = make([]model.Secret, 0)
	// SecretResults the final result that will be displayed
	SecretResults = &model.SecretResults{}
)

// Search search secrets in all file contents that does not exceed the max-file-size argument
func Search() {

	if *parser.Arguments.DisableSecretSearch {
		secrets = nil
	} else {
		for _, content := range file.Contents {

			file, _ := os.Open(content.Path)
			if file == nil {
				continue
			}

			// continue if the path is directory
			fs, _ := os.Stat(content.Path)
			if fs.Mode().IsDir() {
				continue
			}

			buf, err := os.ReadFile(content.Path)
			if err != nil {
				err = errors.New("secrets: " + err.Error())
				parser.Errors = append(parser.Errors, &err)
			}
			stat, err := file.Stat()

			if isExcluded(file.Name()) {
				continue
			}

			if (stat.Size() >= parser.Arguments.SecretMaxFileSize && !util.IsText(buf)) || strings.HasSuffix(stat.Name(), ".tar") {
				file.Close()
				continue
			}

			if stat, err := file.Stat(); !stat.IsDir() {

				if err != nil {
					err = errors.New("secrets: " + err.Error())
					parser.Errors = append(parser.Errors, &err)
				}

				scanner := bufio.NewScanner(file)

				lineNumber := 1
				for scanner.Scan() {
					scannerText := scanner.Text()
					if match := regexp.MustCompile(*parser.Arguments.SecretContentRegex).FindString(scannerText); len(match) > 0 {
						secrets = append(secrets, model.Secret{
							ContentRegexName: match,
							FilePath:         parser.TrimUntilLayer(model.Location{Path: content.Path, LayerHash: content.LayerHash}),
							LineNumber:       fmt.Sprintf("%d", lineNumber),
							FileName:         stat.Name(),
						})
					}

					lineNumber++
					if err := scanner.Err(); err != nil {
						if err == bufio.ErrTooLong {
							continue
						}
						err = errors.New("secrets: " + err.Error())
						parser.Errors = append(parser.Errors, &err)
					}
				}

			}

			file.Close()

		}

		SecretResults.Configuration = model.SecretConfig{
			Disabled:    *parser.Arguments.DisableSecretSearch,
			SecretRegex: *parser.Arguments.SecretContentRegex,
			Excludes:    parser.Arguments.ExcludedFilenames,
			MaxFileSize: parser.Arguments.SecretMaxFileSize,
		}
		SecretResults.Secrets = secrets
	}

	defer parser.WG.Done()
}

// Check if filename is excluded from search
func isExcluded(filename string) bool {
	if parser.Arguments.ExcludedFilenames == nil {
		return false
	}
	for _, exclude := range *parser.Arguments.ExcludedFilenames {
		if strings.Contains(filename, exclude) {
			return true
		}
	}
	return false
}
