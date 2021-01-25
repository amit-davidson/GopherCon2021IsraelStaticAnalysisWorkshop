## 2. IR
### 2.1 What is IR?
An **intermediate representation** (**IR**) is the code used internally by a compiler to represent source code. An IR is designed to be conducive for further processing, such as optimization and translation. A "good" IR must be _accurate_ – capable of representing the source code without loss of information  – and _independent_ of any particular source or target language.
   
### 2.2 What is SSA?
SSA stands for static single assignment. It's a property of an IR **that requires each variable to be assigned exactly once**, and every variable be defined before it is used. 
The primary usefulness of SSA comes from how it simplifies the properties of variables and improves compilers optimizations.

For example, consider this piece of code:
```
y := 1
y := 2
x := y
```

In an SSA form, it'll be translated to:
```
y1 := 1
y2 := 2
x1 := y2
```
Humans can see that the first assignment is unnecessary and that the value of `y`  used in the third line comes from the
second assignment of `y`. In SSA form, both of these are immediate

### 2.3 SSA package members
The package tools/go/ssa defines the representation of elements of Go programs in SSA format.
The key types form a hierarchical structure.

Program - A Program is a partial or complete Go program converted to an SSA form.

<img src="https://i.imgur.com/DpzHQib.png" width="50%" height="50%" />

Package - A Package is a single analyzed Go package containing Members for all package-level functions, variables, constants, and types it declares.

<img src="https://i.imgur.com/stQ9izj.png" width="50%" height="50%" />

Function - Function represents the parameters, results, and code of a function or method.

<img src="https://i.imgur.com/5KLBY6r.png" width="50%" height="50%" />

Basic Block - BasicBlock represents an SSA basic block. A set of instructions that are executed and can't jump somewhere else. Basic blocks are connected using conditions and goto statements.
 
<img src="https://i.imgur.com/dBLj172.png" width="50%" height="50%" />

Control Flow Graph (CFG) - In a control-flow graph, each node in the graph represents a basic block. Together, they compose all paths that might be traversed through a program during its execution.

<img src="https://i.imgur.com/K1u4MZ0.png" width="50%" height="50%" />

Instruction - a statement that consumes values and performs computation. For example, `Call`, `Return`, `TypeAssert`, etc

<img src="https://i.imgur.com/DvheFlc.png" width="50%" height="50%" />

Value - an expression that yields a value. For example, Function calls are both `Instruction` and `Value` since they both consume values and yield a value.

<img src="https://i.imgur.com/oJg97Re.png" width="50%" height="50%" />

And when combined:

<img src="https://i.imgur.com/W02MErA.png" width="70%" height="70%" />

The package contains other [types](https://pkg.go.dev/golang.org/x/tools/go/ssa#pkg-overview) - Include language keywords such as `Defer`, `If` but also lower level primitives like `MakeChan` and `Alloc`. 

### 2.4 Viewing SSA
We can use this  [SSA visualizer](http://golang-ssaview.herokuapp.com/)  to view the SSA form of programs.

> You can also use `go.tools/cmd/ssadump` in view SSA in your CLI

Let's consider this program:
```
package main

import (
    "fmt"
	"math"
	"os"
)

func main() {
	shapeType := os.Args[1]
	squareOrCircleArea(shapeType)
}

func squareOrCircleArea(shapeType string) {
	r := 2.0
	area := r * r
	if shapeType == "circle" {
		area *= math.Pi
	}
	fmt.Printf("Total area is: %g", area)
}
```

I'll focus on the `squareOrCircleArea` function, but you can inspect the full ssa program in the visualizer.
```go
func squareOrCircleArea(shapeType string):
0:                                                                entry P:0 S:2
        t0 = 2:float64 * 2:float64                                      float64
        t1 = shapeType == "circle":string                                  bool
        if t1 goto 1 else 2
1:                                                              if.then P:1 S:1
        t2 = t0 * 3.14159:float64                                       float64
        jump 2
2:                                                              if.done P:2 S:0
        t3 = phi [0: t0, 1: t2] #area                                   float64
        t4 = new [1]interface{} (varargs)                       *[1]interface{}
        t5 = &t4[0:int]                                            *interface{}
        t6 = make interface{} <- float64 (t3)                       interface{}
        *t5 = t6
        t7 = slice t4[:]                                          []interface{}
        t8 = fmt.Printf("Total area is: %g":string, t7...)   (n int, err error)
        return
```

Looking at the first basic block (0) we can see straight away that the variable names were replaced with `t` followed by a number.
Also, the assignment to `r` is missing and it's values are already used in the assignment to `area`. This is the result 
of constant propagation and dead code elimination indicating this code is already optimized.

In the end of the block, we can see a conditional goto (as opposed to the conventional if structure) to the correct
basic block, according to the shape type.

```go
0:                                                                entry P:0 S:2
        t0 = 2:float64 * 2:float64                                      float64
        t1 = shapeType == "circle":string                                  bool
        if t1 goto 1 else 2
```
In the source code, we multiply the value of area with PI and assign it back to the area. In SSA form, each variable is 
assigned once. We can see that `t0` is no longer used and instead `t2` is declared, even though in high level they point
to the same variable.   
```go
1:
        t2 = t0 * 3.14159:float64                                       float64
        jump 2
```

In the last block we see an instruction called `phi`. This instruction represents an SSA φ-node which combines values
that differ across incoming control-flow edges and yields a new value. We won't delve deeper, but in short, it says
the value can be either `t0` or `t2`, depending on the control flow.

At that point, we're ready to print the variable, but there are many instructions between `t3` and the `fmt.Printf` function.
IR is much more verbose and includes instructions that may by represented with single "action" in the source code. In 
this case, `fmt.Printf` uses variadic parameters. Behind the scenes, we have to declare a list of interfaces, convert 
our `float64` to the `interface{}` type and only then pass it to the function.   
```go
2:                                                              if.done P:2 S:0
        t3 = phi [0: t0, 1: t2] #area                                   float64
        t4 = new [1]interface{} (varargs)                       *[1]interface{}
        t5 = &t4[0:int]                                            *interface{}
        t6 = make interface{} <- float64 (t3)                       interface{}
        *t5 = t6
        t7 = slice t4[:]                                          []interface{}
        t8 = fmt.Printf("Total area is: %g":string, t7...)   (n int, err error)
        return
```

### Exercise
In the folder CodeExamples there are some interesting programs. Using our SSA visualizer from earlier, take each of 
the program and look at their SSA. I added comments with notes with explaining the important points. 


### 2.5 SSA vs AST
AST shows us the structure of the code. How different statements in the code relate to each other. SSA, on the other hand, shows us how the code flows.

When applying this logic to static analysis, we'll see that SSA is used for more complex analysis where we need to determine the flow of the data. In contrast, AST will be used for simpler, more structure related analyses.

### 2.6 Writing our analyzer!
https://github.com/ipfs/go-ipfs/issues/2043
