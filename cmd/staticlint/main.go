// custom multichecker
// is used to combine different analyzers/linters and use them with 1 command
// for usage, see root README.md section "linting & static checks"
// uses config.json to specify https://staticcheck.dev/ rules
// also uses 2 custom analyzers errcheck.New and noosexit.New
package main

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/gennadyterekhov/metrics-storage/cmd/staticlint/analyzers/errcheck"
	"github.com/gennadyterekhov/metrics-storage/cmd/staticlint/analyzers/noosexit"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/quickfix"
	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"
)

const ConfigFileName = `config.json`

type Config struct {
	QuickfixRules    map[string]bool `json:"honnef.co/go/tools/quickfix"`
	StylecheckRules  map[string]bool `json:"honnef.co/go/tools/stylecheck"`
	SimpleRules      map[string]bool `json:"honnef.co/go/tools/simple"`
	StaticcheckRules map[string]bool `json:"honnef.co/go/tools/staticcheck"`
}

func main() {
	cfg := getConfig()

	multichecker.Main(
		getAnalyzers(cfg)...,
	)
}

func getConfig() *Config {
	currentFile, err := os.Executable()
	if err != nil {
		panic(err)
	}

	data, err := os.ReadFile(filepath.Join(filepath.Dir(currentFile), ConfigFileName))
	if err != nil {
		panic(err)
	}

	var cfg Config
	if err = json.Unmarshal(data, &cfg); err != nil {
		panic(err)
	}
	return &cfg
}

func getAnalyzers(cfg *Config) []*analysis.Analyzer {
	checks := []*analysis.Analyzer{
		errcheck.New(),
		noosexit.New(),
		printf.Analyzer,
		shadow.Analyzer,
		structtag.Analyzer,
	}
	for _, v := range staticcheck.Analyzers {
		if cfg.StaticcheckRules[v.Analyzer.Name] {
			checks = append(checks, v.Analyzer)
		}
	}
	for _, v := range simple.Analyzers {
		if cfg.SimpleRules[v.Analyzer.Name] {
			checks = append(checks, v.Analyzer)
		}
	}
	for _, v := range stylecheck.Analyzers {
		if cfg.StylecheckRules[v.Analyzer.Name] {
			checks = append(checks, v.Analyzer)
		}
	}
	for _, v := range quickfix.Analyzers {
		if cfg.QuickfixRules[v.Analyzer.Name] {
			checks = append(checks, v.Analyzer)
		}
	}

	return checks
}
