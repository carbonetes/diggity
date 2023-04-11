package relationship

import (
	"path/filepath"
	"strings"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/parser/alpine"
	"github.com/carbonetes/diggity/pkg/parser/bom"
	"github.com/carbonetes/diggity/pkg/parser/debian"
)

const (
	ownershipType = "ownership-by-file-overlap"
)

func FindOwnerships() []model.Relationship {
	relationships := []model.Relationship{}

	for _, pkg := range bom.Packages {
		// Check if package contains metadata
		if pkg.Metadata == nil {
			continue
		}

		switch pkg.Type {
		case "apk":
			if apkRelationships(pkg) != nil {
				relationships = append(relationships, *apkRelationships(pkg))
			}
		case "deb":
			if debianRelationships(pkg) != nil {
				relationships = append(relationships, *debianRelationships(pkg))
			}
		case "python":
			if pythonRelationships(pkg) != nil {
				relationships = append(relationships, *pythonRelationships(pkg))
			}
		default:
			return nil
		}
	}

	return relationships
}

func apkRelationships(pkg *model.Package) *model.Relationship {
	metadata := (map[string]interface{})(pkg.Metadata.(alpine.Manifest))
	_, exists := metadata["Files"]
	if !exists {
		return nil
	}

	files := metadata["Files"].([]model.File)

	for _, file := range files {
		if file.Path == "" {
			continue
		}

		p := isPackage(filepath.Base(file.Path))
		if p == nil {
			continue
		}

		if p.ID == pkg.ID {
			continue
		}

		relationship := model.Relationship{
			Parent: pkg.ID,
			Child:  p.ID,
			Type:   ownershipType,
		}
		return &relationship
	}

	return nil
}

func debianRelationships(pkg *model.Package) *model.Relationship {
	metadata := (map[string]interface{})(pkg.Metadata.(debian.Metadata))
	_, exists := metadata["Conffiles"]
	if !exists {
		return nil
	}

	files := metadata["Conffiles"].([]map[string]interface{})

	for _, file := range files {
		if file["path"] == "" {
			continue
		}

		p := isPackage(filepath.Base(file["path"].(string)))
		if p == nil {
			continue
		}

		if p.ID == pkg.ID {
			continue
		}

		relationship := model.Relationship{
			Parent: pkg.ID,
			Child:  p.ID,
			Type:   ownershipType,
		}
		return &relationship
	}

	return nil
}

func pythonRelationships(pkg *model.Package) *model.Relationship {
	metadata := pkg.Metadata.(map[string]interface{})
	_, exists := metadata["files"]
	if !exists {
		return nil
	}

	files := metadata["files"].([]map[string]string)

	for _, file := range files {
		if file["path"] == "" {
			continue
		}

		p := isPackage(filepath.Base(file["path"]))
		if p == nil {
			continue
		}

		if p.ID == pkg.ID {
			continue
		}

		relationship := model.Relationship{
			Parent: pkg.ID,
			Child:  p.ID,
			Type:   ownershipType,
		}
		return &relationship
	}

	return nil
}

func isPackage(path string) *model.Package {
	_, exists := packageList[path]
	if !exists {
		return nil
	}
	return packageList[path]
}

func getPackageList(result map[string]*model.Package) map[string]*model.Package {
	pkgList := make(map[string]*model.Package)
	for key, val := range result {
		name := strings.Split(key, ":")[0]
		pkgList[name] = val
	}
	return pkgList
}
