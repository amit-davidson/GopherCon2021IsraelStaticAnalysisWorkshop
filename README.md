# Static Analysis with Go - A Practitioner's Guide
[![made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](http://golang.org)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat-square)](http://makeapullrequest.com)

Hi, and welcome to Static Analysis with Go - A Practitioner's Guide. 
This is a workshop about writing static code analyzers in Go for Go.
In this lecture you will learn about how to write a static code analysis in Go, and implement one yourself. 

By the end of this workshop, you'll have a better understanding of the Go packages related to writing static code
analyzers and you'll also know how to write a code analyzer yourself. 

I will start the lecture by giving an overview of static analyzers and how compilers work. Afterwards, you will learn
about 2 different representations of the code (AST and SSA), and write an analyzer in each of those. 
You will also learn about the analysis API making writing analyzers easier and then we'll finish with a discussion.

### Requirements:
Install the repo:
```bash
git clone https://github.com/amit-davidson/GopherCon2021IsraelStaticAnalysisWorkshop.git
```

Install [Go tools](https://github.com/golang/tools):
```bash
go get -u golang.org/x/tools/...
```

Validate the installation by running
```bash
ssadump -h
```
and making sure you get a help message that starts with: `Usage of ssadump:`

### Contents:
1. [Introduction to compilers and program analysis](https://github.com/amit-davidson/GopherCon2021IsraelStaticAnalysisWorkshop/tree/master/Intro)
2. [Compiler front end and static analysis with AST In Go](https://github.com/amit-davidson/GopherCon2021IsraelStaticAnalysisWorkshop/blob/master/CompilerFrontEndASTInGo)
3. [Compiler middle end and static analysis with SSA In Go](https://github.com/amit-davidson/GopherCon2021IsraelStaticAnalysisWorkshop/blob/master/CompilerMiddleEndSSAInGo)
3. [The analysis API](https://github.com/amit-davidson/GopherCon2021IsraelStaticAnalysisWorkshop/tree/master/AnalysisApi)
3. [Conclusion](https://github.com/amit-davidson/GopherCon2021IsraelStaticAnalysisWorkshop/tree/master/Conclusion)