package alpine

import (
	"regexp"
	"testing"

	"github.com/carbonetes/diggity/pkg/model"
)

type (
	ApkPackageResult struct {
		pkg      *model.Package
		expected *model.Package
	}
	AlpineFilesResult struct {
		content  string
		expected []model.File
	}
)

var (
	apkPackage1 = model.Package{
		Name:    "ca-certificates-bundle",
		Type:    apkType,
		Version: "20220614-r0",
		Path:    installedPackagesPath,
		Locations: []model.Location{
			{
				Path:      installedPackagesPath,
				LayerHash: "9b7240956cfbfefddcd91a2195bfb2ed2cd17bdff81f21111849d643dfaf8131",
			},
		},
		Description: "Pre generated bundle of Mozilla certificates",
		Licenses: []string{
			"MPL-2.0",
			"MIT",
		},
		CPEs: []string{
			"cpe:2.3:a:ca-certificates-bundle:ca-certificates-bundle:20220614-r0:*:*:*:*:*:*:*",
			"cpe:2.3:a:ca-certificates-bundle:ca_certificates-bundle:20220614-r0:*:*:*:*:*:*:*",
			"cpe:2.3:a:ca-certificates-bundle:ca_certificates_bundle:20220614-r0:*:*:*:*:*:*:*",
			"cpe:2.3:a:ca_certificates-bundle:ca_certificates_bundle:20220614-r0:*:*:*:*:*:*:*",
			"cpe:2.3:a:ca_certificates-bundle:ca-certificates_bundle:20220614-r0:*:*:*:*:*:*:*",
			"cpe:2.3:a:ca_certificates-bundle:ca-certificates-bundle:20220614-r0:*:*:*:*:*:*:*",
			"cpe:2.3:a:ca_certificates_bundle:ca-certificates-bundle:20220614-r0:*:*:*:*:*:*:*",
			"cpe:2.3:a:ca_certificates_bundle:ca_certificates-bundle:20220614-r0:*:*:*:*:*:*:*",
			"cpe:2.3:a:ca_certificates_bundle:ca_certificates_bundle:20220614-r0:*:*:*:*:*:*:*",
		},
		PURL: model.PURL("pkg:apk/alpine/ca-certificates-bundle@20220614-r0?arch=x86_64\u0026upstream=ca-certificates\u0026distro=alpine"),
		Metadata: Manifest{
			"Architecture":         "x86_64",
			"BuildTimestamp":       "1659254961",
			"GitCommitHashApk":     "bb51fa7743320ac61f76e181cca84daa9977573e",
			"License":              "MPL-2.0 AND MIT",
			"Maintainer":           "Natanael Copa \u003cncopa@alpinelinux.org\u003e",
			"PackageDescription":   "Pre generated bundle of Mozilla certificates",
			"PackageInstalledSize": "233472",
			"PackageName":          "ca-certificates-bundle",
			"PackageOrigin":        "ca-certificates",
			"PackageSize":          "125920",
			"PackageURL":           "https://www.mozilla.org/en-US/about/governance/policies/security-group/certs/",
			"PackageVersion":       "20220614-r0",
			"Provides":             "ca-certificates-cacert=20220614-r0",
			"PullChecksum":         "Q1huqjigIP7ZNHBueDUmNnT6PpToI=",
			"PullDependencies":     "libressl2.7-libcrypto",
		},
	}
	apkPackage2 = model.Package{
		Name:    "zlib",
		Type:    apkType,
		Version: "1.2.12-r3",
		Path:    installedPackagesPath,
		Locations: []model.Location{
			{
				Path:      installedPackagesPath,
				LayerHash: "9b7240956cfbfefddcd91a2195bfb2ed2cd17bdff81f21111849d643dfaf8131",
			},
		},
		Description: "compression/decompression Library",
		Licenses: []string{
			"Zlib",
		},
		CPEs: []string{
			"cpe:2.3:a:zlib:zlib:1.2.12-r3:*:*:*:*:*:*:*",
		},
		PURL: model.PURL("pkg:apk/alpine/zlib@1.2.12-r3?arch=x86_64\u0026upstream=zlib\u0026distro=alpine"),
		Metadata: Manifest{
			"Architecture":         "x86_64",
			"BuildTimestamp":       "1660030129",
			"GitCommitHashApk":     "57ce38bde7ce42964b664c137935cf2de803ac44",
			"License":              "Zlib",
			"Maintainer":           "Natanael Copa \u003cncopa@alpinelinux.org\u003e",
			"PackageDescription":   "A compression/decompression Library",
			"PackageInstalledSize": "110592",
			"PackageName":          "zlib",
			"PackageOrigin":        "zlib",
			"PackageSize":          "53346",
			"PackageURL":           "https://zlib.net/",
			"PackageVersion":       "1.2.12-r3",
			"Provides":             "so:libz.so.1=1.2.12",
			"PullChecksum":         "Q1Ekuqm/0CPywDCKEbEwhsPCw+z9E=",
			"PullDependencies":     "so:libc.musl-x86_64.so.1",
		},
	}
	apkPackage3 = model.Package{
		Name:    "musl-utils",
		Type:    apkType,
		Version: "1.2.3-r0",
		Path:    installedPackagesPath,
		Locations: []model.Location{
			{
				Path:      installedPackagesPath,
				LayerHash: "9b7240956cfbfefddcd91a2195bfb2ed2cd17bdff81f21111849d643dfaf8131",
			},
		},
		Description: "the musl c library (libc) implementation",
		Licenses: []string{
			"MIT",
			"BSD",
			"GPL2+",
		},
		CPEs: []string{
			"cpe:2.3:a:musl-utils:musl-utils:1.2.3-r0:*:*:*:*:*:*:*",
			"cpe:2.3:a:musl-utils:musl_utils:1.2.3-r0:*:*:*:*:*:*:*",
			"cpe:2.3:a:musl_utils:musl_utils:1.2.3-r0:*:*:*:*:*:*:*",
			"cpe:2.3:a:musl_utils:musl-utils:1.2.3-r0:*:*:*:*:*:*:*",
		},
		PURL: model.PURL("pkg:apk/alpine/musl-utils@1.2.3-r0?arch=x86_64\u0026upstream=musl\u0026distro=alpine"),
		Metadata: Manifest{
			"Architecture":         "x86_64",
			"BuildTimestamp":       "1649396308",
			"GitCommitHashApk":     "ee13d43a53938d8a04ba787b9423f3270a3c14a7",
			"License":              "MIT BSD GPL2+",
			"Maintainer":           "Timo Teräs \u003ctimo.teras@iki.fi\u003e",
			"PackageDescription":   "the musl c library (libc) implementation",
			"PackageInstalledSize": "135168",
			"PackageName":          "musl-utils",
			"PackageOrigin":        "musl",
			"PackageSize":          "36938",
			"PackageURL":           "https://musl.libc.org/",
			"PackageVersion":       "1.2.3-r0",
			"Provides":             "cmd:getconf=1.2.3-r0 cmd:getent=1.2.3-r0 cmd:iconv=1.2.3-r0 cmd:ldconfig=1.2.3-r0 cmd:ldd=1.2.3-r0",
			"PullChecksum":         "Q1VVfxM3uSO0X38HWpj1LN0E61fxo=",
			"PullDependencies":     "libiconv",
		},
	}

	apkContent1 = "C:Q1OYJOhU51ILwDQRwNiYgqtq986/o=\n" +
		"P:scanelf\n" +
		"V:1.2.8-r0\n" +
		"A:x86_64\n" +
		"S:36506\n" +
		"I:94208\n" +
		"T:Scan ELF binaries for stuff\n" +
		"U:https://wiki.gentoo.org/wiki/Hardened/PaX_Utilities\n" +
		"L:GPL-2.0-only\n" +
		"o:pax-utils\n" +
		"m:Natanael Copa <ncopa@alpinelinux.org>\n" +
		"t:1608504607\n" +
		"c:375980196b4c31373fcffaf3aba1bd6b65744dc4\n" +
		"D:so:libc.musl-x86_64.so.1\n" +
		"p:cmd:scanelf\n" +
		"r:pax-utils\n" +
		"F:usr\n" +
		"F:usr/bin\n" +
		"R:scanelf\n" +
		"a:0:0:755\n" +
		"Z:Q12n+1EW8lrQAvNQsRIOBhmja3IFc="
	apkContent2 = "C:Q1JhyrSSh6Nws4iebsM6/VHCxwPfQ=\n" +
		"P:libc-utils\n" +
		"V:0.7.2-r3\n" +
		"A:x86_64\n" +
		"S:1228\n" +
		"I:4096\n" +
		"T:Meta package to pull in correct libc\n" +
		"U:https://alpinelinux.org\n" +
		"L:BSD-2-Clause AND BSD-3-Clause\n" +
		"o:libc-dev\n" +
		"m:Natanael Copa <ncopa@alpinelinux.org>\n" +
		"t:1585632275\n" +
		"c:60424133be2e79bbfeff3d58147a22886f817ce2\n" +
		"D:musl-utils"
	apkFiles = []model.File{
		{Path: "usr"},
		{Path: "usr/bin"},
		{
			Digest: map[string]string{
				"algorithm": "sha1",
				"value":     "Q12n+1EW8lrQAvNQsRIOBhmja3IFc=",
			},
			OwnerGID:    "0",
			OwnerUID:    "0",
			Path:        "scanelf",
			Permissions: "755",
		},
	}
)

func TestParseAlpineFiles(t *testing.T) {
	tests := []AlpineFilesResult{
		{apkContent1, apkFiles},
		{apkContent2, nil},
	}

	for _, test := range tests {
		output := parseAlpineFiles(test.content)
		if len(output) != len(test.expected) {
			t.Errorf("Test Failed: Expected Packages of length %+v, Received: %+v.", len(test.expected), len(output))
		}

		for i, file := range output {
			if file.OwnerGID != test.expected[i].OwnerGID ||
				file.OwnerUID != test.expected[i].OwnerUID ||
				file.Path != test.expected[i].Path ||
				file.Permissions != test.expected[i].Permissions {
				t.Errorf("Test Failed: Expected output of %v, received: %v", test.expected[i], file)
			}

			if file.Digest == nil {
				if test.expected[i].Digest != nil {
					t.Errorf("Test Failed: Expected Digest not nil.")
				}
				continue
			}

			if file.Digest.(map[string]string)["algorithm"] != test.expected[i].Digest.(map[string]string)["algorithm"] ||
				file.Digest.(map[string]string)["value"] != test.expected[i].Digest.(map[string]string)["value"] {
				t.Errorf("Test Failed: Expected output of %v, received: %v", test.expected[i], file)
			}
		}
	}
}

func TestInitAlpinePackage(t *testing.T) {
	var pkg1, pkg2, pkg3 model.Package
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")

	tests := []ApkPackageResult{
		{&pkg1, &apkPackage1},
		{&pkg2, &apkPackage2},
		{&pkg3, &apkPackage3},
	}

	for _, test := range tests {
		initAlpinePackage(test.pkg)
		if test.pkg.Metadata == nil {
			t.Error("Test Failed: Metadata must not be nil.")
		}
		if !r.MatchString(test.pkg.ID) {
			t.Errorf("Test Failed: Output of %v must be a valid UUID, received: %v", test.expected.ID, test.pkg.ID)
		}
		if test.pkg.Type != test.expected.Type ||
			test.pkg.Path != test.expected.Path {
			t.Errorf("Test Failed: Expected output of %v, received: %v", test.expected.Path, test.pkg.Path)
		}
	}
}

func TestParseAlpineURL(t *testing.T) {
	pkg1 := model.Package{
		Name:     apkPackage1.Name,
		Version:  apkPackage1.Version,
		Metadata: apkPackage1.Metadata,
	}
	pkg2 := model.Package{
		Name:     apkPackage2.Name,
		Version:  apkPackage2.Version,
		Metadata: apkPackage2.Metadata,
	}
	pkg3 := model.Package{
		Name:     apkPackage3.Name,
		Version:  apkPackage3.Version,
		Metadata: apkPackage3.Metadata,
	}

	tests := []ApkPackageResult{
		{&pkg1, &apkPackage1},
		{&pkg2, &apkPackage2},
		{&pkg3, &apkPackage3},
	}

	for _, test := range tests {
		parseAlpineURL(test.pkg)
		if test.pkg.PURL != test.expected.PURL {
			t.Errorf("Test Failed: Expected output of %v, received: %v", test.expected.PURL, test.pkg.PURL)
		}
	}
}
