package output

import (
	"math/rand"
	"testing"
	"time"

	"github.com/carbonetes/diggity/internal/model"
	"github.com/carbonetes/diggity/internal/parser"
)

var (
	tablePackage = model.Package{
		Name: "apk-tools",
	}
	tablePackage2 = model.Package{
		Name: "busybox",
	}
	tablePackage3 = model.Package{
		Name: "gitlab.com/yawning/obfs4.git",
	}
	tablePackage4 = model.Package{
		Name: "lzo",
	}
	tablePackage5 = model.Package{
		Name: "musl",
	}
	tablePackage6 = model.Package{
		Name: "scanelf",
	}
	tablePackage7 = model.Package{
		Name: "tzdata",
	}
	tablePackage8 = model.Package{
		Name: "xz",
	}
	tablePackage9 = model.Package{
		Name: "yum",
	}
	tablePackage0 = model.Package{
		Name: "zlib",
	}
)

func TestSortPackages(t *testing.T) {
	expected := []*model.Package{
		&tablePackage,
		&tablePackage2,
		&tablePackage3,
		&tablePackage4,
		&tablePackage5,
		&tablePackage6,
		&tablePackage7,
		&tablePackage8,
		&tablePackage9,
		&tablePackage0,
	}

	parser.Packages = expected
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(parser.Packages), func(i, j int) { parser.Packages[i], parser.Packages[j] = parser.Packages[j], parser.Packages[i] })

	sortPackages()

	for i, p := range parser.Packages {
		if p.Name != expected[i].Name {
			t.Errorf("Test Failed: Expected output of %v, received: %v", expected[i].Name, p.Name)
		}
	}
}
