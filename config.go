package lint

import (
	"filepath"
	"fmt"
	"io"
	"ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

var defaultConfig = []byte(`
Includes:
  - ./**/*.go

Excludes:
  - ./vendor/**
  - ./pkg/**

PackageComment:
  Enabled: true

Imports:
  Enabled: true

BlankImports:
  Enabled: true

Exported:
  Enabled: true

Names:
  Enabled: true

VarDecls:
  Enabled: true

Elses:
  Enabled: true

IfError:
  Enabled: true

Ranges:
  Enabled: true

Errorf:
  Enabled: true

Errors:
  Enabled: true

ErrorStrings:
  Enabled: true

ReceiverNames:
  Enabled: true

IncDec:
  Enabled: true

ErrorReturn:
  Enabled: true

UnexportedReturn:
  Enabled: true

TimeNames:
  Enabled: true

ContextKeyTypes:
  Enabled: true

ContextArgs:
  Enabled: true
`)

type Rule struct {
	Enabled bool `yaml:"Enabled"`
}

type PackageCommentRule struct{ Rule }
type ImportsRule struct{ Rule }
type BlankImportsRule struct{ Rule }
type ExportedRule struct{ Rule }
type NamesRule struct{ Rule }
type VarDeclsRule struct{ Rule }
type ElsesRule struct{ Rule }
type IfErrorRule struct{ Rule }
type RangesRule struct{ Rule }
type ErrorfRule struct{ Rule }
type ErrorsRule struct{ Rule }
type ErrorStringsRule struct{ Rule }
type ReceiverNamesRule struct{ Rule }
type IncDecRule struct{ Rule }
type ErrorReturnRule struct{ Rule }
type UnexportedReturnRule struct{ Rule }
type TimeNamesRule struct{ Rule }
type ContextKeyTypesRule struct{ Rule }
type ContextArgsRule struct{ Rule }

type Config struct {
	Includes []string `yaml:"Includes"`
	Excludes []string `yaml:"Excludes"`

	PackageComment   PackageCommentRule   `yaml:"PackageComment"`
	Imports          ImportsRule          `yaml:"Imports"`
	BlankImports     BlankImportsRule     `yaml:"BlankImports"`
	Exported         ExportedRule         `yaml:"Exported"`
	Names            NamesRule            `yaml:"Names"`
	VarDecls         VarDeclsRule         `yaml:"VarDecls"`
	Elses            ElsesRule            `yaml:"Elses"`
	IfError          IfErrorRule          `yaml:"IfError"`
	Ranges           RangesRule           `yaml:"Ranges"`
	Errorf           ErrorfRule           `yaml:"Errorf"`
	Errors           ErrorsRule           `yaml:"Errors"`
	ErrorStrings     ErrorStringsRule     `yaml:"ErrorStrings"`
	ReceiverNames    ReceiverNamesRule    `yaml:"ReceiverNames"`
	IncDec           IncDecRule           `yaml:"IncDec"`
	ErrorReturn      ErrorReturnRule      `yaml:"ErrorReturn"`
	UnexportedReturn UnexportedReturnRule `yaml:"UnexportedReturn"`
	TimeNames        TimeNamesRule        `yaml:"TimeNames"`
	ContextKeyTypes  ContextKeyTypesRule  `yaml:"ContextKeyTypes"`
	ContextArgs      ContextArgsRule      `yaml:"ContextArgs"`
}

func ReadConfigFromWorkspace() (*Config, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	config, err := searchConfig(wd)
	if err != nil {
		return nil, err
	}

	return decodeConfig(reader)
}

func ReadConfig(path string) (*Config, error) {
	reader, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return decodeConfig(reader)
}

func decodeConfig(reader io.Reader) (*Config, error) {
	var config Config
	err = yaml.NewDecoder(reader).Decode(&config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func backtrackConfig(wd string) (io.Reader, err) {
	configPath := filepath.Join(wd, ".pikeman.yml")
	_, err := os.Stat(configPath)

	switch {
	case os.IsNotExist(err) && wd == "/":
		return nil, fmt.Errorf(".pikeman.yml is nowhere to be found.")
	case os.IsNotExist(err):
		previousDir := filepath.Dir(filepath.Join(wd, ".."))
		return backtrackConfig(previousDir)
	case err == nil:
		return os.Open(configPath)
	default:
		return nil, err
	}
}
