package apk

import (
	"bufio"
	"regexp"
	"strings"

	"github.com/mitchellh/mapstructure"
)

type ApkIndexRecord struct {
	Checksum         string   `mapstructure:"C"`
	Package          string   `mapstructure:"P"`
	Version          string   `mapstructure:"V"`
	Architecture     string   `mapstructure:"A"`
	Size             string   `mapstructure:"S"`
	InstalledSize    string   `mapstructure:"I"`
	Description      string   `mapstructure:"T"`
	URL              string   `mapstructure:"U"`
	Licenses         []string `mapstructure:"L"`
	Origin           string   `mapstructure:"o"`
	Maintainer       string   `mapstructure:"m"`
	BuildTimestamp   string   `mapstructure:"t"`
	GitCommit        string   `mapstructure:"c"`
	ProviderPriority string   `mapstructure:"k"`
	Dependencies     []string `mapstructure:"D"`
	Provides         []string `mapstructure:"p"`
	Install          string   `mapstructure:"i"`
}

const licenseDelimiter = " AND "

func ParseApkIndexFile(apkDBContent string) ([]ApkIndexRecord, error) {
	var records []ApkIndexRecord
	recordMap := make(map[string]interface{})
	scanner := bufio.NewScanner(strings.NewReader(apkDBContent))
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			var record ApkIndexRecord
			mapstructure.Decode(recordMap, &record)
			records = append(records, record)
			recordMap = make(map[string]interface{})
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key, value := parts[0], strings.TrimSpace(parts[1])
		switch key {
		case "D", "p", "L":
			recordMap[key] = splitValues(value)
		default:
			recordMap[key] = value
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return records, nil
}

func SplitLicense(license string) []string {
	if !strings.Contains(license, licenseDelimiter) {
		return []string{license}
	}
	return strings.Split(license, licenseDelimiter)
}

func splitValues(value string) (values []string) {
	if strings.Contains(value, " ") {
		return strings.Split(value, " ")
	}

	if strings.Contains(value, licenseDelimiter) {
		return strings.Split(value, licenseDelimiter)
	}

	return []string{value}
}

func SplitDependencies(value string) (dependencies []string) {
	props := strings.Split(value, " ")
	constr := regexp.MustCompile(`[=><]`)
	for _, p := range props {
		if constr.MatchString(p) {
			parts := constr.Split(p, 2)
			packageName := parts[0]
			dependencies = append(dependencies, packageName)
		}
	}
	return dependencies
}
