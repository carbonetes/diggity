package parsers

import (
	"github.com/carbonetes/diggity/pkg/scan/parsers/alpine"
	"github.com/carbonetes/diggity/pkg/scan/parsers/apt"
	"github.com/carbonetes/diggity/pkg/scan/parsers/cargo"
	"github.com/carbonetes/diggity/pkg/scan/parsers/dart"
	"github.com/carbonetes/diggity/pkg/scan/parsers/dotnet"
	"github.com/carbonetes/diggity/pkg/scan/parsers/dpkg"
	golang "github.com/carbonetes/diggity/pkg/scan/parsers/go"
	"github.com/carbonetes/diggity/pkg/scan/parsers/java"
	"github.com/carbonetes/diggity/pkg/scan/parsers/npm"
	"github.com/carbonetes/diggity/pkg/scan/parsers/php"
	"github.com/carbonetes/diggity/pkg/scan/parsers/python"
	"github.com/carbonetes/diggity/pkg/scan/parsers/r"
	"github.com/carbonetes/diggity/pkg/scan/parsers/rpm"
	"github.com/carbonetes/diggity/pkg/scan/parsers/ruby"
	"github.com/carbonetes/diggity/pkg/scan/parsers/swift"
	"github.com/carbonetes/diggity/pkg/scan/parsers/terraform"
	"github.com/carbonetes/diggity/pkg/scan/types"
)

// ModularParsers provides access to the new modular parser implementations
type ModularParsers struct{}

// NewModularParsers creates a new modular parsers factory
func NewModularParsers() *ModularParsers {
	return &ModularParsers{}
}

// NewModularPythonParser creates a new modular Python parser
func (m *ModularParsers) NewModularPythonParser() types.Scanner {
	return python.New()
}

// NewModularNPMParser creates a new modular NPM parser
func (m *ModularParsers) NewModularNPMParser() types.Scanner {
	return npm.New()
}

// NewModularJavaParser creates a new modular Java parser
func (m *ModularParsers) NewModularJavaParser() types.Scanner {
	return java.New()
}

// NewModularGoParser creates a new modular Go parser
func (m *ModularParsers) NewModularGoParser() types.Scanner {
	return golang.New()
}

// NewModularCargoParser creates a new modular Cargo parser
func (m *ModularParsers) NewModularCargoParser() types.Scanner {
	return cargo.New()
}

// NewModularPHPParser creates a new modular PHP parser
func (m *ModularParsers) NewModularPHPParser() types.Scanner {
	return php.New()
}

// NewModularRubyParser creates a new modular Ruby parser
func (m *ModularParsers) NewModularRubyParser() types.Scanner {
	return ruby.New()
}

// NewModularDotNetParser creates a new modular .NET parser
func (m *ModularParsers) NewModularDotNetParser() types.Scanner {
	return dotnet.New()
}

// NewModularSwiftParser creates a new modular Swift parser
func (m *ModularParsers) NewModularSwiftParser() types.Scanner {
	return swift.New()
}

// NewModularDartParser creates a new modular Dart parser
func (m *ModularParsers) NewModularDartParser() types.Scanner {
	return dart.New()
}

// NewModularRParser creates a new modular R parser
func (m *ModularParsers) NewModularRParser() types.Scanner {
	return r.New()
}

// NewModularTerraformParser creates a new modular Terraform parser
func (m *ModularParsers) NewModularTerraformParser() types.Scanner {
	return terraform.New()
}

// NewModularAptParser creates a new modular APT parser
func (m *ModularParsers) NewModularAptParser() types.Scanner {
	return apt.New()
}

// NewModularRpmParser creates a new modular RPM parser
func (m *ModularParsers) NewModularRpmParser() types.Scanner {
	return rpm.New()
}

// NewModularAlpineParser creates a new modular Alpine parser
func (m *ModularParsers) NewModularAlpineParser() types.Scanner {
	return alpine.New()
}

// NewModularDpkgParser creates a new modular DPKG parser
func (m *ModularParsers) NewModularDpkgParser() types.Scanner {
	return dpkg.New()
}

// GetAllModularParsers returns all available modular parsers
func (m *ModularParsers) GetAllModularParsers() []types.Scanner {
	return []types.Scanner{
		m.NewModularPythonParser(),
		m.NewModularNPMParser(),
		m.NewModularJavaParser(),
		m.NewModularGoParser(),
		m.NewModularCargoParser(),
		m.NewModularPHPParser(),
		m.NewModularRubyParser(),
		m.NewModularDotNetParser(),
		m.NewModularSwiftParser(),
		m.NewModularDartParser(),
		m.NewModularRParser(),
		m.NewModularTerraformParser(),
		m.NewModularAptParser(),
		m.NewModularRpmParser(),
		m.NewModularAlpineParser(),
		m.NewModularDpkgParser(),
	}
}
