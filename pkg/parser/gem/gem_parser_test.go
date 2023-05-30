package gem

// import (
// 	"path/filepath"

// 	"github.com/carbonetes/diggity/pkg/model"

// 	"testing"
// )

// type (
// 	GemPurlResult struct {
// 		pkg      *model.Package
// 		expected model.PURL
// 	}
// 	InitGemPackageResult struct {
// 		pkg      *model.Package
// 		metadata Metadata
// 		expected *model.Package
// 	}
// )

// var (
// 	gemPackage1 = model.Package{
// 		Name:    "ostruct",
// 		Type:    "gem",
// 		Version: "0.5.2",
// 		Path:    "ostruct",
// 		Locations: []model.Location{
// 			{
// 				Path:      filepath.Join("usr", "local", "lib", "ruby", "gems", "3.1.0", "specifications", "default", "ostruct-0.5.2.gemspec"),
// 				LayerHash: "dc213d55f35dfc8beb0d34f1c56ea2fa2c2dddac029582af5a9bb5e3ca613c94",
// 			},
// 		},
// 		Description: "Class to build custom data structures, similar to a Hash.",
// 		Licenses: []string{
// 			"Ruby, BSD2Clause",
// 		},

// 		CPEs: []string{
// 			"cpe:2.3:a:ostruct:ostruct:0.5.2:*:*:*:*:*:*:*",
// 		},
// 		PURL: model.PURL("pkg:gem/ostruct@0.5.2"),
// 		Metadata: Metadata{
// 			"authors": []string{
// 				"MarcAndre Lafortune",
// 			},
// 			"bindir":      "exe",
// 			"date":        "2022-10-26",
// 			"description": "Class to build custom data structures, similar to a Hash.",
// 			"email":       "[ruby-core@marc-andre.ca]",
// 			"files": []string{
// 				"libostructrb",
// 			},
// 			"homepage": "https://github.com/ruby/ostruct",
// 			"licenses": []string{
// 				"Ruby, BSD2Clause",
// 			},
// 			"name":                      "ostruct",
// 			"require_paths":             "[lib]",
// 			"required_ruby_version":     "Gem::Requirement.new(\u003e= 2.5.0)",
// 			"required_rubygems_version": "Gem::Requirement.new(\u003e= 0) if s.respond_to? :required_rubygems_version=",
// 			"rubygems_version":          "3.3.7",
// 			"specification_version":     "4  end  if s.respond_to? :add_runtime_dependency then",
// 			"summary":                   "Class to build custom data structures, similar to a Hash.  if s.respond_to? :specification_version then",
// 			"version":                   "0.5.2",
// 		},
// 	}
// 	gemPackage2 = model.Package{
// 		Name:    "net-imap",
// 		Type:    "gem",
// 		Version: "0.2.3",
// 		Path:    "net-imap",
// 		Locations: []model.Location{
// 			{
// 				Path:      filepath.Join("usr", "local", "lib", "ruby", "gems", "3.1.0", "specifications", "net-imap-0.2.3.gemspec"),
// 				LayerHash: "dc213d55f35dfc8beb0d34f1c56ea2fa2c2dddac029582af5a9bb5e3ca613c94",
// 			},
// 		},
// 		Description: "Ruby client api for Internet Message Access Protocol",
// 		Licenses: []string{
// 			"Ruby, BSD2Clause",
// 		},

// 		CPEs: []string{
// 			"cpe:2.3:a:net-imap:net-imap:0.2.3:*:*:*:*:*:*:*",
// 		},
// 		PURL: model.PURL("pkg:gem/net-imap@0.2.3"),
// 		Metadata: Metadata{
// 			"authors": []string{
// 				"Shugo Maeda",
// 			},
// 			"bindir":               "exe",
// 			"date":                 "2022-01-06",
// 			"description":          "Ruby client api for Internet Message Access Protocol",
// 			"email":                "[shugo@ruby-lang.org]",
// 			"homepage":             "https://github.com/ruby/net-imap",
// 			"installed_by_version": "3.3.7 if s.respond_to? :installed_by_version  if s.respond_to? :specification_version then",
// 			"licenses": []string{
// 				"Ruby, BSD2Clause",
// 			},
// 			"name":                      "net-imap",
// 			"require_paths":             "[lib]",
// 			"required_ruby_version":     "Gem::Requirement.new(\u003e= 2.6.0)",
// 			"required_rubygems_version": "Gem::Requirement.new(\u003e= 0) if s.respond_to? :required_rubygems_version=",
// 			"rubygems_version":          "3.3.7",
// 			"specification_version":     "4  end  if s.respond_to? :add_runtime_dependency then",
// 			"summary":                   "Ruby client api for Internet Message Access Protocol",
// 			"version":                   "0.2.3",
// 		},
// 	}
// 	gemPackage3 = model.Package{
// 		Name:    "digest",
// 		Type:    "gem",
// 		Version: "3.1.0",
// 		Path:    "digest",
// 		Locations: []model.Location{
// 			{
// 				Path:      filepath.Join("usr", "local", "lib", "ruby", "gems", "3.1.0", "specifications", "default", "digest-3.1.0.gemspec"),
// 				LayerHash: "dc213d55f35dfc8beb0d34f1c56ea2fa2c2dddac029582af5a9bb5e3ca613c94",
// 			},
// 		},
// 		Description: "Provides a framework for message digest libraries.",
// 		Licenses: []string{
// 			"Ruby, BSD2Clause",
// 		},

// 		CPEs: []string{
// 			"cpe:2.3:a:digest:digest:3.1.0:*:*:*:*:*:*:*",
// 		},
// 		PURL: model.PURL("pkg:gem/digest@3.1.0"),
// 		Metadata: Metadata{
// 			"authors": []string{
// 				"Akinori MUSHA",
// 			},
// 			"bindir":      "exe",
// 			"date":        "2022-10-26",
// 			"description": "Provides a framework for message digest libraries.",
// 			"email":       "[knu@idaemons.org]",
// 			"extensions":  "[ext/digest/bubblebabble/extconf.rb, ext/digest/extconf.rb, ext/digest/md5/extconf.rb, ext/digest/rmd160/extconf.rb, ext/digest/sha1/extconf.rb, ext/digest/sha2/extconf.rb]",
// 			"files": []string{
// 				"digestrb",
// 				"digestso",
// 				"digestbubblebabbleso",
// 				"digestloaderrb",
// 				"digestmd5so",
// 				"digestrmd160so",
// 				"digestsha1so",
// 				"digestsha2rb",
// 				"digestsha2so",
// 				"digestsha2loaderrb",
// 				"digestversionrb",
// 				"extdigestbubblebabbleextconfrb",
// 				"extdigestextconfrb",
// 				"extdigestmd5extconfrb",
// 				"extdigestrmd160extconfrb",
// 				"extdigestsha1extconfrb",
// 				"extdigestsha2extconfrb",
// 			},
// 			"homepage": "https://github.com/ruby/digest",
// 			"licenses": []string{
// 				"Ruby, BSD2Clause",
// 			},
// 			"name":                      "digest",
// 			"require_paths":             "[lib]",
// 			"required_ruby_version":     "Gem::Requirement.new(\u003e= 2.5.0)",
// 			"required_rubygems_version": "Gem::Requirement.new(\u003e= 0) if s.respond_to? :required_rubygems_version=",
// 			"rubygems_version":          "3.3.7",
// 			"summary":                   "Provides a framework for message digest libraries.end",
// 			"version":                   "3.1.0",
// 		},
// 	}
// 	gemMetadata1 = Metadata{
// 		"authors": []string{
// 			"MarcAndre Lafortune",
// 		},
// 		"bindir":      "exe",
// 		"date":        "2022-10-26",
// 		"description": "Class to build custom data structures, similar to a Hash.",
// 		"email":       "[ruby-core@marc-andre.ca]",
// 		"files": []string{
// 			"libostructrb",
// 		},
// 		"homepage": "https://github.com/ruby/ostruct",
// 		"licenses": []string{
// 			"Ruby, BSD2Clause",
// 		},
// 		"name":                      "ostruct",
// 		"require_paths":             "[lib]",
// 		"required_ruby_version":     "Gem::Requirement.new(\u003e= 2.5.0)",
// 		"required_rubygems_version": "Gem::Requirement.new(\u003e= 0) if s.respond_to? :required_rubygems_version=",
// 		"rubygems_version":          "3.3.7",
// 		"specification_version":     "4  end  if s.respond_to? :add_runtime_dependency then",
// 		"summary":                   "Class to build custom data structures, similar to a Hash.  if s.respond_to? :specification_version then",
// 		"version":                   "0.5.2",
// 	}
// 	gemMetadata2 = Metadata{
// 		"authors": []string{
// 			"Shugo Maeda",
// 		},
// 		"bindir":               "exe",
// 		"date":                 "2022-01-06",
// 		"description":          "Ruby client api for Internet Message Access Protocol",
// 		"email":                "[shugo@ruby-lang.org]",
// 		"homepage":             "https://github.com/ruby/net-imap",
// 		"installed_by_version": "3.3.7 if s.respond_to? :installed_by_version  if s.respond_to? :specification_version then",
// 		"licenses": []string{
// 			"Ruby, BSD2Clause",
// 		},
// 		"name":                      "net-imap",
// 		"require_paths":             "[lib]",
// 		"required_ruby_version":     "Gem::Requirement.new(\u003e= 2.6.0)",
// 		"required_rubygems_version": "Gem::Requirement.new(\u003e= 0) if s.respond_to? :required_rubygems_version=",
// 		"rubygems_version":          "3.3.7",
// 		"specification_version":     "4  end  if s.respond_to? :add_runtime_dependency then",
// 		"summary":                   "Ruby client api for Internet Message Access Protocol",
// 		"version":                   "0.2.3",
// 	}
// 	gemMetadata3 = Metadata{
// 		"authors": []string{
// 			"Akinori MUSHA",
// 		},
// 		"bindir":      "exe",
// 		"date":        "2022-10-26",
// 		"description": "Provides a framework for message digest libraries.",
// 		"email":       "[knu@idaemons.org]",
// 		"extensions":  "[ext/digest/bubblebabble/extconf.rb, ext/digest/extconf.rb, ext/digest/md5/extconf.rb, ext/digest/rmd160/extconf.rb, ext/digest/sha1/extconf.rb, ext/digest/sha2/extconf.rb]",
// 		"files": []string{
// 			"digestrb",
// 			"digestso",
// 			"digestbubblebabbleso",
// 			"digestloaderrb",
// 			"digestmd5so",
// 			"digestrmd160so",
// 			"digestsha1so",
// 			"digestsha2rb",
// 			"digestsha2so",
// 			"digestsha2loaderrb",
// 			"digestversionrb",
// 			"extdigestbubblebabbleextconfrb",
// 			"extdigestextconfrb",
// 			"extdigestmd5extconfrb",
// 			"extdigestrmd160extconfrb",
// 			"extdigestsha1extconfrb",
// 			"extdigestsha2extconfrb",
// 		},
// 		"homepage": "https://github.com/ruby/digest",
// 		"licenses": []string{
// 			"Ruby, BSD2Clause",
// 		},
// 		"name":                      "digest",
// 		"require_paths":             "[lib]",
// 		"required_ruby_version":     "Gem::Requirement.new(\u003e= 2.5.0)",
// 		"required_rubygems_version": "Gem::Requirement.new(\u003e= 0) if s.respond_to? :required_rubygems_version=",
// 		"rubygems_version":          "3.3.7",
// 		"summary":                   "Provides a framework for message digest libraries.end",
// 		"version":                   "3.1.0",
// 	}
// )

// func TestReadGemContent(t *testing.T) {
// 	gemspecPath := filepath.Join("..", "..", "..", "docs", "references", "ruby", "bigdecimal-1.4.1.gemspec")
// 	testLocation := model.Location{Path: gemspecPath}
// 	pkgs := new([]model.Package)
// 	err := parseGemPackage(&testLocation, pkgs)
// 	if err != nil {
// 		t.Error("Test Failed: Error occurred while reading Gem content.")
// 	}
// }

// func TestInitGemPackages(t *testing.T) {
// 	var pkg1, pkg2, pkg3 model.Package
// 	tests := []InitGemPackageResult{
// 		{&pkg1, gemMetadata1, &gemPackage1},
// 		{&pkg2, gemMetadata2, &gemPackage2},
// 		{&pkg3, gemMetadata3, &gemPackage3},
// 	}
// 	for _, test := range tests {
// 		output := initGemPackages(test.pkg, test.metadata)
// 		if output.Name != test.expected.Name ||
// 			output.Version != test.expected.Version ||
// 			output.Description != test.expected.Description ||
// 			string(output.PURL) != string(test.expected.PURL) {

// 			t.Errorf("Test Failed:\n Expected output of %v \n, Received: %v \n", test.expected, output)
// 		}
// 		for i := range output.CPEs {
// 			if output.CPEs[i] != test.expected.CPEs[i] {
// 				t.Errorf("Test Failed:\n Expected output of %v \n, Received: %v \n", test.expected.CPEs[i], output.CPEs[i])
// 			}
// 		}
// 	}
// }

// func TestParseGemPackageURL(t *testing.T) {
// 	tests := []GemPurlResult{
// 		{&gemPackage1, model.PURL("pkg:gem/ostruct@0.5.2")},
// 		{&gemPackage2, model.PURL("pkg:gem/net-imap@0.2.3")},
// 		{&gemPackage3, model.PURL("pkg:gem/digest@3.1.0")},
// 	}
// 	for _, test := range tests {
// 		parseGemPackageURL(test.pkg)
// 		if test.pkg.PURL != test.expected {
// 			t.Errorf("Test Failed: Expected an output of %v, received: %v", test.expected, test.pkg.PURL)
// 		}
// 	}
// }
