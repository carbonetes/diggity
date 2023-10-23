package maven

import (
	"path/filepath"
	"strings"

	"github.com/carbonetes/diggity/pkg/model"
)

// Parse pom properties
func parsePomProperties(data string, pkg *model.Package, path string) {

	var value string
	var attribute string

	pomProperties := make(Manifest)
	pomProperties["location"] = filepath.Join(pkg.Name, path)
	pomProperties["name"] = ""

	lines := strings.Split(data, "\n")
	for _, keyValue := range lines {
		if strings.Contains(keyValue, "=") {
			keyValues := strings.Split(keyValue, "=")
			attribute = keyValues[0]
			value = keyValues[1]
		}

		if len(attribute) > 0 && attribute != " " {
			pomProperties[attribute] = strings.Replace(value, "\r\n", "", -1)
			pomProperties[attribute] = strings.Replace(value, "\r ", "", -1)
			pomProperties[attribute] = strings.TrimSpace(pomProperties[attribute])
		}
	}

	pkg.Metadata.(Metadata)["PomProperties"] = pomProperties
}
