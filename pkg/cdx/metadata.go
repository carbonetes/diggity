package cdx

import (
	"fmt"
	"strings"
	"time"

	"github.com/CycloneDX/cyclonedx-go"
	"github.com/carbonetes/diggity/internal/helper"
	"github.com/carbonetes/diggity/internal/log"
	reader "github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/package-url/packageurl-go"
)

func setBasicMetadata() *cyclonedx.Metadata {
	return &cyclonedx.Metadata{
		Timestamp: time.Now().Format(time.RFC3339),
		Tools: &cyclonedx.ToolsChoice{
			Components: &[]cyclonedx.Component{
				{
					BOMRef:  fmt.Sprintf("pkg:%s/%s@%s", vendor, name, diggityVersion),
					Group:   vendor,
					Type:    cyclonedx.ComponentTypeApplication,
					Author:  author,
					Name:    name,
					Version: diggityVersion,
				},
			},
		},
		Authors: &[]cyclonedx.OrganizationalContact{
			{
				Name:  author,
				Email: email,
			},
		},
		Component: &cyclonedx.Component{},
	}
}

func SetImageMetadata(image v1.Image, ref reader.Reference, imageTag string) *cyclonedx.Component {
	config, err := image.ConfigFile()
	if err != nil {
		log.Errorf("error getting image config: %s", err)
	}

	layers, err := image.Layers()
	if err != nil {
		log.Errorf("error getting image layers: %s", err)
	}

	hash, err := image.Digest()
	if err != nil {
		log.Errorf("error getting image digest: %s", err)
	}

	var name, version, digest, arch, os, created string
	nameVersion := strings.SplitN(imageTag, ":", 2)
	if len(nameVersion) == 2 {
		name = nameVersion[0]
		version = nameVersion[1]
	}

	if len(name) == 0 || len(version) == 0 {
		log.Errorf("error getting image name or version")
		return nil
	}

	qualifiers := map[string]string{}
	purl := packageurl.NewPackageURL("pkg", packageurl.TypeOCI, name, version, nil, "")
	properties := &[]cyclonedx.Property{}
	arch, os, created = config.Architecture, config.OS, config.Created.String()
	digest = hash.String()
	if len(arch) != 0 {
		*properties = append(*properties, cyclonedx.Property{
			Name:  "diggity:image:architecture",
			Value: arch,
		})
		qualifiers["arch"] = arch
	}
	if len(os) != 0 {
		*properties = append(*properties, cyclonedx.Property{
			Name:  "diggity:image:os",
			Value: os,
		})
		qualifiers["os"] = os
	}

	if len(created) != 0 {
		*properties = append(*properties, cyclonedx.Property{
			Name:  "diggity:image:created",
			Value: created,
		})
	}
	if len(digest) != 0 {
		*properties = append(*properties, cyclonedx.Property{
			Name:  "diggity:image:digest",
			Value: digest,
		})

	}
	hashes := &[]cyclonedx.Hash{
		{
			Algorithm: helper.DetectCDXHashAlgorithm(digest),
			Value:     digest,
		},
	}
	for _, layer := range layers {
		hash, err := layer.Digest()
		if err != nil {
			log.Errorf("error getting layer digest: %s", err)
		}
		*properties = append(*properties, cyclonedx.Property{
			Name:  "diggity:image:layer",
			Value: hash.String(),
		})

		*hashes = append(*hashes, cyclonedx.Hash{
			Algorithm: helper.DetectCDXHashAlgorithm(hash.String()),
			Value:     hash.String(),
		})
	}

	return &cyclonedx.Component{
		BOMRef:     purl.String(),
		Type:       cyclonedx.ComponentTypeContainer,
		Name:       name,
		Version:    version,
		Properties: properties,
		Hashes:     hashes,
	}
}
