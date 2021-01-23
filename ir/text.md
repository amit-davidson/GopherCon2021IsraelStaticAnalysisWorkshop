2. IR
	2.1 What is IR?
	An **intermediate representation** (**IR**) is the code used internally by a compiler to represent source code. An IR is designed to be conducive for further processing, such as optimization and translation. A "good" IR must be _accurate_ – capable of representing the source code without loss of information  – and _independent_ of any particular source or target language.
	
	2.2 What is SSA?
		SSA stands for static single assignment. It's a property of an IR **that requires each variable be assigned exactly once**, and every variable be defined before it is used. 
		The primary usefulness of SSA comes from how it simplifies the properties of variables and improves compilers optimizations.

For example, consider this piece of code:
```
y := 1
y := 2
x := y
```

In an SSA form it'll be translated to:
```
y1 := 1
y2 := 2
x1 := y2
```
Humans can see that the first assignment is not necessary, and that the value of `y` being used in the third line comes from the second assignment of `y`. In SSA form, both of these are immediate


2.3 tools/go/SSA package members
Program - A Program is a partial or complete Go program converted to SSA form.
![Program](https://i.imgur.com/FHbYxeU.png =450x400)

Package - A Package is a single analyzed Go package containing Members for all package-level functions, variables, constants and types it declares.
![Package](https://i.imgur.com/eLzMEHR.png =450x400)

Function - Function represents the parameters, results, and code of a function or method.
![Function](https://i.imgur.com/FqN1GdN.png =600x400)

Basic Block - BasicBlock represents an SSA basic block. A set of instructions that are executed and can't jump somewhere else. Basic blocks are connected using conditions and goto statements. 
![Function](https://i.imgur.com/XGrpRkH.png =600x400)
Control Flow Graph (CFG) - In a control-flow graph each node in the graph represents a basic block. Together, they all paths that might be traversed through a program during its execution.
![Function](https://i.imgur.com/jpmXl4P.png =700x200)

Instruction - a statement that consumes values and performs computation. For example, `Call`, `Return`, `TypeAssert`, etc
![Function](https://i.imgur.com/VJ5mxF3.png =600x400)
Value - an expression that yields a value. Function calls for example are both `Instruction` and `Value` since they both consume values but also yield a value.

![Function](https://i.imgur.com/UlKSNVu.png =600x400)
[Other functions types](https://pkg.go.dev/golang.org/x/tools/go/ssa#pkg-overview) - Include language keywords such as `Defer`, `If` but also lower level primivites like `MakeChan` and `Alloc`. 

2.4 Viewing SSA
We can use this  [SSA visualizer](http://golang-ssaview.herokuapp.com/)  to view the SSA form of programs.

> You can also use `go.tools/cmd/ssadump` in view SSA in your CLI

Let's consider this program:
```
package main  
  
import (  
   "fmt"  
   "math/rand"
 )  
  
func main() {  
   a := 4  
  if rand.Int() > 0 {  
      fmt.Println(a)  
   }  
}
```

It's SSA representation looks like the following. At the top, we can see the package members, the `main` function, and the implicit `init` function.
After that we can see the `init` function it's [attributes](https://pkg.go.dev/golang.org/x/tools/go/ssa#Function), and it's CFG. We'll ignore that for now and focus on the `main` function CFG:
```
package main.go:
  func  init       func()
  var   init$guard bool
  func  main       func()

# Name: main.go.init
# Package: main.go
# Synthetic: package initializer
func init():
0:                                                                entry P:0 S:2
	t0 = *init$guard                                                   bool
	if t0 goto 2 else 1
1:                                                           init.start P:1 S:1
	*init$guard = true:bool
	t1 = math/rand.init()                                                ()
	t2 = fmt.init()                                                      ()
	jump 2
2:                                                            init.done P:2 S:0
	return

# Name: main.go.main
# Package: main.go
# Location: main.go:8:6
# Locals:
#   0:	t0 int
func main():
0:                                                                entry P:0 S:2
    # `a` is defined
	t0 = local int (a)                                                 *int
	# `a` is assigned with 4
	*t0 = 4:int
	# A variable used to hold the `rand` result is defined and assigned - doesn't appear in the source code
	t1 = math/rand.Int()                                                int
	t2 = t1 > 0:int                                                    bool
	# According to conditions results, gotos are used to determine the flow of the function
	if t2 goto 1 else 2
1:                                                              if.then P:1 S:1
	t3 = *t0                                                            int
	# an interface used to hold varargs is defined
	t4 = new [1]interface{} (varargs)                       *[1]interface{}
	t5 = &t4[0:int]                                            *interface{}
	# `a` is converted to an interface so it can be passed to the function 
	t6 = make interface{} <- int (t3)                           interface{}
	*t5 = t6
	t7 = slice t4[:]                                          []interface{}
	t8 = fmt.Println(t7...)                              (n int, err error)
	jump 2
2:                                                              if.done P:2 S:0
	t9 = *t0                                                            int
	# Instruction inserted by the compiler
	rundefers
	return
```
As you can see some instructions that we don't expect to see at the soruce code level were inserted by the compiler in the SSA format. Like the example of the if condition, they are a more verbose form of the same instruction.

2.5 SSA vs AST
	AST, shows us the structure of the code. How different statements in the code relate to each other. SSA on the the hand shows us how the code flows.

When applying this logic to static analysis we'll see that SSA will be used for more complex analysis where we need to determine the flow of the data, where as AST, will be used for simpler, more structure related, analyses.

2.6 Writing our analyzer!
https://github.com/ipfs/go-ipfs/issues/2043