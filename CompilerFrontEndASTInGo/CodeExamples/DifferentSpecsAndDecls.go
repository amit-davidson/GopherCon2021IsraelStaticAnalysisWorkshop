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

// Look here:https://golang.org/src/go/ast/ast.go?s=26852:26891#L851 and here:https://golang.org/src/go/ast/ast.go?s=29093:29170#L929
// to see the difference between spec and decl. See what type falls under each category

//2. Look at the difference between the "name" and "name2".
