package output

import (
	"sort"
	"testing"

	"github.com/carbonetes/diggity/pkg/model"
)

var (
	printPackage1 = model.Package{
		ID:      "8fe93afb-86f2-4639-a3eb-6c4e787f210b",
		Name:    "lzo",
		Type:    "rpm",
		Version: "2.08",
	}
	printPackage2 = model.Package{
		ID:      "9583e9ec-df1d-484a-b560-8e1415ea92c2",
		Name:    "gitlab.com/yawning/obfs4.git",
		Type:    "go-module",
		Version: "v0.0.0-20220204003609-77af0cba934d",
	}
	printPackage3 = model.Package{
		ID:      "bdbd600f-dbdf-49a1-a329-a339f1123ffd",
		Name:    "scanelf",
		Type:    "apk",
		Version: "1.3.4-r0",
	}
	printPackage4 = model.Package{
		ID:      "418ee75b-cb1a-4abe-aad6-d757c7a91610",
		Name:    "scanf",
		Type:    "gem",
		Version: "1.0.0",
	}
	printDuplicate = model.Package{
		ID:      "418ee75b-cb1a-4abe-aad6-d757c7a91610",
		Name:    "scanf",
		Type:    "gem",
		Version: "1.0.0",
	}
	printNewVersion = model.Package{
		ID:      "519ee75c-cb1d-4abe-bad7-d758c7a91611",
		Name:    "scanf",
		Type:    "gem",
		Version: "2.0.0",
	}
)

func TestDepulicate(t *testing.T) {
	pkgs := &[]model.Package{printPackage1, printPackage2, printPackage3, printPackage4, printDuplicate}
	expected := &[]model.Package{printPackage1, printPackage2, printPackage3, printPackage4}

	_ = Depulicate(pkgs)
	// SortResults(pkgs, result)

	if len(*pkgs) != len(*expected) {
		t.Errorf("Test Failed: Expected Packages of length %+v, Received: %+v.", len(*expected), len(*pkgs))
	}

	sort.Slice(expected, func(i, j int) bool {
		return (*expected)[i].Name < (*expected)[j].Name
	})

	for i, p := range *pkgs {
		if p.Name != (*expected)[i].Name {
			t.Errorf("Test Failed: Expected output of %v, received: %v", (*expected)[i].Name, p.Name)
		}
	}
}

// func TestSortResults(t *testing.T) {
// 	pkgs := &[]model.Package{printPackage1, printPackage2, printPackage3, printPackage4, printNewVersion}
// 	expected := &[]model.Package{printPackage2, printPackage1, printPackage3, printPackage4, printNewVersion}

// 	result := Depulicate(pkgs)
// 	SortResults(pkgs, result)

// 	for i, p := range *pkgs {
// 		if p.Name != (*expected)[i].Name {
// 			t.Errorf("Test Failed: Expected output of %v, received: %v", (*expected)[i].Name, p.Name)
// 		}
// 	}
// }
