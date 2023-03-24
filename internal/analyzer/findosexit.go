package analyzer

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
)

var ErrCheckAnalyzer = &analysis.Analyzer{
	Name: "findOsExit",
	Doc:  "find os exit",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			switch x := node.(type) {
			case *ast.SelectorExpr:
				if call, ok := x.X.(*ast.Ident); ok {
					if call.Name == "os" && x.Sel.Name == "Exit" {
						pass.Reportf(call.NamePos, "do not use os.Exit() in main")
					}
				}
			}
			return true
		})
	}
	return nil, nil
}
