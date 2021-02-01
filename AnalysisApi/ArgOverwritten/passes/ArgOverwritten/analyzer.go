package ArgOverwritten

import (
	"fmt"
	"go/ast"
	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "ArgOverwritten",
	Doc:  doc,
	Run:  run,
}

const (
	doc = "ArgOverwritten finds function arguments being overwritten"
)

func report(pass *analysis.Pass, ident *ast.Ident) {
	message := fmt.Sprintf("\"%s\" overwrites func parameter", ident.Name)
	pass.Report(analysis.Diagnostic{
		Pos:     ident.Pos(),
		Message: message,
	})
}

func run(pass *analysis.Pass) (interface{}, error) {
	visitor := func(node ast.Node) bool {
		var typ *ast.FuncType
		var body *ast.BlockStmt
		switch fn := node.(type) {
		case *ast.FuncDecl: // Regular function
			typ = fn.Type
			body = fn.Body
		case *ast.FuncLit: // Anonymous function
			typ = fn.Type
			body = fn.Body
		}
		if typ == nil || body == nil { // Exclude other types but also external functions with missing body
			return true
		}
		if len(typ.Params.List) == 0 {
			return true
		}

		for _, field := range typ.Params.List {
			for _, arg := range field.Names {
				obj := pass.TypesInfo.ObjectOf(arg)
				ast.Inspect(body, func(node ast.Node) bool {
					switch stmt := node.(type) {
					case *ast.AssignStmt:
						for _, lhs := range stmt.Lhs {
							ident, ok := lhs.(*ast.Ident)
							if ok && pass.TypesInfo.ObjectOf(ident) == obj {
								report(pass, ident)
							}
						}
					case *ast.IncDecStmt:
						ident, ok := stmt.X.(*ast.Ident)
						if ok && pass.TypesInfo.ObjectOf(ident) == obj {
							report(pass, ident)
						}
					}
					return true
				})
			}
		}
		return true
	}
	for _, f := range pass.Files {
		ast.Inspect(f, visitor)
	}
	return nil, nil
}
