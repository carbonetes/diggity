package maven

import (
	"encoding/xml"
	"os"
	"strings"

	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/parser/util"
	"github.com/google/uuid"
)

//TODO: reduce the code complexities here

func parsePomXML(location model.Location, layerPath string, dir *string, result *map[string]model.Package) error {
	file, err := os.ReadFile(location.Path)
	if err != nil {
		return err
	}
	if err = xml.Unmarshal(file, &JavaPomXML); err != nil {
		return err
	}

	if len(JavaPomXML.Dependencies) > 0 {
		if len(*dir) > 0 {
			for _, dep := range JavaPomXML.Dependencies {
				if dep.ArtifactID != "" && !strings.Contains(dep.Version, "$") {
					pkg := new(model.Package)
					pkg.Metadata = Metadata{}
					pkg.ID = uuid.NewString()
					pkg.Name = dep.ArtifactID
					pkg.Path = util.TrimUntilLayer(model.Location{
						Path:      layerPath,
						LayerHash: location.LayerHash,
					})
					pkg.Version = dep.Version
					pkg.Type = Type
					paths := strings.Split(location.Path, string(os.PathSeparator))
					pkg.Locations = append(pkg.Locations, model.Location{
						Path:      paths[len(paths)-1],
						LayerHash: location.LayerHash,
					})
					pkg.Metadata.(Metadata)["ManifestLocation"] = Manifest{"path": pkg.Path}
					pkg.Metadata.(Metadata)["PomProject"] = Manifest{
						"name":    pkg.Name,
						"version": pkg.Version,
						"groupID": dep.GroupID,
					}
					parseJavaURL(pkg)
					cpe.NewCPE23(pkg, pkg.Name, pkg.Name, pkg.Version)
					generateAdditionalCPE(dep.GroupID, pkg.Name, pkg.Version, pkg)
					checkPackage(pkg, location.LayerHash, result)
				}
			}
		} else {
			if _, exists := (*result)[JavaPomXML.ArtifactID+":"+JavaPomXML.Version+":"+location.LayerHash]; exists {
				_tmpPackage := (*result)[JavaPomXML.ArtifactID+":"+JavaPomXML.Version+":"+location.LayerHash]
				_tmpPackage.Metadata.(Metadata)["PomProject"] = Manifest{
					"name":    JavaPomXML.Name,
					"version": JavaPomXML.Version,
					"groupID": JavaPomXML.GroupID,
				}
				cpe.NewCPE23(&_tmpPackage, JavaPomXML.ArtifactID, JavaPomXML.ArtifactID, JavaPomXML.Version)
				generateAdditionalCPE(JavaPomXML.GroupID, JavaPomXML.ArtifactID, JavaPomXML.Version, &_tmpPackage)
				(*result)[JavaPomXML.ArtifactID+":"+JavaPomXML.Version+":"+location.LayerHash] = _tmpPackage
			}
		}
	}
	return nil
}
