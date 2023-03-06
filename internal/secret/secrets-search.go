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
	"github.com/carbonetes/diggity/internal/parser/bom"
	parserUtil "github.com/carbonetes/diggity/internal/parser/util"

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

	if *bom.Arguments.DisableSecretSearch {
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
				bom.Errors = append(bom.Errors, &err)
			}
			stat, err := file.Stat()

			if isExcluded(file.Name()) {
				continue
			}

			if (stat.Size() >= bom.Arguments.SecretMaxFileSize && !util.IsText(buf)) || strings.HasSuffix(stat.Name(), ".tar") {
				file.Close()
				continue
			}

			if stat, err := file.Stat(); !stat.IsDir() {

				if err != nil {
					err = errors.New("secrets: " + err.Error())
					bom.Errors = append(bom.Errors, &err)
				}

				scanner := bufio.NewScanner(file)

				lineNumber := 1
				for scanner.Scan() {
					scannerText := scanner.Text()
					if match := regexp.MustCompile(*bom.Arguments.SecretContentRegex).FindString(scannerText); len(match) > 0 {
						secrets = append(secrets, model.Secret{
							ContentRegexName: match,
							FilePath:         parserUtil.TrimUntilLayer(model.Location{Path: content.Path, LayerHash: content.LayerHash}),
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
						bom.Errors = append(bom.Errors, &err)
					}
				}

			}

			file.Close()

		}

		SecretResults.Configuration = model.SecretConfig{
			Disabled:    *bom.Arguments.DisableSecretSearch,
			SecretRegex: *bom.Arguments.SecretContentRegex,
			Excludes:    bom.Arguments.ExcludedFilenames,
			MaxFileSize: bom.Arguments.SecretMaxFileSize,
		}
		SecretResults.Secrets = secrets
	}

	defer bom.WG.Done()
}

// Check if filename is excluded from search
func isExcluded(filename string) bool {
	if bom.Arguments.ExcludedFilenames == nil {
		return false
	}
	for _, exclude := range *bom.Arguments.ExcludedFilenames {
		if strings.Contains(filename, exclude) {
			return true
		}
	}
	return false
}
