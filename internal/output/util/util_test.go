package util

import (
	"math/rand"
	"testing"
	"time"

	"github.com/carbonetes/diggity/internal/model"
	"github.com/carbonetes/diggity/internal/parser/bom"
)

var (
	sortPackage = model.Package{
		Name: "apk-tools",
	}
	sortPackage2 = model.Package{
		Name: "busybox",
	}
	sortPackage3 = model.Package{
		Name: "gitlab.com/yawning/obfs4.git",
	}
	sortPackage4 = model.Package{
		Name: "lzo",
	}
	sortPackage5 = model.Package{
		Name: "musl",
	}
	sortPackage6 = model.Package{
		Name: "scanelf",
	}
	sortPackage7 = model.Package{
		Name: "tzdata",
	}
	sortPackage8 = model.Package{
		Name: "xz",
	}
	sortPackage9 = model.Package{
		Name: "yum",
	}
	sortPackage0 = model.Package{
		Name: "zlib",
	}
)

func TestSortPackages(t *testing.T) {
	expected := []*model.Package{
		&sortPackage,
		&sortPackage2,
		&sortPackage3,
		&sortPackage4,
		&sortPackage5,
		&sortPackage6,
		&sortPackage7,
		&sortPackage8,
		&sortPackage9,
		&sortPackage0,
	}

	bom.Packages = expected
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(bom.Packages), func(i, j int) { bom.Packages[i], bom.Packages[j] = bom.Packages[j], bom.Packages[i] })

	SortPackages()

	for i, p := range bom.Packages {
		if p.Name != expected[i].Name {
			t.Errorf("Test Failed: Expected output of %v, received: %v", expected[i].Name, p.Name)
		}
	}
}
