package golang

import (
	"strings"

	"github.com/carbonetes/diggity/pkg/model"
)

const (
	Type         string = "go-module"
	NoFileErrWin string = "The system cannot find the file specified."
	NoFileErrMac string = "No such file or directory exists."
)

// Split go package path
func SplitPath(path string) []string {
	return strings.Split(path, "/")
}

// Parse PURL
func SetPurl(pkg *model.Package) {
	pkg.PURL = model.PURL("pkg:golang/" + pkg.Name + "@" + pkg.Version)
}
