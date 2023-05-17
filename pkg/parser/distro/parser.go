package distro

import (
	"strings"

	"github.com/carbonetes/diggity/pkg/model"
)

// Parse Linux distro
func parseLinuxDistribution(filenames []string) (*model.Distro, error) {
	distro := new(model.Distro)
	for _, filename := range filenames {
		metadata, err := parseMetadata(filename)
		if err != nil {
			return nil, err
		}

		if metadata == nil {
			continue
		}
		var ids []string
		for _, id := range strings.Split((*metadata)["ID_LIKE"], " ") {
			id = strings.TrimSpace(id)
			if id == "" {
				continue
			}
			ids = append(ids, id)
		}

		distro = &model.Distro{
			PrettyName:         (*metadata)["PRETTY_NAME"],
			Name:               (*metadata)["NAME"],
			ID:                 (*metadata)["ID"],
			IDLike:             ids,
			Version:            (*metadata)["VERSION"],
			VersionID:          (*metadata)["VERSION_ID"],
			DistribID:          (*metadata)["DISTRIB_ID"],
			DistribDescription: (*metadata)["DISTRIB_DESCRIPTIONN"],
			DistribCodename:    (*metadata)["DISTRIB_CODENAME"],
			HomeURL:            (*metadata)["HOME_URL"],
			SupportURL:         (*metadata)["SUPPORT_URL"],
			BugReportURL:       (*metadata)["BUG_REPORT_URL"],
			PrivacyPolicyURL:   (*metadata)["PRIVACY_POLICY_URL"],
		}

		if len(distro.Name) == 0 {
			continue
		}
		break
	}
	return distro, nil
}
