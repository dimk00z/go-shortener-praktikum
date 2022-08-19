// Package for staticlint check
// For using dhould be compiled|installed
// run staticlint <path_to_files>
package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"staticlint/analyzers"

	critic "github.com/go-critic/go-critic/checkers/analyzer"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/analysis/lint"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"
)

// Config — config path.
const Config = `config/config.json`
const (
	staticAnylPref = "SA"
	styleAnylPref  = "ST"
)

// ConfigData describes struct of config.
type ConfigData struct {
	Staticcheck []string
	Stylecheck  []string
}

var DefaultConfig = ConfigData{
	Staticcheck: []string{"SA"},
	Stylecheck:  []string{"ST"},
}

// Load analyzers by config
func loadRequiredSaticChecks(cfg ConfigData) []*analysis.Analyzer {
	requiredChecks := make([]*analysis.Analyzer, 0)

	analyzers := make(map[string][]*lint.Analyzer)
	analyzers[staticAnylPref] = staticcheck.Analyzers
	analyzers[styleAnylPref] = stylecheck.Analyzers

	configOpt := make(map[string][]string)
	configOpt[staticAnylPref] = cfg.Staticcheck
	configOpt[styleAnylPref] = cfg.Stylecheck

	for _, checkType := range []string{staticAnylPref, styleAnylPref} {
		if len(configOpt[checkType]) == 1 && configOpt[checkType][0] == checkType {
			for _, a := range analyzers[checkType] {
				requiredChecks = append(requiredChecks, a.Analyzer)
			}
		} else {
			checks := make(map[string]bool)
			for _, v := range configOpt[checkType] {
				checks[v] = true
			}
			for _, v := range analyzers[checkType] {
				if checks[v.Analyzer.Name] {
					requiredChecks = append(requiredChecks, v.Analyzer)
				}
			}
		}
	}
	return requiredChecks
}

func main() {
	log.Println("Staticlint started")

	appfile, err := os.Executable()
	if err != nil {
		log.Fatalf("os.Executable - %v ", err)
	}
	var cfg ConfigData

	data, err := os.ReadFile(filepath.Join(filepath.Dir(appfile), Config))
	if err != nil {
		log.Printf("os.ReadFile - %v", err)
		log.Printf("Loaded default config - %v", err)
		cfg = DefaultConfig
	} else {
		if err = json.Unmarshal(data, &cfg); err != nil {
			log.Printf("json.Unmarshal - %v", err)
		}
	}

	staticlintChecks := []*analysis.Analyzer{
		printf.Analyzer,
		shadow.Analyzer,
		shift.Analyzer,
		structtag.Analyzer,
		critic.Analyzer,
		analyzers.OsExitAnalyzer,
	}
	staticlintChecks = append(staticlintChecks, loadRequiredSaticChecks(cfg)...)

	log.Printf("Loaded current analyzers - %v", staticlintChecks)
	multichecker.Main(staticlintChecks...)

}
