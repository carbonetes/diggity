package rubygem

import (
	"bufio"
	"regexp"
	"strings"
)

type Gemspec struct {
	Name        string   `json:"name"`
	Version     string   `json:"version"`
	Authors     []string `json:"authors"`
	Description string   `json:"description"`
	Licenses    []string `json:"licenses"`
}

func parseGemspec(content []byte) (*Gemspec, error) {
	gemspec := &Gemspec{}
	scanner := bufio.NewScanner(strings.NewReader(string(content)))

	re := regexp.MustCompile(`s\.(\w+)\s*=\s*(.+)`)

	for scanner.Scan() {
		line := scanner.Text()
		matches := re.FindStringSubmatch(line)
		if len(matches) > 2 {
			key := matches[1]
			value := strings.Trim(matches[2], ` "%q{}`)

			switch key {
			case "name":
				gemspec.Name = value
			case "version":
				gemspec.Version = value
			case "authors":
				gemspec.Authors = append(gemspec.Authors, value)
			case "description":
				gemspec.Description = value
			case "licenses":
				gemspec.Licenses = append(gemspec.Licenses, value)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return gemspec, nil
}

func cleanMetadata(metadata *Gemspec) {

	// Clean up string fields
	if len(metadata.Name) > 0 {
		metadata.Name = cleanStringField(metadata.Name)
	}

	if len(metadata.Version) > 0 {
		metadata.Version = cleanStringField(metadata.Version)
	}

	if len(metadata.Description) > 0 {
		metadata.Description = cleanStringField(metadata.Description)
	}

	// Convert string representations of arrays to actual arrays
	if len(metadata.Authors) > 0 {
		metadata.Authors = cleanArrayField(metadata.Authors[0])
	}

	if len(metadata.Licenses) > 0 {
		metadata.Licenses = cleanArrayField(metadata.Licenses[0])
	}
}

func cleanArrayField(field string) []string {
	field = strings.Trim(field, "[]")
	items := strings.Split(field, "\", \"")
	for i, item := range items {
		items[i] = cleanStringField(item)
	}
	return items
}

func cleanStringField(field string) string {
	field = strings.ReplaceAll(field, "\"", "")
	field = strings.ReplaceAll(field, ".freeze", "")
	return field
}

func readManifestFile(content []byte) [][]string {
	var attributes [][]string
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) < 5 {
			continue
		}

		if strings.Count(line[:5], " ") != 4 {
			continue
		}

		props := strings.Fields(line)

		if len(props) != 2 {
			continue
		}
		attributes = append(attributes, props)
	}

	return attributes
}
