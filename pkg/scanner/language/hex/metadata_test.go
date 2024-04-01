package hex

import "testing"

func TestReadMixFile(t *testing.T) {
	content := []byte(`name1:version1:hash1:hashExt1
name2:version2:hash2:hashExt2
name3:version3:hash3:hashExt3
invalid line
name4:version4:hash4:hashExt4
`)

	expectedPackages := []HexMetadata{
		{
			Name:       "name1",
			Version:    "version1",
			PkgHash:    "hash1",
			PkgHashExt: "hashExt1",
		},
		{
			Name:       "name2",
			Version:    "version2",
			PkgHash:    "hash2",
			PkgHashExt: "hashExt2",
		},
		{
			Name:       "name3",
			Version:    "version3",
			PkgHash:    "hash3",
			PkgHashExt: "hashExt3",
		},
		{
			Name:       "name4",
			Version:    "version4",
			PkgHash:    "hash4",
			PkgHashExt: "hashExt4",
		},
	}

	packages := readMixFile(content)

	if len(packages) != len(expectedPackages) {
		t.Errorf("Test Failed: Expected %d packages, got %d", len(expectedPackages), len(packages))
	}

	for i, pkg := range packages {
		if pkg != expectedPackages[i] {
			t.Errorf("Test Failed: Expected package %v, got %v", expectedPackages[i], pkg)
		}
	}
}
