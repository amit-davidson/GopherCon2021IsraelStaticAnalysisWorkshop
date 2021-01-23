1. Intro
	This is a workshop about writing static code analyzers in Go for Go. I will start the lecture by giving an overview of static analyzers and how compilers work. 

	The reality is that static analysis is closely related to compilers. They perform static analysis over the code to validate it's correctness and optimize it, and we'll also need to use the compilers capabalities to parse the code so we can analyze it.
	
1.1.1 What is program analysis?
	
From wikipedia:

> Program analysis is the process of automatically analyzing the behavior of computer programs regarding a property such as correctness, robustness, safety and liveness. Program analysis focuses on two major areas: program optimization and program correctness. The first focuses on improving the program’s performance while reducing the resource usage while the latter focuses on ensuring that the program does what it is supposed to do.

	
	
1.1.2 Benefits of program analysis
		Program analysis can help detect bugs before they reach production. You've probably have already used them and didn't even thought about. A linter is used to flag styling errors, profiler to find performance issues and aid in optimizing them, and even tests to validate program correctness. 
		
1.1.3 Static analyses vs Dynamic analyses

> Program analysis can be performed without executing the program (static program analysis), during runtime (dynamic program analysis) or in a combination of both.

As stated, the main difference between static and dynamic analyses, is that dynamic is performed at runtime where as static works without running it.

In practice, the difference is much more obvious. The main benefit of dynamic code analysis is that it finds bugs that can actually occur. Also, they are usually easier to write. The problem with dynamic code analyses is that they make your coder slower (I'm looking at you `go race`). 

Static analyses on the other hand, can also find bugs **that may**/**before they** occur. The reason is that static analyses can evaluate paths of your code or worloads that don't happen in production. Evaluating all this data comes at the expanse of the time and resources required to perform analyses and also inaccurracy of the final result.

1.2 Compilers
			1.2.1 Overview
		 From wikipedia:
		> In computing, a compiler is a computer program that translates computer code written in one programming language (the source language) into another language (the target language). The name "compiler" is primarily used for programs that translate source code from a high-level programming language to a lower level language (e.g., assembly language, object code, or machine code) to create an executable program.

Compilers are divided into 3 stages - front end, middle end and the back end.
![Compiler overview](https://i.imgur.com/9x70LAl.png)

Front end - The front end scans the input and verifies syntax and semantics according to a specific source language. For example, it makes sure the code contains only familiar charcthers, it validates that an if statement contains a condition and it's not wrapped with parentheses and it also does type checking.
The front end transforms the input program into an intermediate representation (IR) for further processing by the middle end

Middle end - The middle end performs optimizations on the intermediate representation in order to improve the performance and the quality of the produced machine code. Popular optimzations include: [dead code elimination](https://en.wikipedia.org/wiki/Dead_code_elimination), [constant propagation](https://en.wikipedia.org/wiki/Constant_folding) and [inlining](https://en.wikipedia.org/wiki/Inline_expansion). 

Back end - The back end is responsible for the CPU architecture specific optimizations and for code generation - converting IR to machine code.

1.2.2 Frontend
	Most commonly today, the frontend is broken into three phases: lexical analysis (also known as lexing or scanning), syntax analysis (also known as scanning or parsing), and semantic analysis.

![Compiler frontend overview](https://i.imgur.com/muZGoQt.png)

Lexing - converts a sequence of characters into a sequence of tokens. A token is a pair consisting of a token name and token value. Common token names include keyword, separator, identifier, literal	and some of their respectively `while`, `{`, `x`, `"music"`.

Parsing - involves parsing the token sequence to identify the of the program. This phase builds a parse tree or an abstract syntax tree, which replaces the linear sequence of tokens with a tree structure.

Semantic Analysis - This phase performs checks such as type checking and rejecting incorrect programs. It also constructs the symbol table used to map between identifiers and information relating to their declaration or appearance in the source.

![Compiler frontend](https://i.imgur.com/biUHNJq.png)

1.2.3 Middle end
		As explained, the middle end performs optimzations regardless of the source code language and the target machine. As opposed to the front end phase, the middle end analyses are more complex. By estimating how the code and the data will flow, the compiler does optimizations ranging from the scope of a function to  the entire program (interprocedural). 

For example, using constant propagation optimization we will get the following: 
```
  int x = 14;
  int y = 7 - x / 2;
  return y * (28 / x + 2);
```
->
```
  int x = 14;
  int y = 7 - 14 / 2;
  return y * (28 / 14 + 2);
```
->
```
  int x = 14;
  int y = 0;
  return 0;
```
(And using dead code elimination we can optimize away `x` and `y`)