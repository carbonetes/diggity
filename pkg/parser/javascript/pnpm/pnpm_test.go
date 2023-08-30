package pnpm_test

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/parser/bom"
	"github.com/carbonetes/diggity/pkg/parser/javascript/pnpm"
)

var (
	args   = model.NewArguments()
	target = filepath.Join("..","..", "..", "..", "docs", "references", "pnpm")
)

func TestNpm(t *testing.T) {
	if _, err := os.Stat(target); os.IsNotExist(err) {
		t.Error(errors.New("Pnpm reference not found"))
		t.FailNow()
	}
	args.Dir = &target

	req, err := bom.InitParsers(args)
	if err != nil {
		t.Fatal(err)
	}
	req.WG.Add(1)
	pnpm.FindPnpmPackagesFromContent(req)
	req.WG.Wait()
	if len(*req.Errors) > 0 {
		for _, err := range *req.Errors {
			t.Error(err)
		}
	}

	if len(*req.SBOM.Packages) == 0 {
		t.Error(errors.New("No package has been found!"))
	}

	for index, p := range *req.SBOM.Packages {
		checkPackageFields(t, p, index)
	}
}

func checkPackageFields(t *testing.T, p model.Package, index int) {
	if len(p.ID) == 0 {
		t.Error(errors.New("Empty package id has been detected at index " + fmt.Sprint(index)))
	}
	if len(p.Name) == 0 {
		t.Error(errors.New("Empty package name has been detected at index " + fmt.Sprint(index)))
	}
	if len(p.Version) == 0 {
		t.Error(errors.New("Empty package version has been detected at index " + fmt.Sprint(index)))
	}
	if len(p.Type) == 0 {
		t.Error(errors.New("Empty package type has been detected at index " + fmt.Sprint(index)))
	}
	if len(p.CPEs) == 0 {
		t.Error(errors.New("Empty package cpe has been detected at index " + fmt.Sprint(index)))
	}
	if p.Metadata == nil {
		t.Error(errors.New("Nil package metadata has been detected at index " + fmt.Sprint(index)))
	}
}
