package stream

import (
	"time"

	"github.com/CycloneDX/cyclonedx-go"
	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/pkg/types"
	v1 "github.com/google/go-containerregistry/pkg/v1"
)

func SetDefaultValues() {
	store.Set(FileListStoreKey, []string{})
	store.Set(ComponentsStoreKey, []types.Component{})
	store.Set(SecretsStoreKey, []types.Secret{})
	store.Set(SBOMStoreKey, types.NewSBOM())
	store.Set(ParametersStoreKey, types.Parameters{})
	store.Set(ConfigStoreKey, types.Config{})
	store.Set(CycloneDXComponentsStoreKey, []cyclonedx.Component{})
	store.Set(OSReleasesStoreKey ,[]types.OSRelease{})
}

func AddComponent(component types.Component) {
	data, exist := store.Get(ComponentsStoreKey)

	components, ok := data.([]types.Component)

	if !ok {
		log.Error("Received invalid component slice from store")
	}

	if !exist {
		store.Set(ComponentsStoreKey, []types.Component{component})
	}
	for _, c := range components {
		if c.Name == component.Name && c.Version == component.Version && c.Type == component.Type {
			return
		}
	}
	components = append(components, component)
	store.Set(ComponentsStoreKey, components)

}

func AddSecret(secret types.Secret) {
	data, exist := store.Get(SecretsStoreKey)

	secrets, ok := data.([]types.Secret)

	if !ok {
		log.Error("Received invalid secret slice from store")
	}

	if !exist {
		store.Set(SecretsStoreKey, []types.Secret{secret})
	}

	secrets = append(secrets, secret)
	store.Set(SecretsStoreKey, secrets)
}

func AddCdxComponent(component cyclonedx.Component) {
	data, exist := store.Get(CycloneDXComponentsStoreKey)

	components, ok := data.([]cyclonedx.Component)

	if !ok {
		log.Error("Received invalid component slice from store")
	}

	if !exist {
		store.Set(CycloneDXComponentsStoreKey, []cyclonedx.Component{component})
	}

	components = append(components, component)
	store.Set(CycloneDXComponentsStoreKey, components)
}

func AddFile(file string) {
	data, exist := store.Get(FileListStoreKey)

	files, ok := data.([]string)

	if !ok {
		log.Error("Received invalid file slice from store")
	}

	if !exist {
		store.Set(FileListStoreKey, []string{file})
	}

	files = append(files, file)
	store.Set(FileListStoreKey, files)
}

func SetParameters(params types.Parameters) {
	store.Set(ParameterScanTypeStoreKey, params.ScanType)
	store.Set(ParameterInputStoreKey, params.Input)
	store.Set(ParameterOutputFormatStoreKey, params.OutputFormat)
	store.Set(ParameterQuietStoreKey, params.Quiet)
	store.Set(ParameterScannersStoreKey, params.Scanners)
	store.Set(ParameterAllowFileListingStoreKey, params.AllowFileListing)
	store.Set(ScanStartStoreKey, time.Now())
	store.Set(ParametersStoreKey, params)
}

func SetImageInstance(image v1.Image) {
	store.Set(ImageInstanceStoreKey, image)
}

func SetScanElapsed(duration float64) {
	store.Set(ScanElapsedStoreKey, duration)
}

func SetConfig(config types.Config) {
	store.Set(ConfigStoreKey, config)
}

func AddOSRelease(release types.OSRelease) {
	data, exist := store.Get(OSReleasesStoreKey)

	releases, ok := data.([]types.OSRelease)

	if !ok {
		log.Error("Received invalid OS release slice from store")
	}

	if !exist {
		store.Set(OSReleasesStoreKey, []types.OSRelease{release})
	}

	releases = append(releases, release)
	store.Set(OSReleasesStoreKey, releases)
}