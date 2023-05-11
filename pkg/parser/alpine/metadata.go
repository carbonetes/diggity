package alpine

import (
	"strings"

	"github.com/carbonetes/diggity/pkg/model"
)

type Metadata map[string]interface{}

func parseMetadata(attributes []string, listFiles bool) Metadata {
	var value string
	var key string
	var metadata = make(Metadata)

	for _, attribute := range attributes {
		if strings.Contains(attribute, ":") && !strings.Contains(attribute, ":=") {
			keyValues := strings.SplitN(attribute, ":", 2)
			key = keyValues[0]
			value = keyValues[1]
		} else {
			value = strings.TrimSpace(value + attribute)
		}

		switch key {
		case "A":
			{
				metadata["Architecture"] = value
			}
		case "C":
			{
				metadata["PullChecksum"] = value
			}
		case "D", "r":
			{
				metadata["PullDependencies"] = value
			}
		case "I":
			{
				metadata["PackageInstalledSize"] = value
			}
		case "L":
			{
				metadata["License"] = value
			}
		case "M":
			{
				metadata["Permissions"] = value
			}
		case "P":
			{
				metadata["Name"] = value
			}
		case "S":
			{
				metadata["Size"] = value
			}
		case "T":
			{
				metadata["Description"] = value
			}
		case "U":
			{
				metadata["URL"] = value
			}
		case "V":
			{
				metadata["Version"] = value
			}
		case "c":
			{
				metadata["GitCommitHashApk"] = value
			}
		case "m":
			{
				metadata["Maintainer"] = value
			}
		case "o":
			{
				metadata["Origin"] = value
			}
		case "p":
			{
				metadata["Provides"] = value
			}
		case "t":
			{
				metadata["BuildTimestamp"] = value
			}
		}
		if listFiles {
			files := getAlpineFiles(attribute)
			metadata["Files"] = make([]model.File, len(files))
			if len(files) > 0 {
				metadata["Files"] = append(metadata["Files"].([]model.File), files...)
			}

		}

	}

	return metadata
}
