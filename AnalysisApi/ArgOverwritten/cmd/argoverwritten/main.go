package main

import (
	"github.com/amit-davidson/ArgOverwritten/passes/ArgOverwritten"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(ArgOverwritten.Analyzer)
}
