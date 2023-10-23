package distro_test

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/parser/common"
	"github.com/carbonetes/diggity/pkg/parser/distro"
)

var (
	alpine = filepath.Join("..", "..", "..", "docs", "references", "release", "alpine")
	debian = filepath.Join("..", "..", "..", "docs", "references", "release", "debian")
	rpm    = filepath.Join("..", "..", "..", "docs", "references", "release", "rpm")
	args   = model.NewArguments()
)

func TestAlpineDistro(t *testing.T) {
	args.Dir = &alpine
	req, err := common.NewParams(args)
	if err != nil {
		t.Error(errors.New("Alpine distro release reference not found"))
		t.FailNow()
	}

	req.WG.Add(1)
	distro.ParseDistro(req)
	req.WG.Wait()

	if len(*req.Errors) > 0 {
		for _, err := range *req.Errors {
			t.Error(err)
		}
	}

	if req.SBOM.Distro == nil {
		t.Error(errors.New("Alpine Distro is nil!"))
	}

	distro := req.SBOM.Distro
	if distro.ID == "" {
		t.Error(errors.New("Alpine distro id is empty!"))
	}
	if distro.Name == "" && distro.PrettyName == "" {
		t.Error(errors.New("Alpine distro name is empty!"))
	}
	if distro.VersionID == "" && distro.Version == "" {
		t.Error(errors.New("Alpine distro version id is empty!"))
	}
}

func TestDebianDistro(t *testing.T) {
	args.Dir = &debian
	req, err := common.NewParams(args)
	if err != nil {
		t.Error(errors.New("Debian distro release reference not found"))
		t.FailNow()
	}

	req.WG.Add(1)
	distro.ParseDistro(req)
	req.WG.Wait()

	if len(*req.Errors) > 0 {
		for _, err := range *req.Errors {
			t.Error(err)
		}
	}

	if req.SBOM.Distro == nil {
		t.Error(errors.New("Debian distro is nil!"))
	}

	distro := req.SBOM.Distro
	if distro.ID == "" {
		t.Error(errors.New("Debian distro id is empty!"))
	}
	if distro.Name == "" && distro.PrettyName == "" {
		t.Error(errors.New("Debian distro name is empty!"))
	}
	if distro.VersionID == "" && distro.Version == "" {
		t.Error(errors.New("Debian distro version id is empty!"))
	}
}

func TestRpmDistro(t *testing.T) {
	args.Dir = &rpm
	req, err := common.NewParams(args)
	if err != nil {
		t.Error(errors.New("Rpm distro release reference not found"))
		t.FailNow()
	}

	req.WG.Add(1)
	distro.ParseDistro(req)
	req.WG.Wait()

	if len(*req.Errors) > 0 {
		for _, err := range *req.Errors {
			t.Error(err)
		}
	}

	if req.SBOM.Distro == nil {
		t.Error(errors.New("Rpm distro is nil!"))
	}

	distro := req.SBOM.Distro
	if distro.ID == "" {
		t.Error(errors.New("Rpm distro id is empty!"))
	}
	if distro.Name == "" && distro.PrettyName == "" {
		t.Error(errors.New("Rpm distro name is empty!"))
	}
	if distro.VersionID == "" && distro.Version == "" {
		t.Error(errors.New("Rpm distro version id is empty!"))
	}
}
