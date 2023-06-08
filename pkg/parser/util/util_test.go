package util_test

import (
	"path/filepath"
	"testing"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/parser/util"
	"gotest.tools/assert"
)

type (
	TrimUntilLayerResult struct {
		location model.Location
		expected string
	}
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

var (
	utilLocation1 = model.Location{
		Path:      filepath.Join("AppData", "Local", "Temp", "3175519915", "diggity-tmp-714678a1-5579-42fb-8e8f-0d8e2101c803", "f1a5f5ce6b163fac7f09b47645c56d2ab676bdcdb268eef06a4d9b782a75bfd0", "var", "lib", "dpkg", "status"),
		LayerHash: "f1a5f5ce6b163fac7f09b47645c56d2ab676bdcdb268eef06a4d9b782a75bfd0",
	}
	utilLocation2 = model.Location{
		Path:      filepath.Join("AppData", "Local", "Temp", "921108149", "diggity-tmp-cb5342d2-f2dd-4eb3-b6c0-0e2c9f023279", "0cd4836a36e094e1870a2e6c2578a7ad9d9cb42a7313944a6d05ab72892fc3c3", "bin", "gost"),
		LayerHash: "0cd4836a36e094e1870a2e6c2578a7ad9d9cb42a7313944a6d05ab72892fc3c3",
	}
	utilLocation3 = model.Location{
		Path:      filepath.Join("AppData", "Local", "Temp", "3175519915", "diggity-tmp-614678a1-5579-42fb-8e8f-0d8e2101c803", "69a15d957a7a6f77e3fe31f330da5f4b6b582f228917a713a7a9e59449a3f413", "var", "lib", "rpm", "Packages"),
		LayerHash: "69a15d957a7a6f77e3fe31f330da5f4b6b582f228917a713a7a9e59449a3f413",
	}
	utilLocation4 = model.Location{
		Path:      filepath.Join("AppData", "Local", "Temp", "3175519915", "diggity-tmp-614678a1-5579-42fb-8e8f-0d8e2101c803", "69a15d957a7a6f77e3fe31f330da5f4b6b582f228917a713a7a9e59449a3f413", "lib", "apk", "db", "installed"),
		LayerHash: "69a15d957a7a6f77e3fe31f330da5f4b6b582f228917a713a7a9e59449a3f413",
	}
	utilLocation5 = model.Location{
		Path: filepath.Join("AppData", "Local", "Temp", "4207199802", "diggity-tmp-c25a6d61-6bb0-4d23-90db-8aee8fe0516c", "1ea8aec45877fad7de4c11ccdf09146ce8ac4be9fe84c8ad036564f5d10b441b",
			"usr", "share", "powershell", ".store", "powershell.linux", "7.1.3", "powershell.linux", "7.1.3", "tools", "net5.0", "any", "pwsh.deps.json"),
		LayerHash: "1ea8aec45877fad7de4c11ccdf09146ce8ac4be9fe84c8ad036564f5d10b441b",
	}
)

func TestTrimUntilLayer(t *testing.T) {
	tests := []TrimUntilLayerResult{
		{utilLocation1, "var/lib/dpkg/status"},
		{utilLocation2, "bin/gost"},
		{utilLocation3, "var/lib/rpm/Packages"},
		{utilLocation4, "lib/apk/db/installed"},
		{utilLocation5, "usr/share/powershell/.store/powershell.linux/7.1.3/powershell.linux/7.1.3/tools/net5.0/any/pwsh.deps.json"},
	}

	for _, test := range tests {
		result := util.TrimUntilLayer(test.location)
		assert.Equal(t, test.expected, result)

	}
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

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := util.SplitContentsByEmptyLine(test.contents)
			assert.DeepEqual(t, test.expected, result)
		})
	}
}
