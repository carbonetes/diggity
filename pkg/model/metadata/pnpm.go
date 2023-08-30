package metadata

type PnpmMetadata struct {
	LockFileVersion string                 `yaml:"lockfileVersion"`
	Dependencies    map[string]interface{} `yaml:"dependencies"`
	Packages        map[string]interface{} `yaml:"packages"`
}
