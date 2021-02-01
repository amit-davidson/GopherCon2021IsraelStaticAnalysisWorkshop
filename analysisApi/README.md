## 4 Analysis API:
### 4.1 tools/go/analysis
The package defines an API for modular static analysis tools. In other words, it's a common interface for all static 
code analyzers.

Why should you care? What's the difference between using the analysis API, and the way we wrote so far?
The analysis API makes writing analyses easier by taking care of:
 - file parsing
 - testing
 - integration with go vet
 - many more
 
Also, it enforces a single pattern for all the static analysis tools such as how analysis are structured and how 
warnings are reported
    
### 4.2 analysis members  
The primary type in the API is `analysis.Analyzer`.  It describes an analysis function: its name, documentation, flags, relationship to other analyzers, and of course, it's logic.

``` go
type Analyzer struct {
   Name             string
   Doc              string
   Requires         []*Analyzer
   Run              func(*Pass) (interface{}, error)
 
   ...
}
```

The `Name` and `Doc` are obvious. They are used to describe what the tool does.

Another interesting is the `Requires` field. It specifies a list of analyses upon which this one depends and whose
results it may access, and it constrains the order in which a driver may run analyses.

> To use SSA in the analysis api, we would have to require the [SSA builder](https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/buildssa)
>  `Requires: []*analysis.Analyzer {
                  buildssa.Analyzer
          }`


The most important one is the `Run` function. It contains the logic that should is executed upon a single package. 
It takes as an argument `*analysis.Pass` and returns a value to be used by other analyzers and an error.

``` go
type Pass struct {
   Fset         *token.FileSet
   Files        []*ast.File
   Pkg          *types.Package
   TypesInfo    *types.Info
   Report       func(Diagnostic)

   ...
}
```

A `Pass` describes a single unit of work: the application of a particular `Analyzer` to a particular package of Go code. The `Pass` provides information to the Analyzer's `Run` function about the analyzed package and provides operations to the `Run` function for reporting diagnostics and other information back to the driver. It also provides `Fset`, `Files`, `Pkg`, and `TypesInfo` that we know from earlier, so we don't have to take care of ourselves.

The `Report` function emits a diagnostic, a message associated with a source position. For most analyses, diagnostics are their primary result. For convenience, `Pass` provides a helper method, `Reportf`, to report a new diagnostic by formatting a string. Diagnostic is defined as:

``` go
type Diagnostic struct {
   Pos      token.Pos
   Category string // optional
   Message  string
}
```

### 4.3 How to use it
First let's define the project structure:
```
│── README.md
│── cmd
│   └── analyzerName
│       └── main.go
│── go.mod
│── go.sum
└── passes
    └── passName
        │── pass.go
        │── pass_test.go
        └── testdata
```

We create a directory where all of our passes reside in named `passes`. Each pass lives in its package, including its logic and tests.
Then we define the usual `cmd` for our executables that contains all the analyzers the module has.

Regarding our analyzers we wrote previously, we had to handle both the logic of the analyzer and the instrumentation around it.
When converting our code to the analysis API, the AST traversal part (the logic) will sit under `passes` and loading the code
part will be taken care of by the analysis API so we can ignore it. 

Next , we need a way to run the Analyzer and to test it.

### 4.4 Running our code
inside `main.go`, we'll add the following code. 

``` go
package main

import (
   "path/to/our/pass"
   "golang.org/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(passName.Analyzer) }
```
Analyzers are provided in the form of packages that a driver program is expected to import. 
The [`singlechecker`](https://pkg.go.dev/golang.org/x/tools/go/analysis/singlechecker) package provides the `main` function for a command that runs one Analyzer. By convention, each Analyzer should be accompanied by a singlechecker-based command defined in its entirety as: This code calls our Analyzer. 
If we wanted our command to run multiple analyzers, we would have to use [`multichecker`](https://pkg.go.dev/golang.org/x/tools/go/analysis/multichecker).

Now we can run it using 
``` bash
go install path/to/analyzer
go vet -vettool=$(which analyzername) path/to/files
```


### 4.5 Testing our code
The [`analysistest`](https://godoc.org/golang.org/x/tools/go/analysis/analysistes) subpackage provides utilities for testing an Analyzer. Using `analysistest.Run`, it is possible to run an analyzer on a package of `testdata` files and check that it reported all the expected diagnostics.
Expectations are expressed using "// want ..." comments in the input code, such as the following:

``` go
package testdata  
  
func main() {  
   _ = func(a int) {  
      a = 5 // want `"a" overwrites func parameter`  
  }  
}
```

### 4.6 Implementing a code analyzer using the analysis api.   
In this section, we'll convert our `ArgsOverwrite` Analyzer from earlier to the analysis API

### 4.7 Congratulations
You have a good understanding of what the analysis API is and how to use it help us in writing analyses in the future.

In the [next section](https://github.com/amit-davidson/GopherCon2021IsraelStaticAnalysisWorkshop/blob/master/conclusion/text.md)
we'll conclude this workshop by touching a point regarding static code analyzers in general and take a look at other code
analyzers written by the Go community.  