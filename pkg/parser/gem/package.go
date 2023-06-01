package gem

import (
	"regexp"
	"strings"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/google/uuid"
)

func newGemLockPackage(attributes []string) *model.Package {
	var pkg model.Package
	pkg.ID = uuid.NewString()
	name, version := attributes[0], strings.Trim(attributes[1], "()")
	pkg.Name = name
	pkg.Type = Type
	pkg.Path = name
	pkg.Version = version
	//generate cpe
	generateCpes(&pkg, nil)
	metadata := make(Metadata)
	metadata["name"] = name
	metadata["version"] = version
	pkg.Metadata = metadata
	return &pkg
}

// Initialize package
func newPackage(metadata Metadata) *model.Package {
	var pkg = new(model.Package)
	var licenses = make([]string, 0)
	re := regexp.MustCompile(`[^\w^,^ ]`)

	pkg.ID = uuid.NewString()
	pkg.Type = Type
	pkg.Name = metadata["name"].(string)
	pkg.Path = metadata["name"].(string)
	pkg.Version = metadata["version"].(string)
	if val, ok := metadata["description"].(string); ok {
		pkg.Description = val
	}
	if val, ok := metadata["licenses"].(string); ok {
		tmpLicenses := re.ReplaceAllString(val, "")
		licenses = append(licenses, tmpLicenses)
	}
	pkg.Licenses = licenses
	pkg.Type = Type

	//parseURL
	setPurl(pkg)

	//check if metadata key is exist. if exist delete key to avoid duplicates
	if _, ok := metadata["metadata"].(string); ok {
		delete(metadata, "metadata")
	}

	//check if authors exists
	if val, ok := metadata["authors"].(string); ok {
		tmpAuthors := re.ReplaceAllString(val, "")
		if strings.Contains(tmpAuthors, ",") {
			arrAuthors := strings.Split(tmpAuthors, ", ")
			metadata["authors"] = arrAuthors
			for _, tmpAuthor := range arrAuthors {
				generateCpes(pkg, &tmpAuthor)
			}
		} else {
			var authors = make([]string, 0)
			authors = append(authors, tmpAuthors)
			metadata["authors"] = authors
			generateCpes(pkg, &tmpAuthors)
		}
	}

	//check if files exists
	if val, ok := metadata["files"].(string); ok {
		tmpFiles := re.ReplaceAllString(val, "")
		if strings.Contains(tmpFiles, ",") {
			metadata["files"] = strings.Split(tmpFiles, ", ")
		} else {
			var files = make([]string, 0)
			files = append(files, tmpFiles)
			metadata["files"] = files
		}
	}
	metadata["licenses"] = licenses
	pkg.Metadata = metadata

	return pkg
}

// Parse PURL
func setPurl(pkg *model.Package) {
	pkg.PURL = model.PURL("pkg" + ":" + Type + "/" + pkg.Name + "@" + pkg.Version)
}
