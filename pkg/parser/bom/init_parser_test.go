package bom

import (
	"reflect"
	"testing"

	"github.com/carbonetes/diggity/pkg/model"
)

type (
	InitParsersResult struct {
		_argument *model.Arguments
		expected  *model.Arguments
	}
)

func TestInitParsers(t *testing.T) {
	argImage1 := "test-image"
	argImage2 := "test-image/test"
	argImage3 := "test-image/test:latest"
	argImage4 := "test-image/testt:xyz"
	argImage5 := "alpine"

	boolTrue := true
	boolFalse := false

	arguments1 := model.Arguments{
		Image:              &argImage1,
		DisableFileListing: &boolTrue,
	}
	arguments2 := model.Arguments{
		Image:              &argImage2,
		DisableFileListing: &boolTrue,
	}
	arguments3 := model.Arguments{
		Image:              &argImage3,
		DisableFileListing: &boolFalse,
	}
	arguments4 := model.Arguments{
		Image:              &argImage4,
		DisableFileListing: &boolTrue,
	}
	arguments5 := model.Arguments{
		Image:              &argImage5,
		DisableFileListing: &boolFalse,
	}

	tests := []InitParsersResult{
		{&arguments1, &arguments1},
		{&arguments2, &arguments2},
		{&arguments3, &arguments3},
		{&arguments4, &arguments4},
		{&arguments5, &arguments5},
	}

	for _, test := range tests {
		InitParsers(*test._argument)

		if !reflect.DeepEqual(Arguments, test.expected) {
			t.Errorf("Test Failed: Arguments must be instantiated from &arguments.")
		}
	}
}
