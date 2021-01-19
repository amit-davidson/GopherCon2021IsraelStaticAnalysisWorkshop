package main

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"log"
	"os"
)

func main() {
	fset := token.NewFileSet()
	path := os.Args[1]
	pkgs, err := parser.ParseDir(fset, path, nil, 0)
	if err != nil {
		log.Fatalf("Failed to parse dir %s: %s", path, err)
	}

	for _, pkg := range pkgs {
		analyzePackage(pkg, fset)
	}
}

func analyzePackage(p *ast.Package, fset *token.FileSet) {
	conf := types.Config{Importer: importer.Default()}
	info := &types.Info{
		Defs: make(map[*ast.Ident]types.Object),
		Uses: make(map[*ast.Ident]types.Object),
	}

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
				obj := info.ObjectOf(arg)
				ast.Inspect(body, func(node ast.Node) bool {
					switch stmt := node.(type) {
					case *ast.AssignStmt:
						for _, lhs := range stmt.Lhs {
							ident, ok := lhs.(*ast.Ident)
							if ok && info.ObjectOf(ident) == obj {
								output(ident, fset.Position(ident.Pos()))
							}
						}
					case *ast.IncDecStmt:
						ident, ok := stmt.X.(*ast.Ident)
						if ok && info.ObjectOf(ident) == obj {
							output(ident, fset.Position(ident.Pos()))
						}
					}
					return true
				})
			}
		}
		return true
	}
	for _, f := range p.Files {
		_, err := conf.Check("fib", fset, []*ast.File{f}, info)
		if err != nil {
			log.Fatal(err)
		}
		ast.Inspect(f, visitor)
	}

}

func output(ident *ast.Ident, pos token.Position) {
	fmt.Printf("\"%s\" overwrites func parameter in pos: %s:%d:%d\n", ident.Name, pos.Filename, pos.Line, pos.Column)
}
