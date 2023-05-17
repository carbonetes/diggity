package distro

import (
	"errors"
	"io/fs"
	"os"
	"regexp"

	"github.com/carbonetes/diggity/pkg/parser/bom"
)

const parserErr = "distro-parser: "

func ParseDistro(req *bom.ParserRequirements) {

	var relatedOsFiles []string
	var err error
	_ = os.Mkdir(*req.DockerTemp, fs.ModePerm)

	osFilesRegex := `etc\/(\S+)-release|etc\\(\S+)-release|usr\\(\S+)-release|usr\/lib\/(\S+)-release|usr\/(\S+)-release`
	fileRegexp, _ := regexp.Compile(osFilesRegex)
	for _, content := range *req.Contents {
		if match := fileRegexp.MatchString(content.Path); match {
			relatedOsFiles = append(relatedOsFiles, content.Path)
		}
	}

	distro, err := parseLinuxDistribution(relatedOsFiles)

	if err != nil {
		err = errors.New(parserErr + err.Error())
		*req.Errors = append(*req.Errors, err)
	}
	req.SBOM.Distro = distro

	defer req.WG.Done()
}
