package debian

import (
	"fmt"
	"path/filepath"
	"reflect"
	"sort"
	"testing"

	"github.com/carbonetes/diggity/internal/model"
)

type (
	DebPurlResult struct {
		_package *model.Package
		arch     string
		expected model.PURL
	}

	DebLicenseResult struct {
		_package *model.Package
		path     string
		expected []string
	}

	ParseDebianFilesResult struct {
		metadata Metadata
		content  string
		expected Metadata
	}

	InitDebPackageResult struct {
		_package *model.Package
		location *model.Location
		metadata Metadata
		expected *model.Package
	}
)

var (
	debPackage1 = model.Package{
		Name:    "libpcre2-8-0",
		Type:    debType,
		Version: "10.36-2",
		Path:    dpkgStatusPath,
		Locations: []model.Location{
			{
				Path:      dpkgStatusPath,
				LayerHash: "f1a5f5ce6b163fac7f09b47645c56d2ab676bdcdb268eef06a4d9b782a75bfd0",
			},
		},
		Description: "New Perl Compatible Regular Expression Library- 8 bit runtime files.",
		Licenses:    []string{},
		CPEs: []string{
			"cpe:2.3:a:libpcre2-8-0:libpcre2-8-0:10.36-2:*:*:*:*:*:*:*",
		},
		PURL: model.PURL("pkg:deb/libpcre2-8-0@10.36-2?arch=s390x"),
		Metadata: Metadata{
			"Architecture":   "s390x",
			"Depends":        "libc6 (\u003e= 2.4)",
			"Description":    "New Perl Compatible Regular Expression Library- 8 bit runtime files.",
			"Homepage":       "https://pcre.org/",
			"Installed-Size": "440",
			"Maintainer":     "Matthew Vernon \u003cmatthew@debian.org\u003e",
			"Multi-Arch":     "same",
			"Package":        "libpcre2-8-0",
			"Priority":       "optional",
			"Section":        "libs",
			"Source":         "pcre2",
			"Status":         "install ok installed",
			"Version":        "10.36-2",
		},
	}
	debPackage2 = model.Package{
		Name:    "e2fsprogs",
		Type:    debType,
		Version: "1.46.2-2",
		Path:    dpkgStatusPath,
		Locations: []model.Location{
			{
				Path:      dpkgStatusPath,
				LayerHash: "f1a5f5ce6b163fac7f09b47645c56d2ab676bdcdb268eef06a4d9b782a75bfd0",
			},
		},
		Description: "ext2/ext3/ext4 file system utilities The ext2, ext3 and ext4 file systems are successors of the original ext (\"extended\") file system. They are the main file sys for hard disks on Debian and other Linux systems. . This package contains programs for creating, checking, and maintaining ext2/3/4-based file systems.",
		Licenses:    []string{},
		CPEs: []string{
			"cpe:2.3:a:e2fsprogs:e2fsprogs:1.46.2-2:*:*:*:*:*:*:*",
		},
		PURL: model.PURL("pkg:deb/e2fsprogs@1.46.2-2?arch=s390x"),
		Metadata: Metadata{
			"Architecture":   "s390x",
			"Depends":        "logsave",
			"Description":    "ext2/ext3/ext4 file system utilities The ext2, ext3 and ext4 file systems are successors of the original ext (\"extended\") file system. They are the main file sys for hard disks on Debian and other Linux systems. . This package contains programs for creating, checking, and maintaining ext2/3/4-based file systems.",
			"Homepage":       "http://e2fsprogs.sourceforge.net",
			"Important":      "yes",
			"Installed-Size": "1519",
			"Maintainer":     "Theodore Y. Ts'o \u003ctytso@mit.edu\u003e",
			"Multi-Arch":     "foreign",
			"Package":        "e2fsprogs",
			"Pre-Depends":    "libblkid1 (\u003e= 2.36), libc6 (\u003e= 2.11), libcom-err2 (\u003e= 1.43.9), libext2fs2 (= 1.46.2-2), libss2 (\u003e= 1.38), libuuid1 (\u003e= 2.16)",
			"Priority":       "required",
			"Recommends":     "e2fsprogs-l10n",
			"Section":        "admin",
			"Status":         "install ok installed",
			"Suggests":       "gpart, parted, fuse2fs, e2fsck-static",
			"Version":        "1.46.2-2",
		},
	}

	debPackage3 = model.Package{
		Name:    "libapt-pkg6.0",
		Type:    debType,
		Version: "2.2.4",
		Path:    dpkgStatusPath,
		Locations: []model.Location{
			{
				Path:      dpkgStatusPath,
				LayerHash: "f1a5f5ce6b163fac7f09b47645c56d2ab676bdcdb268eef06a4d9b782a75bfd0",
			},
			{
				Path:      dpkgStatusPath,
				LayerHash: "f1a5f5ce6b163fac7f09b47645c56d2ab676bdcdb268eef06a4d9b782a75bfd0",
			},
		},
		Description: "package management runtime library This library provides the common functionality for searching and managing packages as well as information about packages.",
		Licenses: []string{
			"GPLv2+",
		},
		CPEs: []string{
			"cpe:2.3:a:libapt-pkg6.0:libapt-pkg6.0:2.2.4:*:*:*:*:*:*:*",
		},
		PURL: model.PURL("pkg:deb/libapt-pkg6.0@2.2.4?arch=s390x"),
		Metadata: Metadata{
			"Architecture":   "s390x",
			"Breaks":         "appstream (\u003c\u003c 0.9.0-3~), apt (\u003c\u003c 1.6~), aptitude (\u003c\u003c 0.8.9), libapt-inst1.5 (\u003c\u003c 0.9.9~)",
			"Depends":        "libbz2-1.0, libc6 (\u003e= 2.27), libgcc-s1 (\u003e= 3.0), libgcrypt20 (\u003e= 1.8.0), liblz4-1 (\u003e= 0.0~r127), liblzma5 (\u003e= 5.1.1alpha+20120614), libstdc++6 libsystemd0 (\u003e= 221), libudev1 (\u003e= 183), libxxhash0 (\u003e= 0.7.1), libzstd1 (\u003e= 1.4.0), zlib1g (\u003e= 1:1.2.2.3)",
			"Description":    "package management runtime library This library provides the common functionality for searching and managing packages as well as information about packages.",
			"Installed-Size": "3412",
			"Maintainer":     "APT Development Team \u003cdeity@lists.debian.org\u003e",
			"Multi-Arch":     "same",
			"Package":        "libapt-pkg6.0",
			"Priority":       "optional",
			"Provides":       "libapt-pkg (= 2.2.4)",
			"Recommends":     "apt (\u003e= 2.2.4)",
			"Section":        "libs",
			"Source":         "apt",
			"Status":         "install ok installed",
			"Version":        "2.2.4",
		},
	}

	debMetadata1 = Metadata{
		"Architecture":   "s390x",
		"Depends":        "libc6 (\u003e= 2.4)",
		"Description":    "New Perl Compatible Regular Expression Library- 8 bit runtime files.",
		"Homepage":       "https://pcre.org/",
		"Installed-Size": "440",
		"Maintainer":     "Matthew Vernon \u003cmatthew@debian.org\u003e",
		"Multi-Arch":     "same",
		"Package":        "libpcre2-8-0",
		"Priority":       "optional",
		"Section":        "libs",
		"Source":         "pcre2",
		"Status":         "install ok installed",
		"Version":        "10.36-2",
	}

	debMetadata2 = Metadata{
		"Architecture":   "s390x",
		"Depends":        "logsave",
		"Description":    "ext2/ext3/ext4 file system utilities The ext2, ext3 and ext4 file systems are successors of the original ext (\"extended\") file system. They are the main file sys for hard disks on Debian and other Linux systems. . This package contains programs for creating, checking, and maintaining ext2/3/4-based file systems.",
		"Homepage":       "http://e2fsprogs.sourceforge.net",
		"Important":      "yes",
		"Installed-Size": "1519",
		"Maintainer":     "Theodore Y. Ts'o \u003ctytso@mit.edu\u003e",
		"Multi-Arch":     "foreign",
		"Package":        "e2fsprogs",
		"Pre-Depends":    "libblkid1 (\u003e= 2.36), libc6 (\u003e= 2.11), libcom-err2 (\u003e= 1.43.9), libext2fs2 (= 1.46.2-2), libss2 (\u003e= 1.38), libuuid1 (\u003e= 2.16)",
		"Priority":       "required",
		"Recommends":     "e2fsprogs-l10n",
		"Section":        "admin",
		"Status":         "install ok installed",
		"Suggests":       "gpart, parted, fuse2fs, e2fsck-static",
		"Version":        "1.46.2-2",
	}

	debMetadata3 = Metadata{
		"Architecture":   "s390x",
		"Breaks":         "appstream (\u003c\u003c 0.9.0-3~), apt (\u003c\u003c 1.6~), aptitude (\u003c\u003c 0.8.9), libapt-inst1.5 (\u003c\u003c 0.9.9~)",
		"Depends":        "libbz2-1.0, libc6 (\u003e= 2.27), libgcc-s1 (\u003e= 3.0), libgcrypt20 (\u003e= 1.8.0), liblz4-1 (\u003e= 0.0~r127), liblzma5 (\u003e= 5.1.1alpha+20120614), libstdc++6 libsystemd0 (\u003e= 221), libudev1 (\u003e= 183), libxxhash0 (\u003e= 0.7.1), libzstd1 (\u003e= 1.4.0), zlib1g (\u003e= 1:1.2.2.3)",
		"Description":    "package management runtime library This library provides the common functionality for searching and managing packages as well as information about packages.",
		"Installed-Size": "3412",
		"Maintainer":     "APT Development Team \u003cdeity@lists.debian.org\u003e",
		"Multi-Arch":     "same",
		"Package":        "libapt-pkg6.0",
		"Priority":       "optional",
		"Provides":       "libapt-pkg (= 2.2.4)",
		"Recommends":     "apt (\u003e= 2.2.4)",
		"Section":        "libs",
		"Source":         "apt",
		"Status":         "install ok installed",
		"Version":        "2.2.4",
	}

	debConffiles1 = "/etc/mke2fs.conf 72b349d890a9b5cca06c7804cd0c8d1d"
	debConffiles2 = "/etc/pam.conf 87fc76f18e98ee7d3848f6b81b3391e5 /etc/pam.d/other 31aa7f2181889ffb00b87df4126d1701"
	debLocation   = model.Location{
		Path:      filepath.Join("AppData", "Local", "Temp", "3175519915", "diggity-tmp-614678a1-5579-42fb-8e8f-0d8e2101c803", "f1a5f5ce6b163fac7f09b47645c56d2ab676bdcdb268eef06a4d9b782a75bfd0", "var", "lib", "dpkg", "status"),
		LayerHash: "f1a5f5ce6b163fac7f09b47645c56d2ab676bdcdb268eef06a4d9b782a75bfd0",
	}
)

func TestInitDebianPackage(t *testing.T) {
	var _package1, _package2, _package3 model.Package

	tests := []InitDebPackageResult{
		{&_package1, &debLocation, debMetadata1, &debPackage1},
		{&_package2, &debLocation, debMetadata2, &debPackage2},
		{&_package3, &debLocation, debMetadata3, &debPackage3},
	}

	for _, test := range tests {
		output := initDebianPackage(test._package, test.location, test.metadata)
		outputMetadata := output.Metadata.(Metadata)
		expectedMetadata := test.expected.Metadata.(Metadata)

		if output.Name != test.expected.Name ||
			output.Version != test.expected.Version ||
			output.Description != test.expected.Description ||
			len(output.CPEs) != len(test.expected.CPEs) ||
			string(output.PURL) != string(test.expected.PURL) ||
			!reflect.DeepEqual(outputMetadata, expectedMetadata) {
			t.Errorf("Test Failed:\n Expected output of %v \n, Received: %v \n", test.expected, output)
		}

		for i := range output.CPEs {
			if output.CPEs[i] != test.expected.CPEs[i] {
				t.Errorf("Test Failed:\n Expected output of %v \n, Received: %v \n", test.expected.CPEs[i], output.CPEs[i])
			}
		}
	}
}

func TestParseDebianFiles(t *testing.T) {
	_debMetadata1 := Metadata{"Conffiles": ""}
	_debMetadata2 := Metadata{"Conffiles": ""}

	tests := []ParseDebianFilesResult{
		{_debMetadata1, debConffiles1, Metadata{
			"Conffiles": []map[string]interface{}{
				{
					"digest": map[string]string{
						"algorithm": "md5",
						"value":     "72b349d890a9b5cca06c7804cd0c8d1d",
					},
					"path": "/etc/mke2fs.conf",
				},
			},
		},
		},
		{_debMetadata2, debConffiles2, Metadata{
			"Conffiles": []map[string]interface{}{
				{
					"digest": map[string]string{
						"algorithm": "md5",
						"value":     "87fc76f18e98ee7d3848f6b81b3391e5",
					},
					"path": "/etc/pam.conf",
				},
				{
					"digest": map[string]string{
						"algorithm": "md5",
						"value":     "31aa7f2181889ffb00b87df4126d1701",
					},
					"path": "/etc/pam.d/other",
				},
			},
		},
		},
	}

	for _, test := range tests {
		parseDebianFiles(test.metadata, test.content)
		outputConffiles := test.metadata["Conffiles"].([]map[string]interface{})
		expectedConffiles := test.expected["Conffiles"].([]map[string]interface{})

		if len(outputConffiles) != len(expectedConffiles) {
			t.Errorf("Test Failed: Slice length must be equal with the expected result. Expected: %v, Received: %v", len(outputConffiles), len(expectedConffiles))
		}

		for i, m := range test.metadata["Conffiles"].([]map[string]interface{}) {
			if fmt.Sprint(m) != fmt.Sprint(expectedConffiles[i]) {
				t.Errorf("Test Failed:\n Expected output of %v \n, Received: %v \n", expectedConffiles[i]["digest"], m["digest"])
			}
		}
	}
}

func TestSearchOnFileSystem(t *testing.T) {
	var _package1, _package2, _package3, _package4 model.Package
	docPath := filepath.Join("..", "..", "..", "docs", "references", "debian", "licenses", "copyright_")

	tests := []DebLicenseResult{
		{&_package1, docPath + "zlib", []string{"Zlib"}},
		{&_package2, docPath + "grep", []string{"GPL-3+"}},
		{&_package3, docPath + "libp11", []string{"ISC",
			"ISC+IBM",
			"same-as-rest-of-p11kit",
			"BSD-3-Clause",
			"permissive-like-automake-output"}},
		{&_package4, docPath + "mawk", []string{}},
	}

	for _, test := range tests {
		searchLicenseOnFileSystem(test._package, test.path)

		if len(test._package.Licenses) != len(test.expected) {
			t.Errorf("Test Failed: Slice length must be equal with the expected result. Expected: %v, Received: %v", len(test.expected), len(test._package.Licenses))
		}

		if len(test._package.Licenses) > 0 {
			sort.Slice(test._package.Licenses, func(i, j int) bool {
				return test._package.Licenses[i] < test._package.Licenses[j]
			})

			sort.Slice(test.expected, func(i, j int) bool {
				return test.expected[i] < test.expected[j]
			})

			for i, license := range test._package.Licenses {

				if license != test.expected[i] {
					t.Errorf("Test Failed:\n Expected output of %v \n, Received: %v \n", test.expected[i], license)
				}
			}
		}
	}
}

func TestParseDebianPackageUrl(t *testing.T) {
	_package1 := model.Package{
		Name:     debPackage1.Name,
		Version:  debPackage1.Version,
		Metadata: debPackage1.Metadata,
	}
	_package2 := model.Package{
		Name:     debPackage2.Name,
		Version:  debPackage2.Version,
		Metadata: debPackage2.Metadata,
	}
	_package3 := model.Package{
		Name:     debPackage3.Name,
		Version:  debPackage3.Version,
		Metadata: debPackage3.Metadata,
	}

	tests := []DebPurlResult{
		{&_package1, "s390x", model.PURL("pkg:deb/libpcre2-8-0@10.36-2?arch=s390x")},
		{&_package2, "s390x", model.PURL("pkg:deb/e2fsprogs@1.46.2-2?arch=s390x")},
		{&_package3, "s390x", model.PURL("pkg:deb/libapt-pkg6.0@2.2.4?arch=s390x")},
	}

	for _, test := range tests {
		parseDebianPackageURL(test._package, test.arch)
		if test._package.PURL != test.expected {
			t.Errorf("Test Failed: Expected output of %v, received: %v", test.expected, test._package.PURL)
		}
	}
}
