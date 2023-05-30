package alpine

import (
	"github.com/carbonetes/diggity/pkg/model"
)

var (
	apkPackage1 = model.Package{
		Name:    "ca-certificates-bundle",
		Type:    Type,
		Version: "20220614-r0",
		Path:    InstalledPackagesPath,
		Locations: []model.Location{
			{
				Path:      InstalledPackagesPath,
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
		Metadata: Metadata{
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
		Type:    Type,
		Version: "1.2.12-r3",
		Path:    InstalledPackagesPath,
		Locations: []model.Location{
			{
				Path:      InstalledPackagesPath,
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
		Metadata: Metadata{
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
		Type:    "apk",
		Version: "1.2.3-r0",
		Path:    InstalledPackagesPath,
		Locations: []model.Location{
			{
				Path:      InstalledPackagesPath,
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
		Metadata: Metadata{
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
