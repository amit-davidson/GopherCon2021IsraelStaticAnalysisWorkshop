package main

import (
	"errors"
	"fmt"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
	"log"
	"os"
)

var ErrNoPackages = errors.New("no packages in the path")
var ErrLoadPackages = errors.New("loading the following file contained errors")
var errMessage = "Infinite Recursion Call"
var loadMode = packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles | packages.NeedImports | packages.NeedDeps | packages.NeedExportsFile | packages.NeedTypes | packages.NeedSyntax | packages.NeedTypesInfo | packages.NeedTypesSizes

func main() {
	path := os.Args[1]
	prog, err := loadPackage(path)
	if err != nil {
		log.Fatalf("Failed to parse dssa %s: %s", path, err)
	}
	funcs := getAllFunctions(prog)
	for _, fn := range funcs {
		outputs := checkInfiniteRecursion(fn)
		for _, output := range outputs {
			fmt.Println(output)
		}
	}

}

func output(fn *ssa.Function) string {
	return fmt.Sprintf("%s: %s:%d", errMessage, fn.Name(), fn.Pos())
}

func checkInfiniteRecursion(fn *ssa.Function) string {
	for _, block := range fn.Blocks {
		for _, ins := range block.Instrs {
			if call, ok := ins.(ssa.CallInstruction); ok {
				if callCommon := call.Common().StaticCallee(); callCommon != nil {
					if callCommon != fn {
						continue
					}

					if _, ok := ins.(*ssa.Go); ok {
						// Recursively spawning goroutines doesn't consume
						// stack space infinitely, so don't flag it.
						return ""
					}

					canReturn := false
					for _, b := range fn.Blocks {
						if block.Dominates(b) {
							continue
						}
						if len(b.Instrs) == 0 {
							continue
						}
						lastInstrInBlock := b.Instrs[len(b.Instrs)-1]
						if _, ok := lastInstrInBlock.(*ssa.Return); ok {
							canReturn = true
							break
						}
					}
					if canReturn {
						return ""
					}
					return output(callCommon)
				}
			}
		}
	}
	return ""
}

func getAllFunctions(prog *ssa.Program) []*ssa.Function {
	funcs := make([]*ssa.Function, 0)
	pkgs := prog.AllPackages()
	for _, pkg := range pkgs {
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
	}
	return funcs
}

func loadPackage(path string) (*ssa.Program, error) {
	conf1 := packages.Config{
		Mode: loadMode,
	}
	loadQuery := fmt.Sprintf("file=%s", path)
	pkgs, err := packages.Load(&conf1, loadQuery)
	if err != nil {
		return nil, err
	}
	if len(pkgs) == 0 {
		return nil, fmt.Errorf("%s: %w", path, ErrNoPackages)
	}

	if len(pkgs[0].Errors) > 0 {
		return nil, fmt.Errorf("%w %s: %s", ErrLoadPackages, path, pkgs[0].Errors[0].Msg)
	}
	ssaProg, _ := ssautil.AllPackages(pkgs, 0)
	ssaProg.Build()
	return ssaProg, nil
}
