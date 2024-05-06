package goexit

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

// OSExitAnalyzer переменная анализатора, который проверяет есть в файле main.go в фунции main прямой вызов os.Exit
var OSExitAnalyzer = &analysis.Analyzer{
	Name: "goexit",
	Doc:  "check main func in main file if it uses direct call of os.Exit",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "main.go", nil, parser.ParseComments)
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}

	ast.Inspect(node, func(n ast.Node) bool {
		if callExpr, ok := n.(*ast.CallExpr); ok {
			if fun, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
				if id, ok := fun.X.(*ast.Ident); ok && id.Name == "os" && fun.Sel.Name == "Exit" {
					pass.Reportf(id.NamePos, "direct call of os.Exit detected")
				}
			}
		}
		return true
	})
	return nil, nil
}
