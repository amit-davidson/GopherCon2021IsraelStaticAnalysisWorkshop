## Compiler front end, AST, and analysis introduction
### 2.1 Go packages overview
There are six relevant packages regarding the compiler front end when talking about static analysis:

- [token](https://golang.org/pkg/go/token/) - Package `token` defines constants representing the lexical tokens of Go
- [scanner](https://golang.org/pkg/go/scanner/) - Package `scanner` implements a scanner for Go source text. It takes a `[]byte` as the source, which can then be tokenized through repeated calls to the `Scan` method.
- [parser](https://golang.org/pkg/go/parser/) - Package `parser` implements a parser for Go source files. The output is an abstract syntax tree (AST) representing the Go source
- [AST](https://golang.org/pkg/go/ast/) - Package `AST` declares the types used to represent syntax trees for Go packages.
- [constant](https://golang.org/pkg/go/constant/) - Package `constant` implements Values representing untyped Go constants and their corresponding operations.
- [types](https://golang.org/pkg/go/types/) - Package `types` declares the data types and implements the algorithms for type-checking of Go packages

The `scanner` package is fed with `[]byte` representing the source code. Its output is a list of tokens defined by the
`token` package, and the parser package uses them to create the `AST` tree. After the tree is constructed,
the parser runs type-checking algorithms run over the tree, validates its correctness, and evaluates constants.

### 2.2 What is AST?
An abstract syntax tree (AST) is a way of representing the syntax of a programming language as a hierarchical tree-like structure. Let's take a look at the following program for an explanation.

``` go
package main
import "fmt"

func main() {
  fmt.Println("hello world")
}
```

We can use this [AST visualizer](http://goast.yuroyoro.net/) to view it's AST.
```
     0  *ast.File {
     1  .  Package: 1:1
     2  .  Name: *ast.Ident {
     3  .  .  NamePos: 1:9
     4  .  .  Name: "main"
     5  .  }
     6  .  Decls: []ast.Decl (len = 2) {
     7  .  .  0: *ast.GenDecl {
     8  .  .  .  TokPos: 3:1
     9  .  .  .  Tok: import
    10  .  .  .  Lparen: -
    11  .  .  .  Specs: []ast.Spec (len = 1) {
    12  .  .  .  .  0: *ast.ImportSpec {
    13  .  .  .  .  .  Path: *ast.BasicLit {
    14  .  .  .  .  .  .  ValuePos: 3:8
    15  .  .  .  .  .  .  Kind: STRING
    16  .  .  .  .  .  .  Value: "\"fmt\""
    17  .  .  .  .  .  }
    18  .  .  .  .  .  EndPos: -
    19  .  .  .  .  }
    20  .  .  .  }
    21  .  .  .  Rparen: -
    22  .  .  }
    23  .  .  1: *ast.FuncDecl {
    24  .  .  .  Name: *ast.Ident {
    25  .  .  .  .  NamePos: 5:6
    26  .  .  .  .  Name: "main"
    27  .  .  .  .  Obj: *ast.Object {
    28  .  .  .  .  .  Kind: func
    29  .  .  .  .  .  Name: "main"
    30  .  .  .  .  .  Decl: *(obj @ 23)
    31  .  .  .  .  }
    32  .  .  .  }
    33  .  .  .  Type: *ast.FuncType {
    34  .  .  .  .  Func: 5:1
    35  .  .  .  .  Params: *ast.FieldList {
    36  .  .  .  .  .  Opening: 5:10
    37  .  .  .  .  .  Closing: 5:11
    38  .  .  .  .  }
    39  .  .  .  }
    40  .  .  .  Body: *ast.BlockStmt {
    41  .  .  .  .  Lbrace: 5:13
    42  .  .  .  .  List: []ast.Stmt (len = 1) {
    43  .  .  .  .  .  0: *ast.ExprStmt {
    44  .  .  .  .  .  .  X: *ast.CallExpr {
    45  .  .  .  .  .  .  .  Fun: *ast.SelectorExpr {
    46  .  .  .  .  .  .  .  .  X: *ast.Ident {
    47  .  .  .  .  .  .  .  .  .  NamePos: 6:3
    48  .  .  .  .  .  .  .  .  .  Name: "fmt"
    49  .  .  .  .  .  .  .  .  }
    50  .  .  .  .  .  .  .  .  Sel: *ast.Ident {
    51  .  .  .  .  .  .  .  .  .  NamePos: 6:7
    52  .  .  .  .  .  .  .  .  .  Name: "Println"
    53  .  .  .  .  .  .  .  .  }
    54  .  .  .  .  .  .  .  }
    55  .  .  .  .  .  .  .  Lparen: 6:14
    56  .  .  .  .  .  .  .  Args: []ast.Expr (len = 1) {
    57  .  .  .  .  .  .  .  .  0: *ast.BasicLit {
    58  .  .  .  .  .  .  .  .  .  ValuePos: 6:15
    59  .  .  .  .  .  .  .  .  .  Kind: STRING
    60  .  .  .  .  .  .  .  .  .  Value: "\"hello world\""
    61  .  .  .  .  .  .  .  .  }
    62  .  .  .  .  .  .  .  }
    63  .  .  .  .  .  .  .  Ellipsis: -
    64  .  .  .  .  .  .  .  Rparen: 6:28
    65  .  .  .  .  .  .  }
    66  .  .  .  .  .  }
    67  .  .  .  .  }
    68  .  .  .  .  Rbrace: 7:1
    69  .  .  .  }
    70  .  .  }
    71  .  }
    72  .  Scope: *ast.Scope {
    73  .  .  Objects: map[string]*ast.Object (len = 1) {
    74  .  .  .  "main": *(obj @ 27)
    75  .  .  }
    76  .  }
    77  .  Imports: []*ast.ImportSpec (len = 1) {
    78  .  .  0: *(obj @ 12)
    79  .  }
    80  .  Unresolved: []*ast.Ident (len = 1) {
    81  .  .  0: *(obj @ 46)
    82  .  }
    83  }
```

Let's focus on the JSON under `*ast.File` representing a Go source file. The file is the root node, and it contains all
the top-level declarations in the file - the import and the main function. Under `mains'` body, we have a
`blockStmt` containing a list of the function statements. Similar to HTML, the dependency of the nodes create a
tree-like structure. 

The syntax is "abstract" in the sense that it does not represent every detail appearing in the real syntax, but rather
just the structural or content-related details. For instance, grouping parentheses are implicit in the tree structure,
so these are not represented as separate nodes.

### 2.3 AST package members
The AST package contains the types used to represent syntax trees in Go. We can divide the members into three categories:
Interfaces, concrete types, and others.

- Concrete Types: The full list is [long](https://golang.org/pkg/go/ast/#ArrayType). Those are the tree nodes, and they contain values such as: `FuncDecl`, `IncDecStmt`, `Ident`, `Comment`, and so on.
- Interfaces: `Node`, `Decl`, `Spec`, `Stmt`, `Expr`
- Others: `Package`, `File`, `Scope`, `Object` 

Well take a look at `AssignStmt`
```go
type AssignStmt struct {
    Lhs    []Expr      // All the variables to left side of assign operator
    TokPos token.Pos   // position of operator
    Tok    token.Token // assignment token. `=`, `:=`, `+=`, `<<=` and so on...
    Rhs    []Expr      // Expressions to right of the assignment operator 
}
```
For example, in the expression: `a := 5`, 
 - Lhs is [a]
 - Rhs is [5]
 - TokPos [3] (the position of the ":" character)
 - Tok [:=]
 
Pretty straight forward. Now we'll look at `Expr` which `AssignStmt` implements. `Expr` is common interface for 
everything that returns a value. As you can see, it only contains the node interface (which is implemented by all the nodes on the AST graph).

```go
type Expr interface {
    Node
    // contains filtered or unexported methods
}
```

From the other's group we'll look at `File`.
```go
type File struct {
    Doc        *CommentGroup   // associated documentation; or nil
    Package    token.Pos       // position of "package" keyword
    Name       *Ident          // package name
    Decls      []Decl          // top-level declarations; or nil
    Scope      *Scope          // package scope (this file only)
    Imports    []*ImportSpec   // imports in this file
    Unresolved []*Ident        // unresolved identifiers in this file
    Comments   []*CommentGroup // list of all comments in the source file
}

```

An `ast.File` is the root node of each of the files we analyze. When analyzing a program, we'll iterate over the files 
and for each we'll pass it to the `ast.Inspect` function for iteration.   

It's worth mentioning again that `ast` package contains only the "abstract" parts so it ignores parentheses, colon, etc...

### 2.4 Exercise:
In the folder `/ast/CodeExamples` there are some interesting programs (well... AST-wise). Using our [AST visualizer](http://goast.yuroyoro.net/)
from earlier, take each of the program and look at their AST. I added comments explaining the important points.   

### 2.5 Loading a program using the parser
To load the program, we need to parse it first
``` go
package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

func main() {

	src := `package main  

import "fmt"  

func main() {  
   fmt.Println("hello world")}  
`

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		fmt.Println(err)
		return
	}
	visitor := func(node ast.Node) bool {
		strLit, ok := (node).(*ast.BasicLit)
		if ok && strLit.Value == "\"hello world\"" {
			pos := fset.Position(strLit.Pos())
			fmt.Printf("We found hello world in pos:%d:%d", pos.Line, pos.Column)
			return false
		}
		return true
	}
	ast.Inspect(f, visitor)
}
```


We first create a `fileSet`, representing a set of source files. `FileSet` has the properties `files` and `base`
recording the used files and all the files' total size. Using the size property of `token.File`, we can easily determine
in which file a statement is, given its position.

```go
fset := token.NewFileSet()  
```

<img src="https://i.imgur.com/AQfkL3E.png" height="50%" width="50%"/>

Then, we call the `parser.ParseFile` function, providing it our `fileSet` to populate it, an empty path, a string as the
source so the parser will use it instead of loading from a file, and a build mode - 0. In this example, we used 0 to 
fully load the program, but any other [mode](https://golang.org/pkg/go/parser/#Mode) can be used.

```go
f, err := parser.ParseFile(fset, "", src, 0)  
 if err != nil {  
    fmt.Println(err)  
    return  
}  
```
> Tip: Instead of iterating file by file, you can load an entire directory using `parser.ParseDir`

Finally, we define a visitor function that will be called with each node inside the AST. We pass our function to
`ast.Inspect` to iterate over all the nodes in depth-first order and print a message when we reach the
`Hello World` string literal with it's position in the code. We return `true` each iteration to keep traversing the tree until we found the desired 
node. Then, we print our message and return false to indicate we're done searching and to exit the traverse function.

```go
visitor := func(node ast.Node) bool {
    strLit, ok := (node).(*ast.BasicLit)
    if ok && strLit.Value == "\"hello world\"" {
        pos := fset.Position(strLit.Pos())
        fmt.Printf("We found hello world in pos:%d:%d", pos.Line, pos.Column)
        return false
    }
    return true
}
ast.Inspect(f, visitor)
```

### 2.6 Exercise 2!
In the file `ast/ArgsOverwriteAnalyzer.go` we have an analyzer that checks if function arguments were modified as in the 
example below. The problem is that there are some parts missing from it. There are comments in the places where you
should add your code according to the comment. You can run the tests `ast/ArgsOverwriteAnalyzer_test.go` to make sure your
tests pass. You can also debug using the test to inspect the AST graph of this program.
If you give up, you can see the result in `ast/result/ArgsOverwriteAnalyzer.go` :)


### 2.7 Congratulations
You have a good understanding of what AST is, the different Go packages used to create static code analyzers that 
interact with it and how to write such analyzers.  

In the [next section](https://github.com/amit-davidson/GopherCon2021IsraelStaticAnalysisWorkshop/blob/master/ir/text.md) 
we'll focus on the middle end level, and see how analyzer "operating" in this level work.
