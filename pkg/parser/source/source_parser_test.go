package source

import (
	"path/filepath"
	"testing"
)

type (
	GetHashPathResult struct {
		input    string
		expected string
	}
)

func TestGetPathHash(t *testing.T) {
	tests := []GetHashPathResult{
		{filepath.Join("..", "..", "..", "docs", "references", "conan"), "082331eb86388735f396923918d45b5f79a3ff4ade3636844392e5138034c237"},
		{filepath.Join("..", "..", "..", "docs", "references", "go"), "7a69337018fed34d88b67c2777d2406892b88d0f01fccdf61c1a40162c26b0d3"},
		{filepath.Join("..", "..", "..", "docs", "references", "maven"), "cb1c23aa8f7f3def6dff445360b8437c40478241469a3e8b5db279217db44156"},
		{filepath.Join("..", "..", "..", "docs", "references", "npm"), "3e3040ea9bb2a3ccd8a02cfcd11b3e0990e3422df6b7f2a7c1ac8df0f932a540"},
		{filepath.Join("..", "..", "..", "docs", "references", "python"), "71129616db9643406a4410929c68b122dff2144b8bdecfa8f7c0c509ae579c54"},
	}

	for _, test := range tests {
		output, err := getPathHash(test.input)
		if err != nil {
			t.Error("Test Failed: Error occurred while parsing path hash.")
		}
		if output != test.expected {
			t.Errorf("Test Failed:\n Expected output of %v \n, Received: %v \n", test.expected, output)
		}
	}
}
