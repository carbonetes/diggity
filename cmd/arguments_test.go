package cmd_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/carbonetes/diggity/cmd"
	"github.com/carbonetes/diggity/pkg/model"
	"gotest.tools/v3/assert"
)

func TestValidateOutputArg(t *testing.T) {
	// Test with valid output types
	outputType := "json,table"
	err := cmd.ValidateOutputArg(outputType)
	if err != nil {
		t.Error(err.Error())
	}

	// Test with invalid output type
	outputType = "invalid"
	err1 := cmd.ValidateOutputArg(outputType)
	expectedError := fmt.Errorf("Invalid output type: %+v \nSupported output types: %+v", outputType, model.OutputList)
	assert.DeepEqual(t, err1.Error(), expectedError.Error())
}

func TestSplitArgs(t *testing.T) {
	// Test with single argument
	args := []string{"arg1"}
	result := cmd.SplitArgs(args)
	if len(result) != 1 || result[0] != "arg1" {
		t.Errorf("Expected single argument to be returned as is")
	}

	// Test with multiple arguments separated by comma
	args = []string{"arg1,arg2", "arg3"}
	result = cmd.SplitArgs(args)
	expectedResult := []string{"arg1", "arg2", "arg3"}
	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("Expected %v but got %v", expectedResult, result)
	}
}
