package secret

import (
	"regexp"
	"strings"

	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/pkg/config"
	"github.com/carbonetes/diggity/pkg/stream"
	"github.com/carbonetes/diggity/pkg/types"
	"github.com/golistic/urn"
)

const Type = "secret"

type MatchPattern struct {
	Name        string
	Description string
	Pattern     *regexp.Regexp
	Keywords    []string
}

var (
	// Secrets      []types.Secret
	secretConfig types.SecretConfig
	rules        []MatchPattern
	whitelist    []*regexp.Regexp
)

func init() {
	secretConfig = config.Config.SecretConfig
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

func New(addr *urn.URN) {
	secretAddr := *addr
	secretAddr.NID = "secret"
	stream.Set(secretAddr.String(), []types.Secret{})
}

func Scan(data interface{}) interface{} {
	payload, ok := data.(types.Payload)
	if !ok {
		log.Error("Secret received unknown file type")
		return nil
	}

	manifest, ok := payload.Body.(types.ManifestFile)
	if !ok {
		log.Error("Secret received unknown file type")
		return nil
	}

	if manifest.Content == nil {
		return nil
	}

	secretAddr := *payload.Address
	secretAddr.NID = "secret"

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
				File:        manifest.Path,
			}

			if len(payload.Layer) > 0 {
				secret.Layer = payload.Layer
			}

			AddSecret(&secretAddr, secret)
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

func AddSecret(addr *urn.URN, secret types.Secret) {
	data, _ := stream.Get(addr.String())
	secrets, _ := data.([]types.Secret)
	secrets = append(secrets, secret)
	stream.Set(addr.String(), secrets)
}
