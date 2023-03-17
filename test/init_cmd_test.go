package test

import (
	"testing"

	cmd "github.com/carbonetes/diggity/cmd"
)

type (
	SplitArgsReult struct {
		input    []string
		expected []string
	}
)

func TestValidateOutputArg(t *testing.T) {
	tests := []string{"json", "table", "cyclonedx-xml", "cyclonedx-json", "spdx-json", "spdx-tag-value", "github-json"}

	for _, test := range tests {
		cmd.ValidateOutputArg(test)
	}
}

func TestSplitArgs(t *testing.T) {
	tests := []SplitArgsReult{
		{[]string{"apk"}, []string{"apk"}},
		{[]string{"apk", "go"}, []string{"apk", "go"}},
		{[]string{"apk,go", "deb"}, []string{"apk", "go", "deb"}},
		{[]string{"apk,go", "deb,java"}, []string{"apk", "go", "deb", "java"}},
		{[]string{"apk,go", "deb,java,rpm"}, []string{"apk", "go", "deb", "java", "rpm"}},
	}

	for _, test := range tests {
		output := cmd.SplitArgs(test.input)
		if len(output) != len(test.expected) {
			t.Errorf("Test Failed: Slice length must be equal with the expected result. Expected: %v, Received: %v", len(test.expected), len(output))
		}

		for i, arg := range output {
			if arg != test.expected[i] {
				t.Errorf("Test Failed: Expected output of %v, Received: %v", test.expected[i], output)
			}
		}
	}
}
