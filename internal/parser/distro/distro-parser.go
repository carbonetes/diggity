package distro

import (
	"bufio"
	"errors"
	"io/fs"
	"os"
	"regexp"
	"strings"

	"github.com/carbonetes/diggity/internal/docker"
	"github.com/carbonetes/diggity/internal/file"
	"github.com/carbonetes/diggity/internal/model"
	"github.com/carbonetes/diggity/internal/parser/bom"
)

var distro *model.Distro

// Distro - returns parsed distro information
func Distro() *model.Distro {
	return distro
}

// ParseDistro parses os release
func ParseDistro() {

	var relatedOsFiles []string
	var err error
	_ = os.Mkdir(docker.Dir(), fs.ModePerm)

	osFilesRegex := `etc\/(\S+)-release|etc\\(\S+)-release|usr\\(\S+)-release|usr\/lib\/(\S+)-release|usr\/(\S+)-release`
	fileRegexp, _ := regexp.Compile(osFilesRegex)
	for _, content := range file.Contents {
		if match := fileRegexp.MatchString(content.Path); match {
			relatedOsFiles = append(relatedOsFiles, content.Path)
		}
	}

	distro, err = parseLinuxDistribution(relatedOsFiles)

	if err != nil {
		err = errors.New("distro-parser: " + err.Error())
		bom.Errors = append(bom.Errors, &err)
	}

	defer bom.WG.Done()
}

// Parse Linux distro
func parseLinuxDistribution(filenames []string) (*model.Distro, error) {

	linuxDistribution := make(map[string]string, 0)
	release := new(model.Distro)
	var value string
	var attribute string

	for _, filename := range filenames {
		file, err := os.Open(filename)

		if err != nil {
			return release, err
		}
		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			keyValue := strings.TrimSpace(scanner.Text())
			if strings.Contains(keyValue, "=") {
				keyValues := strings.SplitN(keyValue, "=", 2)
				attribute = keyValues[0]
				value = keyValues[1]
			}
			if len(attribute) > 0 && attribute != " " {
				value = strings.Replace(value, "\r\n", "", -1)
				value = strings.ReplaceAll(value, "\"", "")
				linuxDistribution[attribute] = strings.Replace(value, "\r ", "", -1)
				linuxDistribution[attribute] = strings.TrimSpace(linuxDistribution[attribute])
			}
		}

		file.Close()
	}

	var ids []string
	for _, id := range strings.Split(linuxDistribution["ID_LIKE"], " ") {
		id = strings.TrimSpace(id)
		if id == "" {
			continue
		}
		ids = append(ids, id)
	}

	release = &model.Distro{
		PrettyName:         linuxDistribution["PRETTY_NAME"],
		Name:               linuxDistribution["NAME"],
		ID:                 linuxDistribution["ID"],
		IDLike:             ids,
		Version:            linuxDistribution["VERSION"],
		VersionID:          linuxDistribution["VERSION_ID"],
		DistribID:          linuxDistribution["DISTRIB_ID"],
		DistribDescription: linuxDistribution["DISTRIB_DESCRIPTIONN"],
		DistribCodename:    linuxDistribution["DISTRIB_CODENAME"],
		HomeURL:            linuxDistribution["HOME_URL"],
		SupportURL:         linuxDistribution["SUPPORT_URL"],
		BugReportURL:       linuxDistribution["BUG_REPORT_URL"],
		PrivacyPolicyURL:   linuxDistribution["PRIVACY_POLICY_URL"],
	}

	return release, nil
}
