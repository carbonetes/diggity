package config

import "github.com/carbonetes/diggity/pkg/types"

var defaultAllowedPatterns = []string{
	"(.*?)(jpg|gif|doc|docx|zip|xls|pdf|bin|svg|socket|vsidx|v2|suo|wsuo|.dll|pdb|exe)$",
}

// Load config with default values and set of rules
func LoadDefaultConfig() types.SecretConfig {
	return types.SecretConfig{
		Whitelist: types.Whitelist{
			Patterns: defaultAllowedPatterns,
			Keywords: []string{},
		},
		Rules: LoadDefaultRules(),
	}
}

// LoadDefaultRules loads the default rules for secret detection
func LoadDefaultRules() []types.Rule {
	var rules []types.Rule

	// These default rules are just sample rules to show how to add new rules
	rules = append(rules, types.Rule{
		ID:          "AWS_ACCESS_KEY_ID",
		Description: "Access Key is part of the security credentials used to authenticate and authorize activities with AWS (Amazon Web Services). These credentials are used to sign programmatic requests that you make to AWS, whether you're using the AWS Management Console, AWS CLI, or AWS SDKs.",
		Pattern:     `\b(?:A3T[A-Z0-9]|AKIA|AGPA|AIDA|AROA|AIPA|ANPA|ANVA|ASIA)[A-Z0-9]{16}\b`,
		Keywords:    []string{"akia", "agpa", "aida", "aroa", "aipa", "anpa", "anva", "asia"},
	})

	rules = append(rules, types.Rule{
		ID:          "PRIVATE_KEY",
		Description: "A private key, also known as a secret key, is a variable in cryptography that is used with an algorithm to encrypt and decrypt data. Secret keys should only be shared with the key's generator or parties authorized to decrypt the data.",
		Pattern:     `(?i)-----BEGIN[ A-Z0-9_-]{0,100}PRIVATE KEY( BLOCK)?-----[\s\S-]*KEY( BLOCK)?----`,
		Keywords:    []string{"-----BEGIN"},
	})

	rules = append(rules, types.Rule{
		ID:          "JWT_TOKEN",
		Description: "JSON Web Token is a compact, URL-safe means of representing claims to be transferred between two parties.",
		Pattern:     `\b(ey[a-zA-Z0-9]{17,}\.ey[a-zA-Z0-9\/\\_-]{17,}\.(?:[a-zA-Z0-9\/\\_-]{10,}={0,2})?)(?:['|\"|\n|\r|\s|\x60|;]|$)`,
		Keywords:    []string{"ey"},
	})

	return rules
}
