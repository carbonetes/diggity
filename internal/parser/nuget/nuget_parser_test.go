package nuget

import (
	"path/filepath"
	"testing"

	"github.com/carbonetes/diggity/internal/model"
)

type (
	DotnetPurlResult struct {
		_package *model.Package
		expected model.PURL
	}
)

var (
	dotnetPackage1 = model.Package{
		Name:    "Microsoft.CodeAnalysis.Common",
		Type:    dotnet,
		Version: "3.7.0",
		Path:    "Microsoft.CodeAnalysis.Common",
		Locations: []model.Location{
			{
				Path:      filepath.Join("usr", "share", "powershell", ".store", "powershell.linux.alpine", "7.1.3", "powershell.linux.alpine", "7.1.3", "tools", "net5.0", "any", "pwsh.deps.json"),
				LayerHash: "bab4d1aab0ea326c9ff258258905fb4ffc7ddbd8bbb444d4d009e8131e01b5c0",
			},
		},
		Description: "",
		CPEs: []string{
			"cpe:2.3:a:Microsoft.CodeAnalysis.Common:Microsoft.CodeAnalysis.Common:3.7.0:*:*:*:*:*:*:*",
			"cpe:2.3:a:Microsoft:Microsoft.CodeAnalysis.Common:3.7.0:*:*:*:*:*:*:*",
			"cpe:2.3:a:CodeAnalysis:Microsoft.CodeAnalysis.Common:3.7.0:*:*:*:*:*:*:*",
			"cpe:2.3:a:Common:Microsoft.CodeAnalysis.Common:3.7.0:*:*:*:*:*:*:*",
		},
		PURL: model.PURL("pkg:dotnet/Microsoft.CodeAnalysis.Common@3.7.0"),
	}

	dotnetPackage2 = model.Package{
		Name:    "System.ServiceModel.NetTcp",
		Type:    dotnet,
		Version: "4.7.0",
		Path:    "System.ServiceModel.NetTcp",
		Locations: []model.Location{
			{
				Path:      filepath.Join("usr", "share", "powershell", ".store", "powershell.linux.alpine", "7.1.3", "powershell.linux.alpine", "7.1.3", "tools", "net5.0", "any", "pwsh.deps.json"),
				LayerHash: "bab4d1aab0ea326c9ff258258905fb4ffc7ddbd8bbb444d4d009e8131e01b5c0",
			},
		},
		Description: "Data compression library with very fast (de)compression",
		CPEs: []string{
			"cpe:2.3:a:System.ServiceModel.NetTcp:System.ServiceModel.NetTcp:4.7.0:*:*:*:*:*:*:*",
			"cpe:2.3:a:System:System.ServiceModel.NetTcp:4.7.0:*:*:*:*:*:*:*",
			"cpe:2.3:a:ServiceModel:System.ServiceModel.NetTcp:4.7.0:*:*:*:*:*:*:*",
			"cpe:2.3:a:NetTcp:System.ServiceModel.NetTcp:4.7.0:*:*:*:*:*:*:*",
		},
		PURL: model.PURL("pkg:dotnet/System.ServiceModel.NetTcp@4.7.0"),
	}

	dotnetPackage3 = model.Package{
		Name:    "MicroBuild.Core",
		Type:    dotnet,
		Version: "0.3.0",
		Path:    "MicroBuild.Core",
		Locations: []model.Location{
			{
				Path:      filepath.Join("usr", "share", "dotnet", "sdk", "5.0.202", "DotnetTools", "dotnet-dev-certs", "5.0.5-servicing.21167.8", "tools", "net5.0", "any", "dotnet-dev-certs.deps.json"),
				LayerHash: "cf50eecb4374e2055f99862cc9aaa047768296a7741765caeeb2040b57d909cb",
			},
		},
		Description: "",
		CPEs: []string{
			"cpe:2.3:a:MicroBuild.Core:MicroBuild.Core:0.3.0:*:*:*:*:*:*:*",
			"cpe:2.3:a:MicroBuild:MicroBuild.Core:0.3.0:*:*:*:*:*:*:*",
			"cpe:2.3:a:Core:MicroBuild.Core:0.3.0:*:*:*:*:*:*:*",
		},
		PURL: model.PURL("pkg:dotnet/MicroBuild.Core@0.3.0"),
	}
)

func TestParseNugetPackages(t *testing.T) {
	dotnetPath := filepath.Join("..", "..", "..", "docs", "references", "nuget", "dotnetTest.deps.json")
	testLocation := model.Location{Path: dotnetPath}
	err := parseNugetPackages(&testLocation)
	if err != nil {
		t.Error("Test Failed: Error occurred while reading Dotnet content.")
	}
}

func TestParseNugetPURL(t *testing.T) {
	tests := []DotnetPurlResult{
		{&dotnetPackage1, model.PURL("pkg:dotnet/Microsoft.CodeAnalysis.Common@3.7.0")},
		{&dotnetPackage2, model.PURL("pkg:dotnet/System.ServiceModel.NetTcp@4.7.0")},
		{&dotnetPackage3, model.PURL("pkg:dotnet/MicroBuild.Core@0.3.0")},
	}

	for _, test := range tests {
		parseNugetPURL(test._package)
		if test._package.PURL != test.expected {
			t.Errorf("Test Failed: Expected an output of %v, received: %v", test.expected, test._package.PURL)
		}
	}
}
