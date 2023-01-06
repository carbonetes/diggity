package parser

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/carbonetes/diggity/internal/model"
)

type (
	ParseLinuxDistroResult struct {
		filenames []string
		expected  *model.Distro
	}
)

func TestDistro(t *testing.T) {

	var distro interface{} = Distro()

	switch distro.(type) {
	case *model.Distro:
		return
	default:
		t.Errorf("Test Failed: Distro must return *model.Distro.")
	}
}

func TestParseLinuxDistribution(t *testing.T) {
	apkFilenames := []string{
		filepath.Join("..", "..", "docs", "references", "release", "alpine", "alpine-release"),
		filepath.Join("..", "..", "docs", "references", "release", "alpine", "os-release"),
	}

	debFilenames := []string{
		filepath.Join("..", "..", "docs", "references", "release", "debian", "os-release"),
	}

	rpmFilenames := []string{
		filepath.Join("..", "..", "docs", "references", "release", "rpm", "centos-release"),
		filepath.Join("..", "..", "docs", "references", "release", "rpm", "centos-release-upstream"),
		filepath.Join("..", "..", "docs", "references", "release", "rpm", "os-release"),
		filepath.Join("..", "..", "docs", "references", "release", "rpm", "system-release-cpe"),
	}

	apkRelease := model.Distro{
		PrettyName:   "Alpine Linux v3.16",
		Name:         "Alpine Linux",
		ID:           "alpine",
		VersionID:    "3.16.2",
		HomeURL:      "https://alpinelinux.org/",
		BugReportURL: "https://gitlab.alpinelinux.org/alpine/aports/-/issues",
	}

	debRelease := model.Distro{
		PrettyName:   "Debian GNU/Linux 11 (bullseye)",
		Name:         "Debian GNU/Linux",
		ID:           "debian",
		Version:      "11 (bullseye)",
		VersionID:    "11",
		HomeURL:      "https://www.debian.org/",
		SupportURL:   "https://www.debian.org/support",
		BugReportURL: "https://bugs.debian.org/",
	}

	rpmRelease := model.Distro{
		PrettyName: "CentOS Linux 8",
		Name:       "CentOS Linux",
		ID:         "centos",
		IDLike: []string{
			"rhel",
			"fedora",
		},
		Version:      "8",
		VersionID:    "8",
		HomeURL:      "https://centos.org/",
		BugReportURL: "https://bugs.centos.org/",
	}

	tests := []ParseLinuxDistroResult{
		{apkFilenames, &apkRelease},
		{debFilenames, &debRelease},
		{rpmFilenames, &rpmRelease},
	}

	for _, test := range tests {
		output, err := parseLinuxDistribution(test.filenames)

		if err != nil {
			t.Error("Test Failed: Error occurred while parsing Linux Distribution.")
		}

		if !reflect.DeepEqual(output, test.expected) {
			t.Errorf("Test Failed:\n Expected output of %v \n, Received: %v \n", test.expected, test.filenames)
		}
	}
}
