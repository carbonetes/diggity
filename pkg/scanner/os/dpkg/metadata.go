package dpkg

import (
	"bufio"
	"strings"

	"github.com/carbonetes/diggity/internal/helper"
	"github.com/mitchellh/mapstructure"
)

type Package struct {
	Name         string `mapstructure:"Package"`
	Status       string
	Priority     string
	Section      string
	Size         string `mapstructure:"Installed-Size"`
	Maintainer   string
	Architecture string
	MultiArch    string `mapstructure:"Multi-Arch"`
	Source       string
	Version      string
	Origin       string
	Replaces     []string
	Provides     []string
	Breaks       []string
	Depends      [][]string
	Suggests     [][]string
	PreDepends   [][]string `mapstructure:"Pre-Depends"`
	Conffiles    []Conffile
	Recommends   []string
	Description  string
	Homepage     string
}

type Conffile struct {
	Path string
	Hash string
}

func ParseDpkgDatabase(dpkgDbContent string) ([]Package, error) {
	packages := make([]Package, 0)
	groups := helper.SplitContentsByEmptyLine(dpkgDbContent)
	for _, group := range groups {
		scanner := bufio.NewScanner(strings.NewReader(group))
		m := make(map[string]interface{})
		var prevKey string
		for scanner.Scan() {
			line := scanner.Text()
			if line == "" {
				continue
			}

			if strings.HasPrefix(line, " ") {
				m[prevKey] = m[prevKey].(string) + "\n" + line
				continue
			}

			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 1 {
				m[prevKey] = m[prevKey].(string) + "\n" + parts[0]
				continue
			}

			key, value := parts[0], parts[1]
			prevKey = key
			m[key] = strings.TrimLeft(value, " ") // remove leading whitespaces from the value
		}

		for key, value := range m {
			switch key {
			case "Depends", "Suggests", "Pre-Depends":
				dependencies := strings.Split(value.(string), ", ")
				altDependencies := make([][]string, len(dependencies))
				for i, dep := range dependencies {
					altDepParts := strings.Split(strings.TrimSpace(dep), " | ")
					for j, altDep := range altDepParts {
						depParts := strings.SplitN(altDep, " ", 2)
						altDepParts[j] = depParts[0] // keep only the package name
					}
					altDependencies[i] = altDepParts
				}
				m[key] = altDependencies

			case "Conffiles":
				conffiles := strings.Split(value.(string), "\n")
				conffiles = conffiles[1:] // remove the first empty line
				conffileList := []Conffile{}
				for _, conffile := range conffiles {
					parts := strings.SplitN(strings.TrimSpace(conffile), " ", 2)
					if parts[0] == "" || parts[1] == "" {
						continue
					}
					conffileList = append(conffileList, Conffile{Path: parts[1], Hash: parts[0]})
				}
				m[key] = conffileList
			case "Provides", "Breaks", "Recommends", "Replaces":
				provides := strings.Split(value.(string), ", ")
				m[key] = provides
			}
		}

		var pkg Package
		_ = mapstructure.Decode(m, &pkg)
		packages = append(packages, pkg)

	}
	return packages, nil
}
