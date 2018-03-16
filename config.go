package lint

import (
	"bytes"
	"io"
	"os"
	"path/filepath"

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

type Rule interface {
	IsEnabled() bool
}

type rule struct {
	Enabled bool `yaml:"Enabled"`
}

func (r rule) IsEnabled() bool { return r.Enabled }

type PackageCommentRule struct {
	rule `yaml:",inline"`
}
type ImportsRule struct {
	rule `yaml:",inline"`
}
type BlankImportsRule struct {
	rule `yaml:",inline"`
}
type ExportedRule struct {
	rule `yaml:",inline"`
}
type NamesRule struct {
	rule `yaml:",inline"`
}
type VarDeclsRule struct {
	rule `yaml:",inline"`
}
type ElsesRule struct {
	rule `yaml:",inline"`
}
type IfErrorRule struct {
	rule `yaml:",inline"`
}
type RangesRule struct {
	rule `yaml:",inline"`
}
type ErrorfRule struct {
	rule `yaml:",inline"`
}
type ErrorsRule struct {
	rule `yaml:",inline"`
}
type ErrorStringsRule struct {
	rule `yaml:",inline"`
}
type ReceiverNamesRule struct {
	rule `yaml:",inline"`
}
type IncDecRule struct {
	rule `yaml:",inline"`
}
type ErrorReturnRule struct {
	rule `yaml:",inline"`
}
type UnexportedReturnRule struct {
	rule `yaml:",inline"`
}
type TimeNamesRule struct {
	rule `yaml:",inline"`
}
type ContextKeyTypesRule struct {
	rule `yaml:",inline"`
}
type ContextArgsRule struct {
	rule `yaml:",inline"`
}

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

func ReadConfigFromWorkingDir() (*Config, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	reader, err := backtrackConfig(wd)
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
	err := yaml.NewDecoder(reader).Decode(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func backtrackConfig(wd string) (io.Reader, error) {
	configPath := filepath.Join(wd, ".pikeman.yml")
	_, err := os.Stat(configPath)

	switch {
	case os.IsNotExist(err) && wd == "/":
		return bytes.NewReader(defaultConfig), nil
	case os.IsNotExist(err):
		previousDir := filepath.Dir(filepath.Join(wd, ".."))
		return backtrackConfig(previousDir)
	case err == nil:
		return os.Open(configPath)
	default:
		return nil, err
	}
}
