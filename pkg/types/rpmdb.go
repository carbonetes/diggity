package types

import (
	rpmdb "github.com/knqyf263/go-rpmdb/pkg"
)

type RpmDB struct {
	Path         string
	Layer        string
	PackageInfos []rpmdb.PackageInfo
}

func (r *RpmDB) ReadDBFile(file string) error {
	db, err := rpmdb.Open(file)
	if err != nil {
		return err
	}

	packageInfos, err := db.ListPackages()
	if err != nil {
		return err
	}

	if len(packageInfos) == 0 {
		return nil
	}

	for _, packageInfo := range packageInfos {
		if packageInfo == nil {
			continue
		}
		r.PackageInfos = append(r.PackageInfos, *packageInfo)
	}

	err = db.Close()
	if err != nil {
		return err
	}

	return nil
}
