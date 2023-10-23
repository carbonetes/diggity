package spdxutils

import (
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/model/metadata"
	"github.com/carbonetes/diggity/pkg/parser/os/apk"
	"github.com/carbonetes/diggity/pkg/parser/os/dpkg"
	"github.com/carbonetes/diggity/pkg/parser/language/gem"
	"github.com/carbonetes/diggity/pkg/parser/language/java/maven"
	spdx22 "github.com/spdx/tools-golang/spdx/v2_2"
)

type (
	StringParserResult struct {
		input    string
		expected string
	}
	PackageParserResult struct {
		pkg      *model.Package
		expected string
	}
	ExternalRefsResult struct {
		pkg      *model.Package
		expected []spdx22.PackageExternalReference
	}
)

var (
	package1 = model.Package{
		ID:      "8fe93afb-86f2-4639-a3eb-6c4e787f210b",
		Name:    "lzo",
		Type:    "rpm",
		Version: "2.08",
		Path:    filepath.Join("var", "lib", "rpm", "Packages"),
		Locations: []model.Location{
			{
				Path:      filepath.Join("var", "lib", "rpm", "Packages"),
				LayerHash: "69a15d957a7a6f77e3fe31f330da5f4b6b582f228917a713a7a9e59449a3f413",
			},
		},
		Description: "Data compression library with very fast (de)compression",
		Licenses: []string{
			"GPLv2+",
		},
		CPEs: []string{
			"cpe:2.3:a:centos:lzo:2.08-14.el8:*:*:*:*:*:*:*",
			"cpe:2.3:a:lzo:lzo:2.08-14.el8:*:*:*:*:*:*:*",
		},
		PURL: model.PURL("pkg:rpm/lzo@2.08arch=x86_64"),
		Metadata: metadata.RPMMetadata{
			Release:      "14.el8",
			Architecture: "x86_64",
			SourceRpm:    "lzo-2.08-14.el8.src.rpm",
			License:      "GPLv2+",
			Size:         198757,
			Name:         "lzo",
			PGP:          "RSA/SHA256, Tue Jul  2 00:01:31 2019, Key ID 05b555b38483c65d",
			Summary:      "Data compression library with very fast (de)compression",
			Vendor:       "CentOS",
			Version:      "2.08",
		},
	}
	package2 = model.Package{
		ID:      "9583e9ec-df1d-484a-b560-8e1415ea92c2",
		Name:    "gitlab.com/yawning/obfs4.git",
		Type:    "go-module",
		Version: "v0.0.0-20220204003609-77af0cba934d",
		Path:    "",
		Locations: []model.Location{
			{
				Path:      filepath.Join("bin", "gost"),
				LayerHash: "0cd4836a36e094e1870a2e6c2578a7ad9d9cb42a7313944a6d05ab72892fc3c3",
			},
		},
		Description: "",
		Licenses:    []string{},
		CPEs: []string{
			"cpe:2.3:a:yawning:obfs4.git:v0.0.0-20220204003609-77af0cba934d:*:*:*:*:*:*:*",
			"cpe:2.3:a:obfs4.git:obfs4.git:v0.0.0-20220204003609-77af0cba934d:*:*:*:*:*:*:*",
		},
		PURL: model.PURL("pkg:go/gitlab.com/yawning/obfs4.git@v0.0.0-20220204003609-77af0cba934d"),
		Metadata: metadata.GoBinMetadata{
			Architecture:     "amd64",
			Compiler:         "gc",
			OS:               "linux",
			GoCompileRelease: "go1.19.1",
			H1Digest:         "h1:tJ8F7ABaQ3p3wjxwXiWSktVDgjZEXkvaRawd2rIq5ws=",
			Path:             "gitlab.com/yawning/obfs4.git",
			Version:          "v0.0.0-20220204003609-77af0cba934d",
		},
	}
	package3 = model.Package{
		ID:      "bdbd600f-dbdf-49a1-a329-a339f1123ffd",
		Name:    "scanelf",
		Type:    "apk",
		Version: "1.3.4-r0",
		Path:    filepath.Join("lib", "apk", "db", "installed"),
		Locations: []model.Location{
			{
				Path:      filepath.Join("lib", "apk", "db", "installed"),
				LayerHash: "1288696addccc4013c5bcf61c1b6c38128a7214a0942976792918b51912d90f7",
			},
			{
				Path:      filepath.Join("lib", "apk", "db", "installed"),
				LayerHash: "1288696addccc4013c5bcf61c1b6c38128a7214a0942976792918b51912d90f7",
			},
		},
		Description: "Scan ELF binaries for stuff",
		Licenses: []string{
			"GPL-2.0-only",
		},
		CPEs: []string{
			"cpe:2.3:a:scanelf:scanelf:1.3.4-r0:*:*:*:*:*:*:*",
		},
		PURL: model.PURL("pkg:alpine/scanelf@1.3.4-r0?arch=x86_64\u0026upstream=pax-utils\u0026distro="),
		Metadata: apk.Metadata{
			"Architecture":         "x86_64",
			"BuildTimestamp":       "1651005390",
			"GitCommitHashApk":     "d7ae612a3cc5f827289d915783b4cbf8c7207947",
			"License":              "GPL-2.0-only",
			"Maintainer":           "Natanael Copa \u003cncopa@alpinelinux.org\u003e",
			"PackageDescription":   "Scan ELF binaries for stuff",
			"PackageInstalledSize": "94208",
			"PackageName":          "scanelf",
			"PackageOrigin":        "pax-utils",
			"PackageSize":          "36745",
			"PackageURL":           "https://wiki.gentoo.org/wiki/Hardened/PaX_Utilities",
			"PackageVersion":       "1.3.4-r0",
			"Provides":             "cmd:scanelf=1.3.4-r0",
			"PullChecksum":         "Q1Gcqe+ND8DFOlhM3R0o5KyZjR2oE=",
			"PullDependencies":     "pax-utils",
		},
	}
	package4 = model.Package{
		ID:      "2180bf4a-def7-471a-bdd0-7db3e82cf2ca",
		Name:    "libgpg-error0",
		Type:    "deb",
		Version: "1.38-2",
		Path:    "",
		Locations: []model.Location{
			{
				Path:      filepath.Join("var", "lib", "dpkg", "status"),
				LayerHash: "f1a5f5ce6b163fac7f09b47645c56d2ab676bdcdb268eef06a4d9b782a75bfd0",
			},
			{
				Path:      filepath.Join("var", "lib", "dpkg", "status"),
				LayerHash: "f1a5f5ce6b163fac7f09b47645c56d2ab676bdcdb268eef06a4d9b782a75bfd0",
			},
		},
		Description: "",
		Licenses: []string{
			"BSD-3-clause",
			"LGPL-2.1+",
			"LGPL-2.1+ or BSD-3-clause",
			"g10-permissive",
			"GPL-3+",
		},
		PURL: model.PURL("pkg:deb/libgpg-error0@1.38-2arch=s390x"),
		Metadata: dpkg.Metadata{
			"Architecture":   "s390x",
			"Breaks":         "libxml2 (\u003c\u003c 2.7.6.dfsg-2), texlive-binaries (\u003c\u003c 2009-12)",
			"Conflicts":      "zlib1 (\u003c= 1:1.0.4-7)",
			"Depends":        "libc6 (\u003e= 2.4)",
			"Description":    "compression library - runtime zlib is a library implementing the deflate compression method found in gzip and PKZIP.  This package includes the shared library.",
			"Essential":      "yes",
			"Homepage":       "http://zlib.net/",
			"Important":      "yes",
			"Installed-Size": "170",
			"Maintainer":     "Mark Brown \u003cbroonie@debian.org\u003e",
			"Multi-Arch":     "same",
			"Package":        "zlib1g",
			"Pre-Depends":    "libaudit1 (\u003e= 1:2.2.1), libblkid1 (\u003e= 2.31.1), libc6 (\u003e= 2.25), libcap-ng0 (\u003e= 0.7.9), libcrypt1 (\u003e= 1:4.1.0), libmount1 (\u003e= 2.34), libpam0g (\u003e= 0.99.7.1), libselinux1 (\u003e= 3.1~), libsmartcols1 (\u003e= 2.34), libsystemd0, libtinfo6 (\u003e= 6), libudev1 (\u003e= 183), libuuid1 (\u003e= 2.16), zlib1g (\u003e= 1:1.1.4)",
			"Priority":       "optional",
			"Protected":      "yes",
			"Provides":       "libz1",
			"Recommends":     "uuid-runtime",
			"Replaces":       "bash-completion (\u003c\u003c 1:2.8), initscripts (\u003c\u003c 2.88dsf-59.2~), login (\u003c\u003c 1:4.5-1.1~), mount (\u003c\u003c 2.29.2-3~), s390-tools (\u003c\u003c 2.2.0-1~), setpriv (\u003c\u003c 2.32.1-0.2~), sysvinit-utils (\u003c\u003c 2.88dsf-59.1~)",
			"Section":        "libs",
			"Source":         "zlib",
			"Status":         "install ok installed",
			"Suggests":       "dosfstools, kbd, util-linux-locales",
			"Version":        "1:1.2.11.dfsg-2+deb11u1",
		},
	}
	package5 = model.Package{
		ID:      "93fe248b-629e-4e6c-841f-20ed4a90bd8f",
		Name:    "jnr-unixsocke",
		Type:    "java",
		Version: "0.18",
		Path:    filepath.Join("BOOT-INF", "lib", "jnr-unixsocket-0.18.jar"),
		Locations: []model.Location{
			{
				Path:      "app.jar",
				LayerHash: "f8e33725a61b2bb436a502a4357f3f499ce10505732ffbddfe1fb7023b12ef4f",
			},
		},
		Description: "Native I/O access for java",
		Licenses: []string{
			"http://www.apache.org/licenses/LICENSE-2.0.txt",
		},
		CPEs: []string{
			"cpe:2.3:a:jnr-unixsocket:jnr-unixsocket:0.18:*:*:*:*:*:*:*",
			"cpe:2.3:a:jnr-unixsocket:jnr_unixsocket:0.18:*:*:*:*:*:*:*",
			"cpe:2.3:a:jnr_unixsocket:jnr_unixsocket:0.18:*:*:*:*:*:*:*",
			"cpe:2.3:a:jnr_unixsocket:jnr-unixsocket:0.18:*:*:*:*:*:*:*",
			"cpe:2.3:a:github:jnr-unixsocket:0.18:*:*:*:*:*:*:*",
			"cpe:2.3:a:github:jnr_unixsocket:0.18:*:*:*:*:*:*:*",
			"cpe:2.3:a:jnr:jnr-unixsocket:0.18:*:*:*:*:*:*:*",
			"cpe:2.3:a:jnr:jnr_unixsocket:0.18:*:*:*:*:*:*:*",
			"cpe:2.3:a:unixsocket:jnr-unixsocket:0.18:*:*:*:*:*:*:*",
			"cpe:2.3:a:unixsocket:jnr_unixsocket:0.18:*:*:*:*:*:*:*",
		},
		PURL: model.PURL("pkg:maven/com.github.jnr/jnr-unixsocket@0.18"),
		Metadata: maven.Metadata{
			"Manifest": maven.Manifest{
				"Archiver-Version":       "Plexus Archiver",
				"Bnd-LastModified":       "1489175442092",
				"Build-Jdk":              "1.8.0_121",
				"Built-By":               "enebo",
				"Bundle-Description":     "Native I/O access for java",
				"Bundle-License":         "http://www.apache.org/licenses/LICENSE-2.0.txt",
				"Bundle-ManifestVersion": "2",
				"Bundle-Name":            "jnr-unixsocket",
				"Bundle-SymbolicName":    "com.github.jnr.unixsocket",
				"Bundle-Version":         "0.18.0",
				"Created-By":             "Apache Maven Bundle Plugin",
				"Export-Package":         "jnr.enxio.channels;version=\"0.16.0\",jnr.unixsocket;version=\"0.18.0\"",
				"Import-Package":         "com.kenai.jffi;version=\"[1.2,2)\",jnr.constants.platform;version=\"[0.9,1)\",jnr.ffi;version=\"[2.1,3)\",jnr.ffi.annotations;version=\"[2.1,3)\",jnr.ffi.byref;version=\"[2.1,3)\",jnr.ffi.mapper;version=\"[2.1,3)\",jnr.ffi.provider.converters;version=\"[2.1,3)\",jnr.ffi.provider.jffi;version=\"[2.1,3)\",jnr.ffi.types;version=\"[2.1,3)\",jnr.posix;version=\"[3.0,4)\"",
				"Manifest-Version":       "1.0",
				"Tool":                   "Bnd-1.50.0",
			},
			"ManifestLocation": map[string]string{
				"path": "BOOT-INF\\lib\\jnr-unixsocket-0.18.jar\\META-INF\\MANIFEST.MF",
			},
			"PomProperties": map[string]string{
				"artifactId": "jnr-unixsocket",
				"groupId":    "com.github.jnr",
				"location":   "META-INF\\maven\\com.github.jnr\\jnr-unixsocket\\pom.properties",
				"name":       "",
				"version":    "0.18",
			},
		},
	}
	package6 = model.Package{
		ID:      "418ee75b-cb1a-4abe-aad6-d757c7a91610",
		Name:    "scanf",
		Type:    "gem",
		Version: "1.0.0",
		Path:    "",
		Locations: []model.Location{
			{
				Path:      filepath.Join("usr", "share", "gems", "specifications", "default", "scanf-1.0.0.gemspec"),
				LayerHash: "a67d9e51873dfbda0e6af0f9971ccea211405916ede446f52b5e7f3ea9d71fc3",
			},
		},
		Description: "scanf is an implementation of the C function scanf(3).",
		Licenses: []string{
			"BSD2Clause",
		},
		CPEs: []string{
			"cpe:2.3:a:scanf:scanf:1.0.0:*:*:*:*:*:*:*",
		},
		PURL: model.PURL("pkg:gem/scanf@1.0.0"),
		Metadata: gem.Metadata{
			"authors":     []string{"David Alan Black"},
			"bindir":      "exe",
			"date":        "2017-12-11",
			"description": "scanf is an implementation of the C function scanf(3).",
			"email":       "[dblack@superlink.net]",
			"files":       []string{"scanfrb"},
			"homepage":    "https://github.com/ruby/scanf",
			"licenses": []string{
				"BSD2Clause",
			},
			"name":                      "scanf",
			"require_paths":             "[lib]",
			"required_ruby_version":     "Gem::Requirement.new(\u003e= 2.3.0)",
			"required_rubygems_version": "Gem::Requirement.new(\u003e= 0) if s.respond_to? :required_rubygems_version=",
			"rubygems_version":          "2.7.6.2",
			"specification_version":     "4",
			"summary":                   "scanf is an implementation of the C function scanf(3).  if s.respond_to? :specification_version then",
			"version":                   "1.0.0",
		},
	}
	package7 = model.Package{
		ID:      "d3d2bb5f-8ef2-40cf-b831-43f9701beb21",
		Name:    "buffer-shims",
		Type:    "npm",
		Version: "1.0.0",
		Path:    "",
		Locations: []model.Location{
			{
				Path:      filepath.Join("node", "node-v4.5.0-linux-x64", "lib", "node_modules", "npm", "node_modules", "readable-stream", "node_modules", "buffer-shims", "package.json"),
				LayerHash: "43aa82be2d1e2fc2f7e553f63aa0db67e0cb6852cf953a84bb95f2ef3cc7dc38",
			},
		},
		Description: "some shims for node buffers",
		Licenses: []string{
			"MIT",
		},
		CPEs: []string{
			"cpe:2.3:a:buffer-shims:buffer-shims:1.0.0:*:*:*:*:*:*:*",
			"cpe:2.3:a:buffer-shims:buffer_shims:1.0.0:*:*:*:*:*:*:*",
			"cpe:2.3:a:buffer_shims:buffer_shims:1.0.0:*:*:*:*:*:*:*",
			"cpe:2.3:a:buffer_shims:buffer-shims:1.0.0:*:*:*:*:*:*:*",
		},
		PURL: model.PURL("pkg:npm/buffer-shims@1.0.0"),
		Metadata: metadata.PackageJSON{
			Version: "1.0.0",
			License: "MIT",
			Name:    "buffer-shims",
			Contributors: []interface{}{
				"Woong Jun <woong.jun@gmail.com> (http://code.woong.org/)",
			},
			Homepage:    "https://github.com/calvinmetcalf/buffer-shims#readme",
			Description: "some shims for node buffers",
			Repository: map[string]interface{}{
				"type": "git",
				"url":  "git+ssh://git@github.com/calvinmetcalf/buffer-shims.git",
			},
			Author: map[string]interface{}{
				"name":  "Woong Jun",
				"email": "<woong.jun@gmail.com>",
				"url":   "(http://code.woong.org/)",
			},
		},
	}
)

func TestExternalRefs(t *testing.T) {
	tests := []ExternalRefsResult{
		{&package1, []spdx22.PackageExternalReference{
			{
				Category: security,
				Locator:  "cpe:2.3:a:centos:lzo:2.08-14.el8:*:*:*:*:*:*:*",
				RefType:  cpeType,
			},
			{
				Category: security,
				Locator:  "cpe:2.3:a:lzo:lzo:2.08-14.el8:*:*:*:*:*:*:*",
				RefType:  cpeType,
			},
			{
				Category: packageManager,
				Locator:  "pkg:rpm/lzo@2.08arch=x86_64",
				RefType:  purlType,
			},
		}},
		{&package2, []spdx22.PackageExternalReference{
			{
				Category: security,
				Locator:  "cpe:2.3:a:yawning:obfs4.git:v0.0.0-20220204003609-77af0cba934d:*:*:*:*:*:*:*",
				RefType:  cpeType,
			},
			{
				Category: security,
				Locator:  "cpe:2.3:a:obfs4.git:obfs4.git:v0.0.0-20220204003609-77af0cba934d:*:*:*:*:*:*:*",
				RefType:  cpeType,
			},
			{
				Category: packageManager,
				Locator:  "pkg:go/gitlab.com/yawning/obfs4.git@v0.0.0-20220204003609-77af0cba934d",
				RefType:  purlType,
			},
		}},
		{&package3, []spdx22.PackageExternalReference{
			{
				Category: security,
				Locator:  "cpe:2.3:a:scanelf:scanelf:1.3.4-r0:*:*:*:*:*:*:*",
				RefType:  cpeType,
			},
			{
				Category: packageManager,
				Locator:  "pkg:alpine/scanelf@1.3.4-r0?arch=x86_64\u0026upstream=pax-utils\u0026distro=",
				RefType:  purlType,
			},
		}},
		{&package4, []spdx22.PackageExternalReference{
			{
				Category: packageManager,
				Locator:  "pkg:deb/libgpg-error0@1.38-2arch=s390x",
				RefType:  purlType,
			},
		}},
		{&package5, []spdx22.PackageExternalReference{
			{
				Category: security,
				Locator:  "cpe:2.3:a:jnr-unixsocket:jnr-unixsocket:0.18:*:*:*:*:*:*:*",
				RefType:  cpeType,
			},
			{
				Category: security,
				Locator:  "cpe:2.3:a:jnr-unixsocket:jnr_unixsocket:0.18:*:*:*:*:*:*:*",
				RefType:  cpeType,
			},
			{
				Category: security,
				Locator:  "cpe:2.3:a:jnr_unixsocket:jnr_unixsocket:0.18:*:*:*:*:*:*:*",
				RefType:  cpeType,
			},
			{
				Category: security,
				Locator:  "cpe:2.3:a:jnr_unixsocket:jnr-unixsocket:0.18:*:*:*:*:*:*:*",
				RefType:  cpeType,
			},
			{
				Category: security,
				Locator:  "cpe:2.3:a:github:jnr-unixsocket:0.18:*:*:*:*:*:*:*",
				RefType:  cpeType,
			},
			{
				Category: security,
				Locator:  "cpe:2.3:a:github:jnr_unixsocket:0.18:*:*:*:*:*:*:*",
				RefType:  cpeType,
			},
			{
				Category: security,
				Locator:  "cpe:2.3:a:jnr:jnr-unixsocket:0.18:*:*:*:*:*:*:*",
				RefType:  cpeType,
			},
			{
				Category: security,
				Locator:  "cpe:2.3:a:jnr:jnr_unixsocket:0.18:*:*:*:*:*:*:*",
				RefType:  cpeType,
			},
			{
				Category: security,
				Locator:  "cpe:2.3:a:unixsocket:jnr-unixsocket:0.18:*:*:*:*:*:*:*",
				RefType:  cpeType,
			},
			{
				Category: security,
				Locator:  "cpe:2.3:a:unixsocket:jnr_unixsocket:0.18:*:*:*:*:*:*:*",
				RefType:  cpeType,
			},
			{
				Category: packageManager,
				Locator:  "pkg:maven/com.github.jnr/jnr-unixsocket@0.18",
				RefType:  purlType,
			},
		}},
		{&package6, []spdx22.PackageExternalReference{
			{
				Category: security,
				Locator:  "cpe:2.3:a:scanf:scanf:1.0.0:*:*:*:*:*:*:*",
				RefType:  cpeType,
			},
			{
				Category: packageManager,
				Locator:  "pkg:gem/scanf@1.0.0",
				RefType:  purlType,
			},
		}},
		{&package7, []spdx22.PackageExternalReference{
			{
				Category: security,
				Locator:  "cpe:2.3:a:buffer-shims:buffer-shims:1.0.0:*:*:*:*:*:*:*",
				RefType:  cpeType,
			},
			{
				Category: security,
				Locator:  "cpe:2.3:a:buffer-shims:buffer_shims:1.0.0:*:*:*:*:*:*:*",
				RefType:  cpeType,
			},
			{
				Category: security,
				Locator:  "cpe:2.3:a:buffer_shims:buffer_shims:1.0.0:*:*:*:*:*:*:*",
				RefType:  cpeType,
			},
			{
				Category: security,
				Locator:  "cpe:2.3:a:buffer_shims:buffer-shims:1.0.0:*:*:*:*:*:*:*",
				RefType:  cpeType,
			},
			{
				Category: packageManager,
				Locator:  "pkg:npm/buffer-shims@1.0.0",
				RefType:  purlType,
			},
		}},
	}

	for _, test := range tests {
		output := ExternalRefs(test.pkg)
		if len(output) != len(test.expected) {
			t.Errorf("Test Failed: Input %v must have an output of %v, received: %v", test.pkg, test.expected, output)
		}
		if len(output) <= 0 {
			return
		}
		for i := range output {
			if output[i].Category != test.expected[i].Category ||
				output[i].Locator != test.expected[i].Locator ||
				output[i].RefType != test.expected[i].RefType {
				t.Errorf("Test Failed: Expected output of %v, received: %v", test.expected[i], output[i])
			}
		}
	}
}
func TestHomepage(t *testing.T) {
	tests := []PackageParserResult{
		{&package1, ""},
		{&package2, ""},
		{&package3, ""},
		{&package4, ""},
		{&package5, ""},
		{&package6, "https://github.com/ruby/scanf"},
		{&package7, "https://github.com/calvinmetcalf/buffer-shims#readme"},
	}

	for _, test := range tests {
		if output := Homepage(test.pkg); output != test.expected {
			t.Errorf("Test Failed: Expected output of %v, received: %v", test.expected, output)
		}
	}
}

func TestLicensesDeclared(t *testing.T) {
	tests := []PackageParserResult{
		{&package1, "NOASSERTION"},
		{&package2, "NONE"},
		{&package3, "GPL-2.0-only"},
		{&package4, "BSD-3-clause AND LGPL-2.1+ AND GPL-3+"},
		{&package5, "NOASSERTION"},
		{&package6, "NOASSERTION"},
		{&package7, "MIT"},
	}

	for _, test := range tests {
		if output := LicensesDeclared(test.pkg); output != test.expected {
			t.Errorf("Test Failed: Expected output of %v, received: %v", test.expected, output)
		}
	}
}

func TestSourceInfo(t *testing.T) {
	tests := []PackageParserResult{
		{&package1, "Information parsed from RPM DB: var/lib/rpm/Packages"},
		{&package2, "Information parsed from go-module information: bin/gost"},
		{&package3, "Information parsed from APK DB: lib/apk/db/installed, lib/apk/db/installed"},
		{&package4, "Information parsed from DPKG DB: var/lib/dpkg/status, var/lib/dpkg/status"},
		{&package5, "Information parsed from java archive: app.jar"},
		{&package6, "Information parsed from gem metadata: usr/share/gems/specifications/default/scanf-1.0.0.gemspec"},
		{&package7, "Information parsed from node module manifest: node/node-v4.5.0-linux-x64/lib/node_modules/npm/node_modules/readable-stream/node_modules/buffer-shims/package.json"},
	}

	for _, test := range tests {
		if output := SourceInfo(test.pkg); output != test.expected {
			t.Errorf("Test Failed: Expected output of %v, received: %v", test.expected, output)
		}
	}
}

func TestDownloadLocation(t *testing.T) {
	tests := []PackageParserResult{
		{&package1, "NOASSERTION"},
		{&package2, "NOASSERTION"},
		{&package3, "https://wiki.gentoo.org/wiki/Hardened/PaX_Utilities"},
		{&package4, "NOASSERTION"},
		{&package5, "NOASSERTION"},
		{&package6, "NOASSERTION"},
		{&package7, "git+ssh://git@github.com/calvinmetcalf/buffer-shims.git"},
	}

	for _, test := range tests {
		if output := DownloadLocation(test.pkg); output != test.expected {
			t.Errorf("Test Failed: Expected output of %v, received: %v", test.expected, output)
		}
	}
}
func TestOriginator(t *testing.T) {
	tests := []PackageParserResult{
		{&package1, "Organization:CentOS"},
		{&package2, ":"},
		{&package3, "Person:Natanael Copa \u003cncopa@alpinelinux.org\u003e"},
		{&package4, "Person:Mark Brown \u003cbroonie@debian.org\u003e"},
		{&package5, ":"},
		{&package6, "Person:David Alan Black"},
		{&package7, "Person:Woong Jun <woong.jun@gmail.com>"},
	}

	for _, test := range tests {
		outputType, outputName := Originator(test.pkg)
		expected := strings.Split(test.expected, ":")
		if outputType != expected[0] || outputName != expected[1] {
			t.Errorf("Test Failed: Expected output of %v, received: %v:%v", test.expected, outputType, outputName)
		}
	}
}

func TestFormatNamespace(t *testing.T) {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	images := []string{"bom", "smartentry/centos", "s390x/debian", "furynix/fedora", "gost/go", "test", "app"}

	for _, image := range images {
		output := FormatNamespace(image)
		uuid := strings.Split(output, image+"-")[1]

		if !r.MatchString(uuid) {
			t.Errorf("Test Failed: Output of %v must contain a valid UUID, received: %v", image, output)
		}
		if !strings.Contains(output, namespace+image) {
			t.Errorf("Test Failed: Output of %v must contain %v, received: %v", image, namespace+image, output)
		}
	}
}

func TestFormatPath(t *testing.T) {
	tests := []StringParserResult{
		{filepath.Join("var", "lib", "rpm", "Packages"), "var/lib/rpm/Packages"},
		{filepath.Join("var", "lib", "dpkg", "status"), "var/lib/dpkg/status"},
		{filepath.Join("usr", "share", "doc"), "usr/share/doc"},
		{filepath.Join("test1", "test2", "test3"), "test1/test2/test3"},
		{filepath.Join("app"), "app"},
	}
	for _, test := range tests {
		if output := FormatPath(test.input); output != test.expected {
			t.Errorf("Test Failed: Input %v must have an output of %v, received: %v", test.input, test.expected, output)
		}
	}
}

func TestFormatTagID(t *testing.T) {
	tests := []PackageParserResult{
		{&package1, "SPDXRef-Package-rpm-lzo-8fe93afb-86f2-4639-a3eb-6c4e787f210b"},
		{&package2, "SPDXRef-Package-go-module-gitlab.com/yawning/obfs4.git-9583e9ec-df1d-484a-b560-8e1415ea92c2"},
		{&package3, "SPDXRef-Package-apk-scanelf-bdbd600f-dbdf-49a1-a329-a339f1123ffd"},
		{&package4, "SPDXRef-Package-deb-libgpg-error0-2180bf4a-def7-471a-bdd0-7db3e82cf2ca"},
		{&package5, "SPDXRef-Package-java-jnr-unixsocke-93fe248b-629e-4e6c-841f-20ed4a90bd8f"},
		{&package6, "SPDXRef-Package-gem-scanf-418ee75b-cb1a-4abe-aad6-d757c7a91610"},
		{&package7, "SPDXRef-Package-npm-buffer-shims-d3d2bb5f-8ef2-40cf-b831-43f9701beb21"},
	}

	for _, test := range tests {
		if output := FormatTagID(test.pkg); output != test.expected {
			t.Errorf("Test Failed: Expected output of %v, received: %v", test.expected, output)
		}
	}
}

func TestCheckLicense(t *testing.T) {
	tests := []StringParserResult{
		{"0BSD", "0BSD"},
		{"BSD-3-Clause", "BSD-3-Clause"},
		{"MIT", "MIT"},
		{"ZPL-1.1", "ZPL-1.1"},
		{"ZLIB", "Zlib"},
		{"TEST-NOT-EXISTING", ""},
		{"   ", ""},
		{"", ""},
	}

	for _, test := range tests {
		if output := CheckLicense(test.input); output != test.expected {
			t.Errorf("Test Failed: Input %v must have an output of %v, received: %v", test.input, test.expected, output)
		}
	}
}

func TestFormatAuthor(t *testing.T) {
	tests := []StringParserResult{
		{"Test Author email@test.com http://test.com/test/", "Test Author email@test.com"},
		{"Test Author email@test.com", "Test Author email@test.com"},
		{"Woong Jun <woong.jun@gmail.com> (http://code.woong.org/)", "Woong Jun <woong.jun@gmail.com>"},
		{"Woong Jun <woong.jun@gmail.com>", "Woong Jun <woong.jun@gmail.com>"},
		{"Test Author", "Test Author"},
		{"Woong Jun", "Woong Jun"},
		{"Test", "Test"},
		{"   ", ""},
		{"", ""},
	}

	for _, test := range tests {
		if output := FormatAuthor(test.input); output != test.expected {
			t.Errorf("Test Failed: Input %v must have an output of %v, received: %v", test.input, test.expected, output)
		}
	}
}
