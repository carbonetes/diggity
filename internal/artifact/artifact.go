package artifact

type Artifact struct {
	// ID is the unique identifier of the artifact
	ID string
	// Name is the name of the artifact
	Name string
	
	// Type
	Type string
	
	/*
	Entities is a map of entities that are part of the artifact.
	Package - Any relevant package information that is part of the artifact.
	Secret - Any string values that can be considered as secrets, such as passwords, tokens, etc. Check if encrypted or not.
	License - License file content or information about license.
	*/
	Entities map[string]interface{}

	// Any important tags
	Attributes map[string]string

	// Related files
	Files []File
}

type File struct {
	Name string
	Path string
	Content string
}