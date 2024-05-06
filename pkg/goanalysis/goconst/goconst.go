package goconst

import (
	"fmt"
	"go/token"

	goconstAPI "github.com/jgautheron/goconst"
	"golang.org/x/tools/go/analysis"
)

// GoConstSettings Структура настроек для анализатора GoConst
type GoConstSettings struct {
	IgnoreStrings       string `mapstructure:"ignore-strings"`
	IgnoreTests         bool   `mapstructure:"ignore-tests"`
	MatchWithConstants  bool   `mapstructure:"match-constant"`
	MinStringLen        int    `mapstructure:"min-len"`
	MinOccurrencesCount int    `mapstructure:"min-occurrences"`
	ParseNumbers        bool   `mapstructure:"numbers"`
	NumberMin           int    `mapstructure:"min"`
	NumberMax           int    `mapstructure:"max"`
	IgnoreCalls         bool   `mapstructure:"ignore-calls"`
}


// NewAnalyzer Создает новый анализатора на основе пакета jgautheron/goconst
func NewAnalyzer() *analysis.Analyzer {
	goconstSettings := GoConstSettings{
		MatchWithConstants:  true,
		MinStringLen:        3,
		MinOccurrencesCount: 3,
		NumberMin:           3,
		NumberMax:           3,
		IgnoreCalls:         true,
	}

	analyzer := &analysis.Analyzer{
		Name: "goconst",
		Doc:  "Finds repeated strings that could be replaced by a constant.",
		Run: func(pass *analysis.Pass) (any, error) {
			err := runGoconst(pass, &goconstSettings)
			if err != nil {
				return nil, err
			}

			return nil, nil
		},
	}
	return analyzer
}

func runGoconst(pass *analysis.Pass, settings *GoConstSettings) error {
	cfg := goconstAPI.Config{
		IgnoreStrings:      settings.IgnoreStrings,
		IgnoreTests:        settings.IgnoreTests,
		MatchWithConstants: settings.MatchWithConstants,
		MinStringLength:    settings.MinStringLen,
		MinOccurrences:     settings.MinOccurrencesCount,
		ParseNumbers:       settings.ParseNumbers,
		NumberMin:          settings.NumberMin,
		NumberMax:          settings.NumberMax,
		ExcludeTypes:       map[goconstAPI.Type]bool{},
	}
	issues, err := goconstAPI.Run(pass.Files, pass.Fset, &cfg)
	if err != nil {
		return err
	}

	for _, i := range issues {
		text := fmt.Sprintf("string %s has %d occurrences", i.Str, i.OccurrencesCount)

		if i.MatchingConst == "" {
			text += ", make it a constant"
		} else {
			text += fmt.Sprintf(", but such constant %s already exists", i.MatchingConst)
		}

		pass.Reportf(token.Pos(i.Pos.Offset), text)
	}

	return nil
}
