package alpm

import "strings"

func parseMetadata(attributes []string) Metadata {
	var metadata = make(Metadata)
	for _, attribute := range attributes {
		if attribute == "" {
			continue
		}
		attribute = strings.TrimSpace(attribute)
		properties := strings.Split(attribute, "\n")
		key := properties[0]
		values := properties[1:]

		// Attributes based on https://gitlab.archlinux.org/pacman/pacman/-/blob/master/lib/libalpm/be_local.c
		switch key {
		case "%NAME%":
			metadata["Name"] = values[0]
		case "%VERSION%":
			metadata["Version"] = values[0]
		case "%BASE%":
			metadata["Base"] = values[0]
		case "%DESC%":
			metadata["Description"] = values[0]
		case "%GROUP%":
			metadata["Group"] = values
		case "%URL%":
			metadata["URL"] = values[0]
		case "%ARCH":
			metadata["Architecture"] = values[0]
		case "%BUILDDATE%":
			metadata["BuildDate"] = values[0]
		case "%INSTALLDATE%":
			metadata["InstallDate"] = values[0]
		case "%PACKAGER%":
			metadata["Packager"] = values[0]
		case "%SIZE%":
			metadata["Size"] = values[0]
		case "%REASON%":
			metadata["Reason"] = values[0]
		case "%LICENSE%":
			metadata["Licenses"] = values
		case "%VALIDATION%":
			metadata["Validation"] = values[0]
		case "%REPLACES%":
			metadata["Replaces"] = values
		case "%DEPENDS%":
			metadata["Depends"] = values
		case "%OPTDEPENDS%":
			metadata["OptDepends"] = values
		case "%MAKEDEPENDS%":
			metadata["MakeDepends"] = values
		case "%CHECKDEPENDS%":
			metadata["CheckDepends"] = values
		case "%CONFLICTS%":
			metadata["Conflicts"] = values
		case "%PROVIDES%":
			metadata["Provides"] = values
		case "%XDATA%":
			metadata["XData"] = values
		}
	}
	return metadata
}
