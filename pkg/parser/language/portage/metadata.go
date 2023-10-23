package portage

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/model/metadata"
)

func initPortageMetadata(p *model.Package, loc string, noFileListing *bool) error {
	var metadata metadata.PortageMetadata
	sizePath := strings.Replace(loc, portageContent, portageSize, -1)

	// Find and parse SIZE file
	file, err := os.Open(sizePath)
	if err != nil {
		if strings.Contains(err.Error(), noFileErrWin) || strings.Contains(err.Error(), noFileErrMac) {
			return nil
		}
		return err
	}

	defer file.Close()

	// Get Size Metadata
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		size, err := strconv.Atoi(scanner.Text())
		if err != nil {
			continue
		}
		metadata.Size = size
	}

	// Get files metadata
	if !*noFileListing {
		if err := getPortageFiles(&metadata, loc); err != nil {
			return nil
		}
	}

	p.Metadata = metadata
	return nil
}
