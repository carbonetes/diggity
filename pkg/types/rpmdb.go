package types

import (
	rpmdb "github.com/knqyf263/go-rpmdb/pkg"
)

type RpmDB struct {
	Path         string
	Layer        string
	PackageInfos []rpmdb.PackageInfo
}

func (r *RpmDB) GetRpmPackageInfos() error {
	db, err := rpmdb.Open(r.Path)
	if err != nil {
		return err
	}

	defer db.Close()

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

	return nil
}
