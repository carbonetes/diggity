package secret

import (
	"regexp"
	"strings"

	"github.com/carbonetes/diggity/internal/config"
	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/pkg/stream"
	"github.com/carbonetes/diggity/pkg/types"
)

const Type = "secret"

type MatchPattern struct {
	Name        string
	Description string
	Pattern     *regexp.Regexp
	Keywords    []string
}

var (
	secretConfig types.SecretConfig
	rules        []MatchPattern
	whitelist    []*regexp.Regexp
	c            *types.Config
)

func init() {
	c = config.Load()
	if c == nil {
		return
	}

	secretConfig = c.SecretConfig
	for _, rule := range secretConfig.Rules {
		rules = append(rules, MatchPattern{
			Name:        rule.ID,
			Description: rule.Description,
			Pattern:     regexp.MustCompile(rule.Pattern),
			Keywords:    rule.Keywords,
		})
	}

	if len(secretConfig.Whitelist.Patterns) > 0 {
		for _, pattern := range secretConfig.Whitelist.Patterns {
			whitelist = append(whitelist, regexp.MustCompile(pattern))
		}
	}
}

func Scan(data interface{}) interface{} {
	if c == nil {
		return nil
	}

	manifest, ok := data.(types.ManifestFile)
	if !ok {
		log.Error("Secret received unknown file type")
	}

	if manifest.Content == nil {
		return nil
	}

	content := string(manifest.Content)
	for _, matcher := range rules {
		for _, keyword := range matcher.Keywords {
			if !strings.Contains(content, keyword) {
				continue
			}
		}
		if match := matcher.Pattern.FindString(content); match != "" {
			secret := types.Secret{
				Match:       matcher.Name,
				Description: matcher.Description,
				Content:     match,
				File: manifest.Path,
			}
			stream.AddSecret(secret)
		}
	}
	return data
}

func CheckRelatedFile(file string) (string, bool, bool) {
	if len(whitelist) == 0 {
		return Type, false, false
	}

	for _, pattern := range whitelist {
		if pattern.MatchString(file) {
			return Type, false, false
		}
	}

	return Type, true, true
}
