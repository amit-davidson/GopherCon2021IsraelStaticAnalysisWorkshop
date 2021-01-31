package main

import (
	"errors"
	"fmt"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
	"log"
	"os"
)

var ErrNoPackages = errors.New("no packages in the path")
var ErrLoadPackages = errors.New("loading the following file contained errors")
var errMessage = "is called in a separate goroutine, but it must be called in the same goroutine as the test"
var loadMode = packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles | packages.NeedImports | packages.NeedExportsFile | packages.NeedTypes | packages.NeedSyntax | packages.NeedTypesInfo | packages.NeedTypesSizes

func main() {
	path := os.Args[1]
	prog, pkg, err := loadPackage(path)
	if err != nil {
		log.Fatalf("Failed to parse dssa %s: %s", path, err)
	}
	funcs := getAllFunctions(pkg)
	for _, fn := range funcs {
		outputs := checkTerminationInsideGoroutine(fn, prog.Fset)
		for _, output := range outputs {
			fmt.Println(output)
		}
	}

}

func output(name string, pos token.Position) string {
	return fmt.Sprintf("T.%s %s: %s:%d:%d", name, errMessage, pos.Filename, pos.Line, pos.Column)
}

func checkTerminationInsideGoroutine(fn *ssa.Function, fset *token.FileSet) string {
	for _, block := range fn.Blocks {
		for _, ins := range block.Instrs {
			gostmt, ok := ins.(*ssa.Go)
			if !ok {
				continue
			}
			gofn := gostmt.Call.StaticCallee()
			if gofn == nil {
				continue
			}
			if gofn.Blocks == nil {
				continue
			}
			for _, block := range gofn.Blocks {
				for _, ins := range block.Instrs {
					call, ok := ins.(*ssa.Call)
					if !ok {
						continue
					}
					callee := call.Call.StaticCallee()
					if callee == nil {
						continue
					}
					recv := callee.Signature.Recv()
					if recv == nil {
						continue
					}
					if types.TypeString(recv.Type(), nil) != "*testing.common" {
						continue
					}
					terminateFn, ok := call.Call.StaticCallee().Object().(*types.Func)
					if !ok {
						continue
					}
					name := terminateFn.Name()
					switch name {
					case "FailNow", "Fatal", "Fatalf", "SkipNow", "Skip", "Skipf":
					default:
						continue
					}
					return output(name, fset.Position(call.Pos()))
				}
			}
		}
	}
	return ""
}

func getAllFunctions(pkg *ssa.Package) []*ssa.Function {
	funcs := make([]*ssa.Function, 0)
	for _, mem := range pkg.Members {
		if fdecl, ok := mem.(*ssa.Function); ok {
			var addAnons func(f *ssa.Function)
			addAnons = func(f *ssa.Function) {
				if f.Name() == "init" {
					return
				}
				funcs = append(funcs, f)
				for _, anon := range f.AnonFuncs {
					addAnons(anon)
				}
			}
			addAnons(fdecl)
		}
	}
	return funcs
}

func loadPackage(path string) (*ssa.Program, *ssa.Package, error) {
	conf1 := packages.Config{
		Mode: loadMode,
	}
	loadQuery := fmt.Sprintf("file=%s", path)
	pkgs, err := packages.Load(&conf1, loadQuery)
	if err != nil {
		return nil, nil, err
	}

	if len(pkgs) == 0 {
		return nil, nil, fmt.Errorf("%s: %w", path, ErrNoPackages)
	}

	if len(pkgs[0].Errors) > 0 {
		return nil, nil, fmt.Errorf("%w %s: %s", ErrLoadPackages, path, pkgs[0].Errors[0].Msg)
	}

	ssaProg, builtPkgs := ssautil.Packages(pkgs, 0)
	ssaProg.Build()

	return ssaProg, builtPkgs[0], nil
}
