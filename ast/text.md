 - AST
	 1.1 Go packages overview
	 There are 6 relevant packages:
	 - ast - Package ast declares the types used to represent syntax trees for Go packages.
	 - token - Package token defines constants representing the lexical tokens of Go 
	 - scanner - Package scanner implements a scanner for Go source text. It takes a `[]byte` as source which can then be tokenized through repeated calls to the Scan method.
	 - parser - Package parser implements a parser for Go source files. The output is an abstract syntax tree (AST) representing the Go source
	 - constant - Package constant implements Values representing untyped Go constants and their corresponding operations.
	 - types - Package types declares the data types and implements the algorithms for type-checking of Go packages

The `scanner` package is fed with `[]byte` representing the source code. It's output is a list of `tokens` defined by the `token` package and they are used by the `parser` package to create the `ast` tree. After the tree is constructed, type checking algorithms defined by `types` run over the tree, validates it correctness and evalutes `constants` accordingly. 

![Go packages overview](https://i.imgur.com/xo2stvz.png)
1.2 What is AST?
		 An abstract syntax tree (AST) is a way of representing the syntax of a 		programming language as a hierarchical tree-like structure. Let's take a look at the following program for explanation. 
	
  ```
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

Let's focus on the JSON under `*ast.File` represetnting a Go source file. The file is the root node and it contains all the top-level declarations in the file - the import and the main function declarations. Under the `main` body we have a `blockStmt` containing a list of the function statements. As you can see, it resembles a tree-like structure.

The syntax is "abstract" in the sense that it does not represent every detail appearing in the real syntax, but rather just the structural or content-related details. For instance, grouping parentheses are implicit in the tree structure, so these are not represented as separate nodes.
![Go packages overview](https://i.imgur.com/oGuNoQZ.png)

1.3 AST package members
The `ast` package contains the types used to represent syntax trees in Go.
We can divide the types into 3 categories: Interfaces, concrete types and others.

Interfaces:	`Node`, `Decl`, `Spec`, `Stmt`, `Expr`
Concrete Types: I won't go over the full list, but it's possible to see it [here](https://golang.org/pkg/go/ast/#pkg-index). In short, those are the nodes of the tree and it contains value such as: `FuncDecl`. `IncDecStmt`, `Ident`, `Comment` and so on.
Others: `Package`, `File`, `Scope`, `Object`

As you can see, the `ast` contains only the "abstract" parts and ignores parentheses, colon, etc...

1.4 Loading a program using the parser
	As explained we need to parse the file first.
	
  ```
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
         fmt.Println("We found hello world!")
         return false  
      }  
      return true  
  }  
   ast.Inspect(f, visitor)  
}
  ```

We first create a `fileSet` which represents a set of source files. `FileSet` has the properties `files` and `base`  recording the files and the total size of all the files respectively. Using the `size` property of `token.File` we can easily determine in what file a statement is given it's position.

Then, we call the `parser.ParseFile` function, providing it our `fileSet` so it can populate it, an empty `path`, a string as the source so the parser will use it instead of loading from a file, and a build mode - `0`	. For this example we used `0`  to load the program fully, but any [other mode](https://golang.org/pkg/go/parser/#Mode) can be used.

![FileSet](https://i.imgur.com/wrO1T0V.png)

> Tip: Instead of iterating file by file, you can load an entire directory using `parser.ParseDir`


Finally, we call `ast.Inspect` to iterate over all the nodes in depth-first order and print a message when we reach the `Hello World` string literal. `true` is returned each iteration to keep traversing the tree until the desired node is found. Then, we return `false` to indicate we're done searching. 

1.5 Writing our first analyzer!
	
