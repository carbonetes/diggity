package alpmdb

import (
	"path/filepath"
	"testing"

	"github.com/carbonetes/diggity/pkg/model"
)

var (
	path  = filepath.Join(installedPackagesPath + "db5.3-" + "5.3.28-2")
	pkgs  = make([]model.Package, 0)
	layer = "9b7240956cfbfefddcd91a2195bfb2ed2cd17bdff81f21111849d643dfaf8131"
)

func TestAlpmdbParser(t *testing.T) {

	initPkg()
	err := readDesc(path, &pkgs, layer)
	if err != nil {
		t.Error("Test Failed: Unable to parse package", err)
	}
}

func initPkg() {
	pkgs = append(pkgs, model.Package{
		Name:    "db5.3",
		Type:    alpmdb,
		Version: "5.3.28-2",
		Path:    filepath.Join(installedPackagesPath + "db5.3-" + "5.3.28-2"),
		Locations: []model.Location{
			{
				Path:      installedPackagesPath,
				LayerHash: "9b7240956cfbfefddcd91a2195bfb2ed2cd17bdff81f21111849d643dfaf8131",
			},
		},
		Description: "The Berkeley DB embedded database system v5.3",
		Licenses: []string{
			"custom:sleepycat",
		},
	})
}
