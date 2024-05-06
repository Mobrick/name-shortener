package prealloc

import (
	"fmt"

	"github.com/alexkohler/prealloc/pkg"
	"golang.org/x/tools/go/analysis"
)

// PreallocSettings Структура настроек для анализатора prealloc
type PreallocSettings struct {
	Simple     bool
	RangeLoops bool `mapstructure:"range-loops"`
	ForLoops   bool `mapstructure:"for-loops"`
}

// NewAnalyzer Создает новый анализатора на основе пакета alexkohler/prealloc
func NewAnalyzer() *analysis.Analyzer {

	preallocSettings := PreallocSettings{
		Simple:     true,
		RangeLoops: true,
		ForLoops:   false,
	}

	analyzer := &analysis.Analyzer{
		Name: "prealloc",
		Doc:  "Prealloc is a Go static analysis tool to find slice declarations that could potentially be preallocated.",
		Run: func(pass *analysis.Pass) (any, error) {
			runPreAlloc(pass, &preallocSettings)
			return nil, nil
		},
	}
	return analyzer
}

func runPreAlloc(pass *analysis.Pass, settings *PreallocSettings) {
	hints := pkg.Check(pass.Files, settings.Simple, settings.RangeLoops, settings.ForLoops)

	for _, hint := range hints {
		pass.Reportf(hint.Pos, fmt.Sprintf("Consider pre-allocating %s", hint.DeclaredSliceName))
	}
}
