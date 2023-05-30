package portage

import (
	"bufio"
	"os"
	"strings"

	"github.com/carbonetes/diggity/pkg/model/metadata"
)

// Get Portage Files
func getPortageFiles(md *metadata.PortageMetadata, loc string) error {
	var files []metadata.PortageFile

	// Parse CONTENT file
	file, err := os.Open(loc)
	if err != nil {
		if strings.Contains(err.Error(), noFileErrWin) || strings.Contains(err.Error(), noFileErrMac) {
			return nil
		}
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		content := scanner.Text()
		if strings.Contains(content, portageObj) {
			files = append(files, parsePortageFile(content))
		}
	}

	md.Files = files

	return nil
}

// Parse Portage Files
func parsePortageFile(content string) metadata.PortageFile {
	var file metadata.PortageFile
	var digest metadata.PortageDigest

	obj := strings.Split(content, " ")
	// digest
	if len(obj) > 2 {
		digest.Algorithm = portageAlgorithm
		digest.Value = obj[2]
	}
	// file
	file.Path = obj[1]
	file.Digest = digest

	return file
}
