package helper

import (
	"regexp"
	"strings"

	"github.com/CycloneDX/cyclonedx-go"
)

// DetectHashAlgorithm identifies the hash algorithm based on the input hash.
func DetectCDXHashAlgorithm(hash string) cyclonedx.HashAlgorithm {
	// Normalize the hash string (remove any spaces or special characters)
	normalizedHash := strings.ToLower(strings.ReplaceAll(hash, " ", ""))
	matched, _ := regexp.MatchString(`^[0-9a-fA-F]+$`, normalizedHash)
	// Check for specific patterns or prefixes
	switch {
	case strings.HasPrefix(normalizedHash, "sha1:"):
		return cyclonedx.HashAlgoSHA1
	case strings.HasPrefix(normalizedHash, "sha256:"):
		return cyclonedx.HashAlgoSHA256
	case strings.HasPrefix(normalizedHash, "sha512:"):
		return cyclonedx.HashAlgoSHA512
	case matched:
		// If it consists of hexadecimal characters, it could be MD5, SHA-1, SHA-256, or SHA-512
		if len(normalizedHash) == 32 {
			return cyclonedx.HashAlgoMD5
		}
		if len(normalizedHash) == 40 {
			return cyclonedx.HashAlgoSHA1
		}
		if len(normalizedHash) == 64 {
			return cyclonedx.HashAlgoSHA256
		}
		if len(normalizedHash) == 128 {
			return cyclonedx.HashAlgoSHA512
		}
	}

	// If no specific pattern matches, consider it unknown
	return "Unknown"
}
