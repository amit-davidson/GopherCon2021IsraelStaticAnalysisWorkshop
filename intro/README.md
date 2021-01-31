## 1.1 program analysis
### 1.1.1 What is program analysis?
   
From Wikipedia:

> Program analysis is the process of automatically analyzing the behavior of computer programs regarding a property such as correctness, robustness, safety and liveness. Program analysis focuses on two major areas: program optimization and program correctness. The first focuses on improving the program's performance while reducing the resource usage while the latter focuses on ensuring that the program does what it is supposed to do.

   
   
### 1.1.2 Benefits of program analysis
Program analysis can help detect bugs before they reach production. You've probably used them already and didn't even 
thought about it. A linter is used to flag styling errors, a profiler to find performance issues and aid in optimizing
them, and even tests to validate program correctness. 

I've gathered a list of issues that could be solved using static analysis tools:

- https://github.com/ipfs/go-ipfs/issues/2043 - Using `t.Fatal` inside a goroutine
- https://github.com/golang/go/issues/23842 - Calling `waitgroup.Add` inside the calling a goroutine
- https://github.com/cockroachdb/cockroach/issues/7972 - Deadlock
- https://github.com/instana/go-sensor/issues/51 - race condition

### 1.1.3 Static analysis vs. Dynamic analysis

> Program analysis can be performed without executing the program (static program analysis), during runtime (dynamic program analysis) or in a combination of both.

The obvious difference between static and dynamic analysis is that dynamic is performed at runtime, whereas static works without running it.

In practice, the difference is much more apparent. The main benefit of dynamic code analysis is that it finds bugs that can actually occur. Also, they are usually easier to write and require fewer resources to run. The problem with dynamic code analysis is that they make the code run slower (I'm looking at you `go race`), which might be intolerant in production. 

On the other hand, static analysis can also find bugs **that may**/**before they** occur. The reason is that static
analysis can evaluate paths of the code or workloads that don't often happen in production. 
Evaluating all the possibilities comes at the expanse of the time and resources required to perform analysis and
inaccuracy of the final result.

## 1.2 Compilers
### 1.2.1 Overview
From Wikipedia:
> In computing, a compiler is a computer program that translates computer code written in one programming language (the source language) into another language (the target language). The name "compiler" is primarily used for programs that translate source code from a high-level programming language to a lower level language (e.g., assembly language, object code, or machine code) to create an executable program.

Compilers are divided into three stages - front end, middle end, and the back end.

- Front end - The front end scans the input and verifies syntax and semantics according to a specific source language.
For example, it makes sure the code contains only familiar characters; it validates that an if statement has a condition
and is not wrapped with parentheses. It also does type checking to make sure the correct types are passed accordingly around the program.
The front end transforms the input program into an intermediate representation (IR) for further processing by the middle end.

- Middle end - The middle end performs optimizations on the intermediate representation to improve the performance and the quality of the produced machine code. Popular optimizations include: [dead code elimination](https://en.wikipedia.org/wiki/Dead_code_elimination), [constant propagation](https://en.wikipedia.org/wiki/Constant_folding) and [inlining](https://en.wikipedia.org/wiki/Inline_expansion). 

- Back end - The back end is responsible for the CPU architecture specific optimizations and code generation - converting IR to machine code.

<img src="https://i.imgur.com/B7ynZWe.png" height="80%" width="80%"/>

### 1.2.2 Frontend
Most commonly today, the frontend is broken into three phases: lexical analysis (also known as lexing or scanning), syntax analysis (also known as scanning or parsing), and semantic analysis.

#### Lexing 
converts a sequence of characters into a sequence of tokens. A token is a pair consisting of a token name and token value.

| Token name | Sample token values              |
|------------|----------------------------------|
| identifier |  x, color, UP                    |
| keyword    |  if, while, return               |
| separator  |  }, (, ;                         |
| operator   |  +, <, =                         |
| literal    |  true, 6.02e23, "music"          |
| comment    |  // can't happen in production   |


#### Parsing 
This involves parsing the token sequence to identify the program. This phase builds a parse tree or an abstract syntax tree, which replaces the linear sequence of tokens with a tree structure.

#### Semantic Analysis 
This phase performs checks such as type checking and rejecting incorrect programs. It also constructs the symbol table
used to map between identifiers and information relating to their declaration or appearance in the source.

For the given code, the type of the variable `n` is inferred in the semantic analysis phase, marking it as `int`
``` go
n := 50 + 70
```



By looking at all the steps, we can see how they come together, and a tree representation of the code is built from source.

<img src="https://i.imgur.com/8sHsDua.png" height="70%" width="70%"/>

### 1.2.3 Middle end
The middle end part performs optimizations regardless of the source code language and the target machine.
As opposed to the front end phase, the middle end analysis is more complex. By estimating how the data will flow through
the code, the compiler does optimizations ranging from the scope of a function to the entire program (interprocedural).

I'll demonstrate an optimization using Constant Propagation mentioned earlier. Constant propagation is the process of substituting the values of known
constants in expressions. Constant propagation eliminates cases in which values are copied from one location or
variable to another to assign their value to another variable.

The algorithm works by doing iterations from top to bottom and and each pass see if it can improve current state. If no
more improvements can be made anymore, the algorithm terminates. 

Propagating X:
```
  int x = 14; <-
  int y = 7 - x / 2;
  return y * (28 / x + 2);
```
```
  int x = 14;
  int y = 7 - 14 / 2; <-
  return y * (28 / x + 2);
```

```
  int x = 14;
  int y = 7 - 14 / 2;
  return y * (28 / 14 + 2); <-
```
Propagating Y:
```
  int x = 14; <-
  int y = 7 - 14 / 2;
  return y * (28 / 14 + 2); 
```
```
  int x = 14;
  int y = 0; <-
  return y * (28 / 14 + 2); 
```

```
  int x = 14;
  int y = 0;
  return 0; <-
```

We can further optimize this code using dead code elimination - The process of removing code that does not affect 
the program results. We can optimize away `x` and `y' which results in 
```
return 0;
```  

### 1.3 Congratulations
You now have a grasp about program analyzers and a good understanding of how compilers are built and work.

In the [next section](https://github.com/amit-davidson/GopherCon2021IsraelStaticAnalysisWorkshop/blob/master/ast/text.md) 
we'll focus on the front end level, and see how analyzer "operating" in this level work 
