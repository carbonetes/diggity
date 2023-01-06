package parser

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/carbonetes/diggity/internal/model"
)

type (
	GenerateAdditionalCPEResult struct {
		vendor   string
		version  string
		product  string
		_package *model.Package
		expected []string
	}
	JavaLicensesResult struct {
		_package *model.Package
		licenses string
		expected []string
	}
	JavaPURLResult struct {
		_package *model.Package
		expected model.PURL
	}
	PomPropertiesResult struct {
		data     string
		_package *model.Package
		path     string
		expected interface{}
	}

	formatVersionMetadataResult struct {
		_package            *model.Package
		version             string
		expectedVersion     string
		expectedMetadataKey string
		expectedMetadataVal string
	}
)

var (
	javaPackages1 = model.Package{
		Name:    "ezmorph",
		Type:    "java",
		Version: "1.0.6",
		Path:    "WEB-INF/lib/ezmorph-1.0.6.jar",
		Locations: []model.Location{
			{
				Path:      filepath.Join("jenkins.war/WEB-INF/lib/ezmorph-1.0.6.jar"),
				LayerHash: "94f6d6f3b7ea362bc836300cf87d510f5fae45139a68c847f51a31661345fa59",
			},
		},
		Description: "",
		Licenses: []string{
			"",
		},
		CPEs: []string{
			"cpe:2.3:a:ezmorph:ezmorph:1.0.6:*:*:*:*:*:*:*",
			"cpe:2.3:a:sf:ezmorph:1.0.6:*:*:*:*:*:*:*",
		},
		PURL: model.PURL("pkg:maven/net.sf.ezmorph/ezmorph@1.0.6"),
		Metadata: JavaMetadata{
			"Manifest": {
				"Archiver-Version": "Plexus Archiver",
				"Build-Jdk":        "1.6.0_11",
				"Built-By":         "aalmiray",
				"Created-By":       "Apache Maven",
				"Manifest-Version": "1.0",
			},
			"ManifestLocation": {
				"path": "WEB-INF\\lib\\ezmorph-1.0.6.jar\\META-INF\\MANIFEST.MF",
			},
			"PomProperties": {
				"artifactId": "ezmorph",
				"groupId":    "net.sf.ezmorph",
				"location":   "META-INF\\maven\\net.sf.ezmorph\\ezmorph\\pom.properties",
				"name":       "",
				"version":    "1.0.6",
			},
		},
	}
	javaPackages2 = model.Package{
		Name:    "mxparser",
		Type:    "java",
		Version: "1.2.2",
		Path:    "WEB-INF/lib/mxparser-1.2.2.jar",
		Locations: []model.Location{
			{
				Path:      filepath.Join("jenkins.war/WEB-INF/lib/mxparser-1.2.2.jar"),
				LayerHash: "94f6d6f3b7ea362bc836300cf87d510f5fae45139a68c847f51a31661345fa59",
			},
		},
		Description: "MXParser is a fork of xpp3_min 1.1.7 containing only the parser with merged changes of the Plexus fork.",
		Licenses: []string{
			"Indiana University Extreme! Lab Software License",
		},
		CPEs: []string{
			"cpe:2.3:a:mxparser:mxparser:1.2.2:*:*:*:*:*:*:*",
			"cpe:2.3:a:github:mxparser:1.2.2:*:*:*:*:*:*:*",
			"cpe:2.3:a:x-stream:mxparser:1.2.2:*:*:*:*:*:*:*",
			"cpe:2.3:a:x_stream:mxparser:1.2.2:*:*:*:*:*:*:*",
			"cpe:2.3:a:io.github.xstream.mxparser:mxparser:1.2.2:*:*:*:*:*:*:*",
			"cpe:2.3:a:io:mxparser:1.2.2:*:*:*:*:*:*:*",
			"cpe:2.3:a:xstream:mxparser:1.2.2:*:*:*:*:*:*:*",
		},
		PURL: model.PURL("pkg:maven/io.github.x-stream/mxparser@1.2.2"),
		Metadata: JavaMetadata{
			"Manifest": {
				"Automatic-Module-Name":  "io.github.xstream.mxparser",
				"Bnd-LastModified":       "1629326138054",
				"Build-Jdk":              "1.8.0_265",
				"Build-Jdk-Spec":         "1.8",
				"Built-By":               "joehni",
				"Bundle-Description":     "MXParser is a fork of xpp3_min 1.1.7 containing only the parser with merged changes of the Plexus fork.",
				"Bundle-License":         "Indiana University Extreme! Lab Software License",
				"Bundle-ManifestVersion": "2",
				"Bundle-Name":            "MXParser",
				"Bundle-SymbolicName":    "mxparser",
				"Bundle-Version":         "1.2.2",
				"Created-By":             "Apache Maven Bundle Plugin",
				"Export-Package":         "io.github.xstream.mxparser;version=\"1.2.2\"",
				"Implementation-Title":   "MXParser",
				"Implementation-Version": "1.2.2",
				"Import-Package":         "org.xmlpull.v1",
				"JAVA_1_4_HOME":          "/opt/blackdown-jdk-1.4.2.03",
				"JAVA_1_5_HOME":          "/opt/sun-jdk-1.5.0.22",
				"JAVA_1_6_HOME":          "/opt/sun-jdk-1.6.0.45",
				"JAVA_1_7_HOME":          "/opt/oracle-jdk-bin-1.7.0.80",
				"JAVA_1_8_HOME":          "/opt/oracle-jdk-bin-1.8.0.202",
				"JAVA_9_HOME":            "/opt/oracle-jdk-bin-9.0.4",
				"Manifest-Version":       "1.0",
				"Specification-Title":    "MXParser",
				"Specification-Version":  "1.2",
				"Tool":                   "Bnd-1.50.0",
				"X-Build-Os":             "Linux",
				"X-Build-Time":           "2021-08-18T22:35:34Z",
				"X-Builder":              "Maven 3.8.1",
				"X-Compile-Source":       "1.4",
				"X-Compile-Target":       "1.4",
			},
			"ManifestLocation": {
				"path": "WEB-INF\\lib\\mxparser-1.2.2.jar\\META-INF\\MANIFEST.MF",
			},
			"PomProperties": {
				"artifactId": "mxparser",
				"groupId":    "io.github.x-stream",
				"location":   "META-INF\\maven\\io.github.x-stream\\mxparser\\pom.properties",
				"name":       "",
				"version":    "1.2.2",
			},
		},
	}
	javaPackages3 = model.Package{
		Name:    "bill-of-materials",
		Type:    "java",
		Version: "0.0.1-SNAPSHOT",
		Path:    "app.jar",
		Locations: []model.Location{
			{
				Path:      filepath.Join("app.jar"),
				LayerHash: "1b87eb290c214948cdbc9f40db0d9c099add021f6a5137dd4b2ea574f85eddd8",
			},
		},
		Description: "",
		Licenses: []string{
			"",
		},
		CPEs: []string{
			"cpe:2.3:a:bill-of-materials:bill-of-materials:0.0.1-SNAPSHOT:*:*:*:*:*:*:*",
			"cpe:2.3:a:bill-of-materials:bill_of-materials:0.0.1-SNAPSHOT:*:*:*:*:*:*:*",
			"cpe:2.3:a:bill-of-materials:bill_of_materials:0.0.1-SNAPSHOT:*:*:*:*:*:*:*",
			"cpe:2.3:a:bill_of-materials:bill_of_materials:0.0.1-SNAPSHOT:*:*:*:*:*:*:*",
			"cpe:2.3:a:bill_of-materials:bill-of_materials:0.0.1-SNAPSHOT:*:*:*:*:*:*:*",
			"cpe:2.3:a:bill_of-materials:bill-of-materials:0.0.1-SNAPSHOT:*:*:*:*:*:*:*",
			"cpe:2.3:a:bill_of_materials:bill-of-materials:0.0.1-SNAPSHOT:*:*:*:*:*:*:*",
			"cpe:2.3:a:bill_of_materials:bill_of-materials:0.0.1-SNAPSHOT:*:*:*:*:*:*:*",
			"cpe:2.3:a:bill_of_materials:bill_of_materials:0.0.1-SNAPSHOT:*:*:*:*:*:*:*",
			"cpe:2.3:a:carbonetes:bill-of-materials:0.0.1-SNAPSHOT:*:*:*:*:*:*:*",
			"cpe:2.3:a:carbonetes:bill_of-materials:0.0.1-SNAPSHOT:*:*:*:*:*:*:*",
			"cpe:2.3:a:carbonetes:bill_of_materials:0.0.1-SNAPSHOT:*:*:*:*:*:*:*",
		},
		PURL: model.PURL("pkg:maven/com.carbonetes/bill-of-materials@0.0.1-SNAPSHOT"),
		Metadata: JavaMetadata{
			"PomProperties": {
				"artifactId": "bill-of-materials",
				"groupId":    "com.carbonetes",
				"location":   "META-INF\\maven\\com.carbonetes\\bill-of-materials\\pom.properties",
				"name":       "",
				"version":    "0.0.1-SNAPSHOT",
			},
		},
	}
)

func TestGenerateAdditionalCPE(t *testing.T) {
	tests := []GenerateAdditionalCPEResult{
		{javaPackages1.Name, javaPackages1.Name, javaPackages1.Version, &javaPackages1, []string{"cpe:2.3:a:ezmorph:ezmorph:1.0.6:*:*:*:*:*:*:*", "cpe:2.3:a:sf:ezmorph:1.0.6:*:*:*:*:*:*:*"}},
		{javaPackages2.Name, javaPackages2.Name, javaPackages2.Version, &javaPackages2, []string{"cpe:2.3:a:mxparser:mxparser:1.2.2:*:*:*:*:*:*:*",
			"cpe:2.3:a:github:mxparser:1.2.2:*:*:*:*:*:*:*",
			"cpe:2.3:a:x-stream:mxparser:1.2.2:*:*:*:*:*:*:*",
			"cpe:2.3:a:x_stream:mxparser:1.2.2:*:*:*:*:*:*:*",
			"cpe:2.3:a:io.github.xstream.mxparser:mxparser:1.2.2:*:*:*:*:*:*:*",
			"cpe:2.3:a:io:mxparser:1.2.2:*:*:*:*:*:*:*",
			"cpe:2.3:a:xstream:mxparser:1.2.2:*:*:*:*:*:*:*"}},
	}

	for _, test := range tests {
		generateAdditionalCPE(test.vendor, test.version, test.product, test._package)
		if len(test._package.CPEs) != len(test.expected) {
			t.Errorf("Test Failed: Expected an output of %v, received: %v", len(test.expected), len(test._package.CPEs))
		}
		for i := range test._package.CPEs {
			if test._package.CPEs[i] != test.expected[i] {
				t.Errorf("Test Failed:\n Expected output of %v \n, Received: %v \n", test.expected[i], test._package.CPEs[i])
			}
		}
	}
}

func TestParseLicenses(t *testing.T) {
	tests := []JavaLicensesResult{
		{&javaPackages2, "Indiana University Extreme! Lab Software License", []string{"Indiana University Extreme! Lab Software License"}},
	}
	for _, test := range tests {
		parseLicenses(test._package)
		if len(test.expected) == 0 && len(test._package.Licenses) != 0 {
			t.Errorf("Test Failed: Slice length must be equal with the expected result. Expected: %v, Received: %v", len(test.expected), len(test._package.Licenses))
		}
		if len(test._package.Licenses) != len(test.expected) {
			t.Errorf("Test Failed: Slice length must be equal with the expected result. Expected: %v, Received: %v", len(test.expected), len(test._package.Licenses))
		}
		for i := range test._package.Licenses {
			if test._package.Licenses[i] != test.expected[i] {
				t.Errorf("Test Failed:\n Expected output of %v \n, Received: %v \n", test.expected[i], test._package.Licenses[i])
			}
		}
	}
}

func TestParseJavaURL(t *testing.T) {
	tests := []JavaPURLResult{
		{&javaPackages1, model.PURL("pkg:maven/net.sf.ezmorph/ezmorph@1.0.6")},
		{&javaPackages2, model.PURL("pkg:maven/io.github.x-stream/mxparser@1.2.2")},
	}
	for _, test := range tests {
		parseJavaURL(test._package)
		if test._package.PURL != test.expected {
			t.Errorf("Test Failed: Expected an output of %v, received: %v", test.expected, test._package.PURL)
		}
	}
}

func TestParsePomProperties(t *testing.T) {
	var data = "artifactId=bill-of-materials\r\ngroupId=com.carbonetes\r\nversion=0.0.1-SNAPSHOT\r\n"
	var path = "META-INF/maven/com.carbonetes/bill-of-materials/pom.properties"
	tests := []PomPropertiesResult{
		{data, &javaPackages3, path, javaPackages3.Metadata.(JavaMetadata)},
	}

	for _, test := range tests {
		parsePomProperties(test.data, test._package, test.path)
		if !reflect.DeepEqual(test.expected, test._package.Metadata) {
			t.Errorf("Testing Failed: Expected an output of %v, received: %v", test.expected, test._package.Metadata)
		}

	}
}

func TestFormatVersionMetadata(t *testing.T) {
	_package1 := javaPackages1
	_package2 := javaPackages2

	tests := []formatVersionMetadataResult{
		{&_package1, `1.0.0New-Metadata: metadata-content`, "1.0.0", "New-Metadata", "metadata-content"},
		{&_package2, `9.0.70Require-Capability: osgi.ee;filter:=\"(\u0026(osgi.ee=JavaSE)(version=1.8))\`, "9.0.70", "Require-Capability", `osgi.ee;filter:=\"(\u0026(osgi.ee=JavaSE)(version=1.8))\`},
	}

	for _, test := range tests {
		if output := formatVersionMetadata(test._package, test.version); output != test.expectedVersion {
			t.Errorf("Testing Failed: Expected an output of %v, received: %v", test.expectedVersion, output)
		}

		if test._package.Metadata.(JavaMetadata)["Manifest"]["Implementation-Version"] != test.expectedVersion {
			t.Errorf("Testing Failed: Expected an output of %v, received: %v", test.expectedVersion, test._package.Metadata.(JavaMetadata)["Manifest"]["Implementation-Version"])
		}

		if test._package.Metadata.(JavaMetadata)["Manifest"][test.expectedMetadataKey] != test.expectedMetadataVal {
			t.Errorf("Testing Failed: Expected an output of %v, received: %v", test.expectedMetadataVal, test._package.Metadata.(JavaMetadata)["Manifest"][test.expectedMetadataKey])
		}
	}
}
