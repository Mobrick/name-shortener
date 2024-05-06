package main

import (
	"strings"

	"github.com/Mobrick/name-shortener/cmd/staticlint/goexit"
	"github.com/Mobrick/name-shortener/pkg/goanalysis/goconst"
	"github.com/Mobrick/name-shortener/pkg/goanalysis/prealloc"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/staticcheck"
)

func main() {

	mychecks := []*analysis.Analyzer{
		printf.Analyzer,
		shadow.Analyzer,
		structtag.Analyzer,
		goconst.NewAnalyzer(),
		prealloc.NewAnalyzer(),
		goexit.OSExitAnalyzer,
	}

	for _, v := range staticcheck.Analyzers {
		switch {
		case strings.HasPrefix(v.Analyzer.Name, "SA"):
		case strings.HasPrefix(v.Analyzer.Name, "S1"):
			mychecks = append(mychecks, v.Analyzer)
		}
	}

	multichecker.Main(
		mychecks...,
	)
}
