## 3. Compiler middle end and static analysis with SSA In Go   
### 3.1 What is SSA?
[SSA](https://en.wikipedia.org/wiki/Static_single_assignment_form) stands for static single assignment. It's a property of an IR **that requires each variable to be assigned exactly once**, and every variable be defined before it is used. 
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

### 3.2 SSA package members
The package [`tools/go/ssa`](https://pkg.go.dev/golang.org/x/tools/go/ssa) defines the representation of elements of Go programs in SSA format.
The key types form a hierarchical structure.

#### [Program](https://pkg.go.dev/golang.org/x/tools/go/ssa#Program) 
A Program is a partial or complete Go program converted to an SSA form.

<img src="https://i.imgur.com/DpzHQib.png" width="50%" height="50%" />

#### [Package](https://pkg.go.dev/golang.org/x/tools/go/ssa#Package) 
A Package is a single analyzed Go package containing Members for all package-level functions, variables, constants, and types it declares.

<img src="https://i.imgur.com/stQ9izj.png" width="50%" height="50%" />

#### [Function](https://pkg.go.dev/golang.org/x/tools/go/ssa#Function)
Function represents the parameters, results, and code of a function or method.

<img src="https://i.imgur.com/5KLBY6r.png" width="50%" height="50%" />

#### [Basic Block](https://pkg.go.dev/golang.org/x/tools/go/ssa#BasicBlock)
BasicBlock represents an SSA basic block. A set of instructions that are executed and can't jump somewhere else. Basic blocks are connected using conditions and goto statements.
 
<img src="https://i.imgur.com/dBLj172.png" width="50%" height="50%" />

Control Flow Graph (CFG) - In a control-flow graph, each node in the graph represents a basic block.
Together, they compose all paths that might be traversed through a program during its execution.

<img src="https://i.imgur.com/xjzOCfb.png" width="70%" height="70%" />

#### [Instruction](https://pkg.go.dev/golang.org/x/tools/go/ssa#Instruction)
a statement that consumes values and performs computation. For example, `Call`, `Return`, `TypeAssert`, etc

<img src="https://i.imgur.com/DvheFlc.png" width="50%" height="50%" />

#### [Value](https://pkg.go.dev/golang.org/x/tools/go/ssa#Value)
an expression that yields a value. For example, function calls are both `Instruction` and `Value` since they both consume values and yield a value.

<img src="https://i.imgur.com/oJg97Re.png" width="50%" height="50%" />

And when combined:

<img src="https://i.imgur.com/W02MErA.png" width="70%" height="70%" />

The package contains other [types](https://pkg.go.dev/golang.org/x/tools/go/ssa#pkg-overview) - Include language keywords such as `Defer`, `If` but also lower level primitives like `MakeChan` and `Alloc`. 

### 3.3 Viewing SSA
We can [`ssadump`](https://pkg.go.dev/golang.org/x/tools/cmd/ssadump) to view the SSA form of programs.
```bash
go get -u golang.org/x/tools/cmd/ssadump
ssadump -build=FI ./CompilerMiddleEndSSAInGo/CodeExamples/Channel/
ssadump -build=FI ./CompilerMiddleEndSSAInGo/CodeExamples/ElseIf/
```
We use the `F` to print the SSA code, and `I` to ignore `init` function.
> You can also use this [SSA visualizer](http://golang-ssaview.herokuapp.com/) in view SSA in your CLI. For this example,
> I chose not to, since it it uses a different [build mode](https://pkg.go.dev/golang.org/x/tools/go/ssa#BuilderMode) then 
> the one we need.

Let's consider this program:
``` go
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

I'll focus on the `squareOrCircleArea` function.
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
Also, the assignment to `r` is missing and it's values are already used in the assignment to `area` (`t0`) in the first 
line. This is the result of constant propagation and dead code elimination indicating this code is already optimized.

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

### 3.4 Exercise
In the folder [`CompilerMiddleEndSSAInGo/CodeExamples`](https://github.com/amit-davidson/GopherCon2021IsraelStaticAnalysisWorkshop/tree/master/CompilerMiddleEndSSAInGo/CodeExamples)
there are some interesting programs. Using our SSA visualizer from earlier, take each of the program and look at their SSA.
I added comments with notes with explaining the important points. You should start first with [`CompilerMiddleEndSSAInGo/CodeExamples/Map`](https://github.com/amit-davidson/GopherCon2021IsraelStaticAnalysisWorkshop/blob/master/CompilerMiddleEndSSAInGo/CodeExamples/Map/Map.go)
and then [`CompilerMiddleEndSSAInGo/CodeExamples/ElseIf`](https://github.com/amit-davidson/GopherCon2021IsraelStaticAnalysisWorkshop/blob/master/CompilerMiddleEndSSAInGo/CodeExamples/ElseIf/ElseIf.go) 


### 3.5 SSA vs AST
The most important difference is that AST reasons about the structure of the code, where SSA reasons about how data 
flows in the code. Why do need both? Each "level" suits for a different problem. You can think of it as satellite vs
terrain modes on maps. They both represent the same source map, but each mode solves a different problem. 

We can summarize the differences using the following table:
|                | SSA                                                                                                                                                                                                                                | AST                                                                                                                                                                                             |
|----------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Why to Choose? | <ul><li>Better when need to handle how data flows through the code.</li><li>Does optimizations to the code such as inlining or constant propagation so some functions or variables might be missing</li><li>Package types are closer to the language </li> | <ul><li> Better when when analyzing the code itself, or you don’t want to reason about the control flow graph.</li> <li>Runs over the source code, so optimizations don’t happen yet.</li>|
| Examples       | <ul><li>Checking a function for infinite recursion</li><li> Checking if all flows after “mutex.Lock” are covered with “mutex.unlock”</li>| <ul><li>Passing the correct types to string format</li><li>Shifts that equal or exceed the width of the integer</li><li>Modifying B.n when benchmarking</li><li>Validate the order of imports according to a convention</li>|
 

### 3.6 Writing our analyzer!
In this section we'll implement an analyzer that warns when `t.Fatal` is used inside a goroutine as described here:
https://github.com/ipfs/go-ipfs/issues/2043

### 3.7 Congratulations
You have a good understanding of what IR and SSA are, the SSA package used to create static code analyzers that 
use it and how to write such analyzers.  

In the [next section](https://github.com/amit-davidson/GopherCon2021IsraelStaticAnalysisWorkshop/tree/master/AnalysisApi)
we'll focus on the analysis API. A package used define a common API for all code analyzers and to make writing analyses easier. 
It also provides us an infrastructure that helps us with all the non-logic code such as loading, testing and running our
analysis. 