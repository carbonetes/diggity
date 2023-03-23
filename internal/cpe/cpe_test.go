package cpe

import (
	"reflect"
	"testing"

	"github.com/carbonetes/diggity/pkg/model"
)

type (
	NewCPE23Result struct {
		pkg      *model.Package
		vendor   string
		product  string
		version  string
		expected []string
	}

	CPEJoinResult struct {
		matchers []string
		expected string
	}

	CPEToStringResult struct {
		baseCPE  CPE
		expected string
	}

	ToCPEResult struct {
		vendor   string
		product  string
		version  string
		expected *CPE
	}

	ExpandBySeparatorsResult struct {
		baseCPE  CPE
		expected []string
	}

	ExpandResult struct {
		baseCPE   CPE
		_field    field
		separator rune
		replace   rune
		expected  []string
	}
)

func TestNewCPE23(t *testing.T) {
	var pkg1, pkg2, pkg3, pkg4, pkg5 model.Package
	tests := []NewCPE23Result{
		{&pkg1, "busybox", "busybox", "1.35.0-r17", []string{
			"cpe:2.3:a:busybox:busybox:1.35.0-r17:*:*:*:*:*:*:*",
		}},
		{&pkg2, "xmlbeans", "xmlbeans", "2.6.0", []string{
			"cpe:2.3:a:xmlbeans:xmlbeans:2.6.0:*:*:*:*:*:*:*",
		}},
		{&pkg3, "centos", "yum", "4.4.2-11.el8", []string{
			"cpe:2.3:a:centos:yum:4.4.2-11.el8:*:*:*:*:*:*:*",
			"cpe:2.3:a:yum:yum:4.4.2-11.el8:*:*:*:*:*:*:*",
		}},
		{&pkg4, "xtaci", "smux", "v1.5.16", []string{
			"cpe:2.3:a:xtaci:smux:v1.5.16:*:*:*:*:*:*:*",
			"cpe:2.3:a:smux:smux:v1.5.16:*:*:*:*:*:*:*",
		}},
		{&pkg5, "libc-utils", "libc-utils", "0.7.2-r3", []string{
			"cpe:2.3:a:libc-utils:libc-utils:0.7.2-r3:*:*:*:*:*:*:*",
			"cpe:2.3:a:libc-utils:libc_utils:0.7.2-r3:*:*:*:*:*:*:*",
			"cpe:2.3:a:libc_utils:libc_utils:0.7.2-r3:*:*:*:*:*:*:*",
			"cpe:2.3:a:libc_utils:libc-utils:0.7.2-r3:*:*:*:*:*:*:*",
			"cpe:2.3:a:libc:libc-utils:0.7.2-r3:*:*:*:*:*:*:*",
			"cpe:2.3:a:libc:libc_utils:0.7.2-r3:*:*:*:*:*:*:*",
		}},
	}

	for _, test := range tests {
		output := NewCPE23(test.pkg, test.vendor, test.product, test.version)

		if len(output.CPEs) != len(test.expected) {
			t.Errorf("Test Failed: Slice length must be equal with the expected result. Expected: %v, Received: %v", len(test.expected), len(output.CPEs))
		}

		for i, cpe := range output.CPEs {
			if cpe != test.expected[i] {
				t.Errorf("Test Failed: Expected output of %v, Received: %v", test.expected[i], cpe)
			}
		}
	}
}

func TestCpeJoin(t *testing.T) {
	tests := []CPEJoinResult{
		{[]string{"cpe:2.3", "a", "alpine_keys", "alpine_keys", "2.4-r1", "*", "*", "*", "*", "*", "*", "*"}, "cpe:2.3:a:alpine_keys:alpine_keys:2.4-r1:*:*:*:*:*:*:*"},
		{[]string{"cpe:2.3", "a", "libcrypto1", "libcrypto1.1", "1.1.1q-r0", "*", "*", "*", "*", "*", "*", "*"}, "cpe:2.3:a:libcrypto1:libcrypto1.1:1.1.1q-r0:*:*:*:*:*:*:*"},
		{[]string{"cpe:2.3", "a", "musl-utils", "musl_utils", "1.2.3-r0", "*", "*", "*", "*", "*", "*", "*"}, "cpe:2.3:a:musl-utils:musl_utils:1.2.3-r0:*:*:*:*:*:*:*"},
		{[]string{"cpe:2.3", "a", "ssl-client", "ssl_client", "1.35.0-r17", "*", "*", "*", "*", "*", "*", "*"}, "cpe:2.3:a:ssl-client:ssl_client:1.35.0-r17:*:*:*:*:*:*:*"},
		{[]string{"cpe:2.3", "a", "zlib", "zlib", "1.2.12-r3", "*", "*", "*", "*", "*", "*", "*"}, "cpe:2.3:a:zlib:zlib:1.2.12-r3:*:*:*:*:*:*:*"},
		{[]string{}, ""},
	}

	for _, test := range tests {
		if output := cpeJoin(test.matchers...); output != test.expected {
			t.Errorf("Test Failed: Expected output of %v, Received: %v", test.expected, output)
		}
	}
}

func TestCpeToString(t *testing.T) {
	tests := []CPEToStringResult{
		{CPE{
			Part:      "a",
			Vendor:    "alpine_keys",
			Product:   "alpine_keys",
			Version:   "2.4-r1",
			Update:    wildcard,
			Edition:   wildcard,
			SWEdition: wildcard,
			TargetSW:  wildcard,
			TargetHW:  wildcard,
			Other:     wildcard,
			Language:  wildcard,
		}, "cpe:2.3:a:alpine_keys:alpine_keys:2.4-r1:*:*:*:*:*:*:*"},
		{CPE{
			Part:      "a",
			Vendor:    "libcrypto1",
			Product:   "libcrypto1.1",
			Version:   "1.1.1q-r0",
			Update:    wildcard,
			Edition:   wildcard,
			SWEdition: wildcard,
			TargetSW:  wildcard,
			TargetHW:  wildcard,
			Other:     wildcard,
			Language:  wildcard,
		}, "cpe:2.3:a:libcrypto1:libcrypto1.1:1.1.1q-r0:*:*:*:*:*:*:*"},
		{CPE{
			Part:      "a",
			Vendor:    "musl-utils",
			Product:   "musl_utils",
			Version:   "1.2.3-r0",
			Update:    wildcard,
			Edition:   wildcard,
			SWEdition: wildcard,
			TargetSW:  wildcard,
			TargetHW:  wildcard,
			Other:     wildcard,
			Language:  wildcard,
		}, "cpe:2.3:a:musl-utils:musl_utils:1.2.3-r0:*:*:*:*:*:*:*"},
		{CPE{
			Part:      "a",
			Vendor:    "ssl-client",
			Product:   "ssl_client",
			Version:   "1.35.0-r17",
			Update:    wildcard,
			Edition:   wildcard,
			SWEdition: wildcard,
			TargetSW:  wildcard,
			TargetHW:  wildcard,
			Other:     wildcard,
			Language:  wildcard,
		}, "cpe:2.3:a:ssl-client:ssl_client:1.35.0-r17:*:*:*:*:*:*:*"},
		{CPE{
			Part:      "a",
			Vendor:    "zlib",
			Product:   "zlib",
			Version:   "1.2.12-r3",
			Update:    wildcard,
			Edition:   wildcard,
			SWEdition: wildcard,
			TargetSW:  wildcard,
			TargetHW:  wildcard,
			Other:     wildcard,
			Language:  wildcard,
		}, "cpe:2.3:a:zlib:zlib:1.2.12-r3:*:*:*:*:*:*:*"},
	}

	for _, test := range tests {
		if output := cpeToString(test.baseCPE); output != test.expected {
			t.Errorf("Test Failed: Expected output of %v, Received: %v", test.expected, output)
		}
	}
}

func TestToCPE(t *testing.T) {
	tests := []ToCPEResult{
		{"test", "test", "v1.0", &CPE{
			Part:      "a",
			Vendor:    "test",
			Product:   "test",
			Version:   "v1.0",
			Update:    wildcard,
			Edition:   wildcard,
			SWEdition: wildcard,
			TargetSW:  wildcard,
			TargetHW:  wildcard,
			Other:     wildcard,
			Language:  wildcard,
		}},
		{"musl", "musl", "1.2.3-r0", &CPE{
			Part:      "a",
			Vendor:    "musl",
			Product:   "musl",
			Version:   "1.2.3-r0",
			Update:    wildcard,
			Edition:   wildcard,
			SWEdition: wildcard,
			TargetSW:  wildcard,
			TargetHW:  wildcard,
			Other:     wildcard,
			Language:  wildcard,
		}},
		{"busybox", "busybox", "1.35.0-r17", &CPE{
			Part:      "a",
			Vendor:    "busybox",
			Product:   "busybox",
			Version:   "1.35.0-r17",
			Update:    wildcard,
			Edition:   wildcard,
			SWEdition: wildcard,
			TargetSW:  wildcard,
			TargetHW:  wildcard,
			Other:     wildcard,
			Language:  wildcard,
		}},
		{"zlib", "zlib", "1.2.12-r3", &CPE{
			Part:      "a",
			Vendor:    "zlib",
			Product:   "zlib",
			Version:   "1.2.12-r3",
			Update:    wildcard,
			Edition:   wildcard,
			SWEdition: wildcard,
			TargetSW:  wildcard,
			TargetHW:  wildcard,
			Other:     wildcard,
			Language:  wildcard,
		}},
		{"", "sed", "4.7-1", &CPE{
			Part:      "a",
			Vendor:    "",
			Product:   "sed",
			Version:   "4.7-1",
			Update:    wildcard,
			Edition:   wildcard,
			SWEdition: wildcard,
			TargetSW:  wildcard,
			TargetHW:  wildcard,
			Other:     wildcard,
			Language:  wildcard,
		}},
	}

	for _, test := range tests {
		if output := toCPE(test.vendor, test.product, test.version); !reflect.DeepEqual(output, test.expected) {
			t.Errorf("Test Failed: Expected output of %v, Received: %v", test.expected, output)
		}
	}
}

func TestExpandCPEsBySeparators(t *testing.T) {
	tests := []ExpandBySeparatorsResult{
		{CPE{
			Part:      "a",
			Vendor:    "libcrypto1.1",
			Product:   "libcrypto1.1",
			Version:   "1.1.1q-r0",
			Update:    "*",
			Edition:   "*",
			SWEdition: "*",
			TargetSW:  "*",
			TargetHW:  "*",
			Other:     "*",
			Language:  "*",
		}, []string{
			"cpe:2.3:a:libcrypto1:libcrypto1.1:1.1.1q-r0:*:*:*:*:*:*:*",
			"cpe:2.3:a:1:libcrypto1.1:1.1.1q-r0:*:*:*:*:*:*:*",
			"cpe:2.3:a:libcrypto1.1:libcrypto1.1:1.1.1q-r0:*:*:*:*:*:*:*",
			"cpe:2.3:a:libcrypto1.1:libcrypto1.1:1.1.1q-r0:*:*:*:*:*:*:*",
		},
		},
		{CPE{
			Part:      "a",
			Vendor:    "ssl_client",
			Product:   "ssl_client",
			Version:   "1.35.0-r17",
			Update:    "*",
			Edition:   "*",
			SWEdition: "*",
			TargetSW:  "*",
			TargetHW:  "*",
			Other:     "*",
			Language:  "*",
		}, []string{
			"cpe:2.3:a:ssl_client:ssl_client:1.35.0-r17:*:*:*:*:*:*:*",
			"cpe:2.3:a:ssl_client:ssl_client:1.35.0-r17:*:*:*:*:*:*:*",
			"cpe:2.3:a:ssl_client:ssl-client:1.35.0-r17:*:*:*:*:*:*:*",
			"cpe:2.3:a:ssl-client:ssl-client:1.35.0-r17:*:*:*:*:*:*:*",
			"cpe:2.3:a:ssl-client:ssl-client:1.35.0-r17:*:*:*:*:*:*:*",
			"cpe:2.3:a:ssl-client:ssl_client:1.35.0-r17:*:*:*:*:*:*:*",
			"cpe:2.3:a:ssl_client:ssl_client:1.35.0-r17:*:*:*:*:*:*:*",
			"cpe:2.3:a:ssl_client:ssl_client:1.35.0-r17:*:*:*:*:*:*:*",
			"cpe:2.3:a:ssl_client:ssl-client:1.35.0-r17:*:*:*:*:*:*:*",
			"cpe:2.3:a:ssl-client:ssl-client:1.35.0-r17:*:*:*:*:*:*:*",
			"cpe:2.3:a:ssl-client:ssl-client:1.35.0-r17:*:*:*:*:*:*:*",
			"cpe:2.3:a:ssl-client:ssl_client:1.35.0-r17:*:*:*:*:*:*:*",
			"cpe:2.3:a:ssl:ssl_client:1.35.0-r17:*:*:*:*:*:*:*",
			"cpe:2.3:a:ssl:ssl-client:1.35.0-r17:*:*:*:*:*:*:*",
		},
		},
		{CPE{
			Part:      "a",
			Vendor:    "zlib",
			Product:   "zlib",
			Version:   "1.2.12-r3",
			Update:    "*",
			Edition:   "*",
			SWEdition: "*",
			TargetSW:  "*",
			TargetHW:  "*",
			Other:     "*",
			Language:  "*",
		}, []string{
			"cpe:2.3:a:zlib:zlib:1.2.12-r3:*:*:*:*:*:*:*",
			"cpe:2.3:a:zlib:zlib:1.2.12-r3:*:*:*:*:*:*:*",
		},
		},
	}

	for _, test := range tests {
		output := expandCPEsBySeparators(test.baseCPE)

		if len(test.expected) == 0 {
			if len(output) != 0 {
				t.Errorf("Test Failed: Slice length must be equal with the expected result. Expected: %v, Received: %v", len(test.expected), len(output))
			}
			continue
		}
		if len(output) != len(test.expected) {
			t.Errorf("Test Failed: Slice length must be equal with the expected result. Expected: %v, Received: %v", len(test.expected), len(output))
		}

		for i, field := range output {
			if field != test.expected[i] {
				t.Errorf("Test Failed: Expected output of %v, Received: %v", test.expected[i], field)
			}
		}
	}
}

func TestExpand(t *testing.T) {
	tests := []ExpandResult{
		{CPE{
			Vendor:  "alpine_keys",
			Product: "alpine_keys",
		}, "Vendor", 45, 95, []string{
			"alpine_keys",
			"alpine-keys",
		},
		},
		{CPE{
			Vendor:  "ANY",
			Product: "lsb-base",
		}, "Product", 45, 95, []string{
			"lsb-base",
			"lsb_base",
		},
		},
		{CPE{
			Vendor:  "centos",
			Product: "kexec-tools",
		}, "Product", 45, 95, []string{
			"kexec-tools",
			"kexec_tools",
		},
		},
		{CPE{}, "Vendor", 45, 95, []string{""}},
		{CPE{}, "Product", 45, 95, []string{""}},
		{CPE{}, "", 45, 95, nil},
	}

	for _, test := range tests {
		output := expand(test.baseCPE, test._field, test.separator, test.replace)

		if len(test.expected) == 0 {
			if len(output) != 0 {
				t.Errorf("Test Failed: Slice length must be equal with the expected result. Expected: %v, Received: %v", len(test.expected), len(output))
			}
			continue
		}
		if len(output) != len(test.expected) {
			t.Errorf("Test Failed: Slice length must be equal with the expected result. Expected: %v, Received: %v", len(test.expected), len(output))
		}

		for i, field := range output {
			if field != test.expected[i] {
				t.Errorf("Test Failed: Expected output of %v, Received: %v", test.expected[i], field)
			}
		}
	}
}
