package util_test

import (
	"testing"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/parser/util"
	"gotest.tools/assert"
)

type (
	IndexOfResult struct {
		array    []string
		s        string
		expected int
	}

	StringSliceContainsResult struct {
		s        []string
		e        string
		expected bool
	}

	FormatLockKeyValResult struct {
		input    string
		expected string
	}

	SplitContentsByEmptyLineResult struct {
		name     string
		contents string
		expected []string
	}
)

func TestTrimUntilLayer(t *testing.T) {
	// Test case 1: location with layer hash in the middle of path
	location := model.Location{
		Path:      "\\diggity-tmp-faba00cf-a55c-4635-b0b6-1f648498e790\\187d3fbd6d78f45ee5d316f07a4e3721c8e4e91b75cbf9ee0ab0ac1bcbef78e8\\lib\\apk\\db\\installed",
		LayerHash: "187d3fbd6d78f45ee5d316f07a4e3721c8e4e91b75cbf9ee0ab0ac1bcbef78e8",
	}
	expected := "lib/apk/db/installed"
	result := util.TrimUntilLayer(location)
	assert.Equal(t, result, expected)
}

func TestIndexOf(t *testing.T) {
	var array1 = []string{"test1", "test2", "test3", "test4", "test5"}
	var array2 = []string{"a", "B", "c", "D", "e"}
	var array3 = []string{"", "", "?", "C:", "Users", "Username", "AppData", "Local", "Temp", "3260872682", "diggity-tmp-64a6619c-a0fe-4208-822f-67300fa7bf89", "9b7240956cfbfefddcd91a2195bfb2ed2cd17bdff81f21111849d643dfaf8131", "bin", "busybox"}
	var arrayempty = []string{}

	tests := []IndexOfResult{
		{array1, "test1", 0},
		{array1, "test3", 2},
		{array1, "test5", 4},
		{array1, "x", -1},
		{array2, "e", 4},
		{array2, "C", -1},
		{array3, "9b7240956cfbfefddcd91a2195bfb2ed2cd17bdff81f21111849d643dfaf8131", 11},
		{array3, "test", -1},
		{arrayempty, "test", -1},
		{arrayempty, "", -1},
	}

	for _, test := range tests {
		result := util.IndexOf(test.array, test.s)
		assert.Equal(t, result, test.expected)
	}

}

func TestStringSliceContains(t *testing.T) {
	tests := []StringSliceContainsResult{
		{[]string{"test1", "test2", "test3"}, "test1", true},
		{[]string{"test1", "test2", "test3"}, "testX", false},
		{[]string{"java"}, "java", true},
		{[]string{"java"}, "alpine", false},
		{[]string{"java", "npm", "deb"}, "npm", true},
		{[]string{"java", "npm", "deb"}, "alpine", false},
		{[]string{""}, "java", false},
	}

	for _, test := range tests {
		result := util.StringSliceContains(test.s, test.e)
		assert.Equal(t, test.expected, result)
	}
}

func TestFormatLockKeyVal(t *testing.T) {
	tests := []FormatLockKeyValResult{
		{`"test"`, "test"},
		{` "test" `, "test"},
		{`"name"`, "name"},
		{`"version"`, "version"},
		{`"checksum"`, "checksum"},
		{` "zerofrom" `, "zerofrom"},
		{`"zerovec-derive"`, "zerovec-derive"},
	}

	for _, test := range tests {
		result := util.FormatLockKeyVal(test.input)
		assert.Equal(t, test.expected, result)
	}
}

func TestSplitContentsByEmptyLine(t *testing.T) {
	tests := []SplitContentsByEmptyLineResult{
		{
			name:     "empty string",
			contents: "",
			expected: []string{""},
		},
		{
			name:     "single line",
			contents: "hello world",
			expected: []string{"hello world"},
		},
		{
			name:     "multiple lines with empty line in between",
			contents: "hello\n\nworld",
			expected: []string{"hello", "world"},
		},
		{
			name:     "multiple lines without empty line in between",
			contents: "hello\nworld",
			expected: []string{"hello\nworld"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := util.SplitContentsByEmptyLine(tt.contents)
			assert.DeepEqual(t, tt.expected, result)
		})
	}
}
