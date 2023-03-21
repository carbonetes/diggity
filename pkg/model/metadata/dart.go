package metadata

type PubspecLockPackage struct {
	Packages map[string]PubspecLockMetadata `yaml:"packages"`
}
type PubspecLockMetadata struct {
	Dependency  string                 `yaml:"dependency"`
	Description PubspecLockDescription `yaml:"description"`
	Source      string                 `yaml:"source"`
	Version     string                 `yaml:"version"`
}
type PubspecLockDescription struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}
