package model

// ImageInfo image information from Docker
type ImageInfo struct {
	DockerConfig   DockerConfig
	DockerManifest []DockerManifest
}

// DockerManifest "manifest.json" from docker image
type DockerManifest struct {
	Config   string      `json:"Config"`
	RepoTags interface{} `json:"RepoTags"`
	Layers   interface{} `json:"Layers"`
}

// DockerConfig "Config" object from DockerManifest
type DockerConfig struct {
	Architecture string         `json:"architecture"`
	Config       Config         `json:"config"`
	Created      string         `json:"created"`
	History      []History      `json:"history"`
	OS           string         `json:"os"`
	RootFS       RootFileSystem `json:"rootFS"`
	Variant      string         `json:"variant"`
}

// Config of DockerConfig
type Config struct {
	Env         []string    `json:"Env"`
	Entrypoint  []string    `json:"Entrypoint"`
	Cmd         []string    `json:"Cmd"`
	Workdir     string      `json:"WorkingDir"`
	ArgsEscaped bool        `json:"ArgsEscaped"`
	OnBuild     interface{} `json:"OnBuild"`
}

// History of DockerConfig
type History struct {
	Created    string `json:"created"`
	CreatedBy  string `json:"created_by"`
	EmptyLayer *bool  `json:"empty_layer,omitempty"`
	Comment    string `json:"comment,omitempty"`
}

// RootFileSystem root fs of DockerConfig
type RootFileSystem struct {
	Type         string   `json:"type"`
	DifferentIDS []string `json:"diff_ids"`
}
