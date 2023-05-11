package alpine

import (
	"testing"

	"github.com/carbonetes/diggity/pkg/model"
)

func TestSetPURL(t *testing.T) {
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
		setPURL(test.pkg)
		if test.pkg.PURL != test.expected.PURL {
			t.Errorf("Test Failed: Expected output of %v, received: %v", test.expected.PURL, test.pkg.PURL)
		}
	}
}