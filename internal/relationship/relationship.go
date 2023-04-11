package relationship

import (
	"github.com/carbonetes/diggity/pkg/model"
)

var packageList map[string]*model.Package

func GetRelationships(results map[string]*model.Package) []model.Relationship {
	packageList = getPackageList(results)
	return FindOwnerships()
}
