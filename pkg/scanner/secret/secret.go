package secret

import (
	"regexp"
	"strings"
	"sync"

	"github.com/carbonetes/diggity/internal/config"
	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/pkg/stream"
	"github.com/carbonetes/diggity/pkg/types"
)

const Type = "secret"

// var SecretsPatterns = map[string]string{
// 	"aws-access-token":           `\b(?:A3T[A-Z0-9]|AKIA|AGPA|AIDA|AROA|AIPA|ANPA|ANVA|ASIA)[A-Z0-9]{16}\b`,
// 	"github-app-token":           `(ghu|ghs)_[0-9a-zA-Z]{36}`,
// 	"google-api-key":             `AIza[0-9A-Za-z\\-_]{35}`,
// 	"slack-token":                `xox[baprs]-([0-9a-zA-Z]{10,48})?`,
// 	"ssh-private-key":            `-----BEGIN (?:EC|PGP|DSA|RSA|OPENSSH) PRIVATE KEY(?: BLOCK)?-----`,
// 	"stripe-api-key":             `(?i)stripe(?:.js)?\W*[sk_live|sk_test]_[0-9a-zA-Z]{24}`,
// 	"github-fine-grained-pat":    `github_pat_[0-9a-zA-Z_]{82}`,
// 	"github-oauth":               `gho_[0-9a-zA-Z]{36}`,
// 	"github-pat":                 `ghp_[0-9a-zA-Z]{36}`,
// 	"github-refresh-token":       `ghr_[0-9a-zA-Z]{36}`,
// 	"google-oauth":               `ya29\.[0-9A-Za-z\-_]+`,
// 	"google-service-account":     `(?i)google-?application-?credentials\.json`,
// 	"google-oauth-refresh-token": `1//0[0-9A-Za-z\-_]+`,
// 	"google-gcp-service-account": `(?i)gcp-?application-?credentials\.json`,
// 	"google-gcp-api-key":         `AIza[0-9A-Za-z\\-_]{35}`,
// 	"gitlab-pat":                 `glpat-[0-9a-zA-Z\-\_]{20}`,
// 	"gitlab-ptt":                 `glptt-[0-9a-f]{40}`,
// 	"gitlab-rrt":                 `GR1348941[0-9a-zA-Z\-\_]{20}`,
// 	"jwt-token":                  `\b(ey[a-zA-Z0-9]{17,}\.ey[a-zA-Z0-9\/\\_-]{17,}\.(?:[a-zA-Z0-9\/\\_-]{10,}={0,2})?)(?:['|\"|\n|\r|\s|\x60|;]|$)`,
// 	"jwt-token-base64encoded":    `\bZXlK(?:(?P<alg>aGJHY2lPaU)|(?P<apu>aGNIVWlPaU)|(?P<apv>aGNIWWlPaU)|(?P<aud>aGRXUWlPaU)|(?P<b64>aU5qUWlP)|(?P<crit>amNtbDBJanBi)|(?P<cty>amRIa2lPaU)|(?P<epk>bGNHc2lPbn)|(?P<enc>bGJtTWlPaU)|(?P<jku>cWEzVWlPaU)|(?P<jwk>cWQyc2lPb)|(?P<iss>cGMzTWlPaU)|(?P<iv>cGRpSTZJ)|(?P<kid>cmFXUWlP)|(?P<key_ops>clpYbGZiM0J6SWpwY)|(?P<kty>cmRIa2lPaUp)|(?P<nonce>dWIyNWpaU0k2)|(?P<p2c>d01tTWlP)|(?P<p2s>d01uTWlPaU)|(?P<ppt>d2NIUWlPaU)|(?P<sub>emRXSWlPaU)|(?P<svt>emRuUWlP)|(?P<tag>MFlXY2lPaU)|(?P<typ>MGVYQWlPaUp)|(?P<url>MWNtd2l)|(?P<use>MWMyVWlPaUp)|(?P<ver>MlpYSWlPaU)|(?P<version>MlpYSnphVzl1SWpv)|(?P<x>NElqb2)|(?P<x5c>NE5XTWlP)|(?P<x5t>NE5YUWlPaU)|(?P<x5ts256>NE5YUWpVekkxTmlJNkl)|(?P<x5u>NE5YVWlPaU)|(?P<zip>NmFYQWlPaU))[a-zA-Z0-9\/\\_+\-\r\n]{40,}={0,2}`,
// }

type MatchPattern struct {
	Match   string
	Pattern *regexp.Regexp
}

var (
	MatchPatterns []MatchPattern
	ExcludedPattern *regexp.Regexp
)

func init() {
	config := stream.GetConfig()
	for name, pattern := range SecretsPatterns {
		MatchPatterns = append(MatchPatterns, MatchPattern{
			Match:   name,
			Pattern: regexp.MustCompile(pattern),
		})
	}
	ExcludedPattern = regexp.MustCompile(`(.*?)(jpg|gif|doc|docx|zip|xls|pdf|bin|svg|socket|vsidx|v2|suo|wsuo|\.dll|\.pdb|\.exe)$`)
}

func Scan(data interface{}) interface{} {
	manifest, ok := data.(types.ManifestFile)
	if !ok {
		log.Error("Secret received unknown file type")
	}

	lines := strings.Split(string(manifest.Content), "\n")
	var wg sync.WaitGroup
	wg.Add(len(lines))
	for index, line := range lines {
		go func(index int, line string) {
			for _, matchPattern := range MatchPatterns {
				if match := matchPattern.Pattern.FindString(line); match != "" {
					secret := types.Secret{
						Match:   matchPattern.Match,
						Content: match,
						File:    manifest.Path,
						Line:    index + 1,
					}
					stream.AddSecret(secret)
				}
			}
		}(index, line)
		wg.Done()
	}
	wg.Wait()
	return data
}

func CheckRelatedFile(file string) (string, bool, bool) {
	if match := ExcludedPattern.FindString(file); match != "" {
		return "", false, false
	}

	return Type, true, true
}
