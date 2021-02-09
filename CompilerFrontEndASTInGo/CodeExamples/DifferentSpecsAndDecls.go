package main

const (
	a = iota
	b
	c = 3
)

type d int

var (
	e = 7
)

func main() {
	const name = 8
	name2 = 7
}

// 1. See how each declaration is represented in the AST and the difference between them. Specifically look which type
// implements each interface (spec or decl) by looking at the docs for go/ast: https://golang.org/pkg/go/ast/.

// You can also look at the code, it gives a good overview: https://golang.org/src/go/ast/ast.go?s=26852:26891#L851
// and here:https://golang.org/src/go/ast/ast.go?s=29093:29170#L929

//2. Look at the difference between the "name" and "name2".
