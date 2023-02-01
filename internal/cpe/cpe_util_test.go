package cpe

import (
	"testing"
)

type (
	RemoveDuplicateCPESResult struct {
		cpes     []string
		expected []string
	}

	ValidateCPEResult struct {
		cpe     string
		isValid bool
	}
)

func TestRemoveDuplicateCPES(t *testing.T) {
	tests := []RemoveDuplicateCPESResult{
		{
			[]string{"cpe:2.3:a:musl-utils:musl-utils:1.2.3-r0:*:*:*:*:*:*:*",
				"cpe:2.3:a:musl-utils:musl_utils:1.2.3-r0:*:*:*:*:*:*:*",
				"cpe:2.3:a:musl_utils:musl_utils:1.2.3-r0:*:*:*:*:*:*:*",
				"cpe:2.3:a:musl-utils:musl_utils:1.2.3-r0:*:*:*:*:*:*:*",
				"cpe:2.3:a:musl_utils:musl_utils:1.2.3-r0:*:*:*:*:*:*:*",
				"cpe:2.3:a:musl_utils:musl-utils:1.2.3-r0:*:*:*:*:*:*:*"},
			[]string{"cpe:2.3:a:musl-utils:musl-utils:1.2.3-r0:*:*:*:*:*:*:*",
				"cpe:2.3:a:musl-utils:musl_utils:1.2.3-r0:*:*:*:*:*:*:*",
				"cpe:2.3:a:musl_utils:musl_utils:1.2.3-r0:*:*:*:*:*:*:*",
				"cpe:2.3:a:musl_utils:musl-utils:1.2.3-r0:*:*:*:*:*:*:*"},
		},
		{
			[]string{
				"cpe:2.3:a:symfony:framework-bundle:v5.4.10:*:*:*:*:*:*:*",
				"cpe:2.3:a:symfony:framework-bundle:v5.4.10:*:*:*:*:*:*:*",
				"cpe:2.3:a:symfony:framework_bundle:v5.4.10:*:*:*:*:*:*:*",
				"cpe:2.3:a:framework-bundle:framework-bundle:v5.4.10:*:*:*:*:*:*:*",
			},
			[]string{
				"cpe:2.3:a:symfony:framework-bundle:v5.4.10:*:*:*:*:*:*:*",
				"cpe:2.3:a:symfony:framework_bundle:v5.4.10:*:*:*:*:*:*:*",
				"cpe:2.3:a:framework-bundle:framework-bundle:v5.4.10:*:*:*:*:*:*:*",
			},
		},
		{
			[]string{
				"cpe:2.3:a:tomarrell:wrapcheck:v1.0.0:*:*:*:*:*:*:*",
				"cpe:2.3:a:wrapcheck:wrapcheck:v1.0.0:*:*:*:*:*:*:*",
			},
			[]string{
				"cpe:2.3:a:tomarrell:wrapcheck:v1.0.0:*:*:*:*:*:*:*",
				"cpe:2.3:a:wrapcheck:wrapcheck:v1.0.0:*:*:*:*:*:*:*",
			},
		},
		{
			[]string{
				"cpe:2.3:a:e2fsprogs:e2fsprogs:1.46.2-2:*:*:*:*:*:*:*",
			},
			[]string{
				"cpe:2.3:a:e2fsprogs:e2fsprogs:1.46.2-2:*:*:*:*:*:*:*",
			},
		},
		{[]string{}, []string{}},
	}

	for _, test := range tests {
		output := RemoveDuplicateCPES(test.cpes)

		if len(output) != len(test.expected) {
			t.Errorf("Test Failed: Slice length must be equal with the expected result. Expected: %v, Received: %v", len(test.expected), len(output))
		}

		for i, cpe := range output {
			if cpe != test.expected[i] {
				t.Errorf("Test Failed: Expected output of %v, Received: %v", test.expected[i], cpe)
			}
		}
	}
}

func TestValidateCPE(t *testing.T) {
	tests := []ValidateCPEResult{
		{"cpe:2.3:a:zlib:zlib:1.2.12-r3:*:*:*:*:*:*:*", true},
		{"cpe:2.3:a:e2fsprogs:e2fsprogs:1.46.2-2:*:*:*:*:*:*:*", true},
		{"cpe:2.3:a:lzo:lzo:2.08-14.el8:*:*:*:*:*:*:*", true},
		{"cpe:2.3:a:test:test:v1.0:*:*:*:*:*:*:*", true},
		{"cpe:2.3:a:test:test:v1.0:*:*:*:*:*:*:?", false},
		{"cpe:2.3:a:test:test:v1.0:*:*:/:*:*:*:*", false},
		{"cpe:2.3:a:test:test:v1.0:*:*:*:*:*:*", false},
		{"abc:2.3:a:test:test:v1.0:*:*:*:*:*:*:*", false},
		{"cpe:test", false},
		{"test", false},
		{"", false},
	}

	for _, test := range tests {
		err := validateCPE(test.cpe)
		if test.isValid && err != nil {
			t.Error("Test Failed: Error occurred for valid CPE string.")
		}
		if !test.isValid && err == nil {
			t.Error("Test Failed: Expected error for invalid CPE string.")
		}
	}
}
