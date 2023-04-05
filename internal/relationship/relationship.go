package relationship

import (
	"github.com/carbonetes/diggity/pkg/model"
)

func GetOwnerships() []model.Relationship {
	return FindOwnerships()
}

func FindOwnerships() []model.Relationship {
	relationships := []model.Relationship{
		{
			Parent: "1234567898765431",
			Child:  "qwertyuiopoiuytrewq",
			Type:   "contains",
		},
	}

	return relationships
}
