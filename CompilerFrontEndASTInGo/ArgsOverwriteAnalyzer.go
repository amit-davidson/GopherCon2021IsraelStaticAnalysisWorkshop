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
	path := os.Args[1]
	fset, pkgs, err := loadProgram(path)
	if err != nil {
		log.Fatalf("Failed to parse dir %s: %s", path, err)
	}

	for _, pkg := range pkgs {
		outputs := analyzePackage(pkg, fset)
		for _, message := range outputs {
			fmt.Println(message)
		}
	}
}

func loadProgram(path string) (*token.FileSet, map[string]*ast.Package, error) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, path, nil, 0)
	if err != nil {
		return nil, nil, err
	}
	return fset, pkgs, nil
}

func populateTypes(conf types.Config, fset *token.FileSet, f *ast.File) (*types.Info, error) {
	info := &types.Info{
		Defs: make(map[*ast.Ident]types.Object),
		Uses: make(map[*ast.Ident]types.Object),
	}
	_, err := conf.Check("", fset, []*ast.File{f}, info)
	if err != nil {
		return nil, err
	}
	return info, nil
}

func analyzePackage(p *ast.Package, fset *token.FileSet) []string {
	var info *types.Info
	var err error
	outputs := make([]string, 0)

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
						var ident *ast.Ident
						// -------------------------------------------------------------------------
						// info.ObjectOf requires *ast.ident. How can I extract *ast.ident from *ast.AssignStmt.Lhs?
						// The flow should be similar to how *ast.Ident is taken in the *ast.IncDecStmt case
						// -------------------------------------------------------------------------
						if info.ObjectOf(ident) == obj {
							outputs = append(outputs, output(ident, fset.Position(ident.Pos())))
						}

					case *ast.IncDecStmt:
						ident, ok := stmt.X.(*ast.Ident)
						if ok && info.ObjectOf(ident) == obj {
							outputs = append(outputs, output(ident, fset.Position(ident.Pos())))
						}
					}
					return true
				})
			}
		}
		return true
	}

	conf := types.Config{Importer: importer.Default()}
	for _, f := range p.Files {
		info, err = populateTypes(conf, fset, f)
		if err != nil {
			log.Fatal(err)
		}
		ast.Inspect(f, visitor)
	}
	return outputs
}

func output(ident *ast.Ident, pos token.Position) string {
	return fmt.Sprintf("\"%s\" overwrites func parameter in pos: %s:%d:%d", ident.Name, pos.Filename, pos.Line, pos.Column)
}
