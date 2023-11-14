package distro

import (
	"strings"

	"github.com/carbonetes/diggity/pkg/types"
)

func parseRelease(manifest types.ManifestFile) (*types.Distro, error) {

	// Parse the os-release content
	lines := strings.Split(string(manifest.Content), "\n")
	distro := types.Distro{}
	for _, line := range lines {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := parts[0]
		value := strings.Trim(parts[1], "\"") // Remove surrounding quotes

		// Populate the Distro struct
		switch key {
		case "PRETTY_NAME":
			distro.PrettyName = value
		case "NAME":
			distro.Name = value
		case "ID":
			distro.ID = value
		case "VERSION":
			distro.Version = value
		case "VERSION_ID":
			distro.VersionID = value
		case "HOME_URL":
			distro.HomeURL = value
		case "SUPPORT_URL":
			distro.SupportURL = value
		case "BUG_REPORT_URL":
			distro.BugReportURL = value
		case "PRIVACY_POLICY_URL":
			distro.PrivacyPolicyURL = value
		}
	}
	return &distro, nil
}
