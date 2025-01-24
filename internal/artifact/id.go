package artifact

import "github.com/google/uuid"

const (
	prefix = "diggity-artifact-"
)

func NewID() string {
	return prefix + NewUUID()
}

func NewUUID() string {
	return uuid.NewString()
}