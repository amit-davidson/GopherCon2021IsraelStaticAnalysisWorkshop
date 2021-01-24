## 1 Analysis API:
### 1.1 tools/go/analysis
The package defines an API for modular static analysis tools. In other words, it's a common interface for all static code analyzers.
      
The primary type in the API is `analysis.Analyzer`.  It describes an analysis function: its name, documentation, flags, relationship to other analyzers, and of course, it's logic.

```
type Analyzer struct {
   Name             string
   Doc              string
   Flags            flag.FlagSet
   Run              func(*Pass) (interface{}, error)
   RunDespiteErrors bool
   ResultType       reflect.Type
   Requires         []*Analyzer
   FactTypes        []Fact
}
```

The `Name` and `Doc` are obvious. They are used to describe what the tool does.

Another interesting field is the `Run` function. It contains the logic that should is executed upon a single package. It takes as an argument `*analysis.Pass` and returns a value to be used by other analyzers and an error.

```
type Pass struct {
   Fset         *token.FileSet
   Files        []*ast.File
   OtherFiles   []string
   IgnoredFiles []string
   Pkg          *types.Package
   TypesInfo    *types.Info
   ResultOf     map[*Analyzer]interface{}
   Report       func(Diagnostic)
   ...
}
```

A `Pass` describes a single unit of work: the application of a particular `Analyzer` to a particular package of Go code. The `Pass` provides information to the Analyzer's `Run` function about the analyzed package and provides operations to the `Run` function for reporting diagnostics and other information back to the driver. It also provides `Fset`, `Files`, `Pkg`, and `TypesInfo` that we know from earlier, so we don't have to take care of ourselves.

The `Report` function emits a diagnostic, a message associated with a source position. For most analyses, diagnostics are their primary result. For convenience, `Pass` provides a helper method, `Reportf`, to report a new diagnostic by formatting a string. Diagnostic is defined as:

```
type Diagnostic struct {
   Pos      token.Pos
   Category string // optional
   Message  string
}
```

### 1.2 How to use it
First let's define the project structure:
<pre>
|-- README.md
|-- cmd
|   `-- analyzerName
|       `-- main.go
|-- go.mod
|-- go.sum
`-- passes
    `-- passName
        |-- pass.go
        |-- pass_test.go
        `-- testdata
</pre>

We create a directory where all of our passes reside in named `passes`. Each pass lives in its package, including its logic and tests. Then we define the usual `cmd` for our executables that contains all the analyzers the module has.

So far, our code sat under `passes`  where each Analyzer had its own pass folder. Now, we need a way to run the Analyzer and to test it.

### 1.3 Running our code
inside `main.go`, we'll add the following code. 

```
package main

import (
   "path/to/our/pass"
   "golang.org/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(passName.Analyzer) }
```
Analyzers are provided in the form of packages that a driver program is expected to import. 
The `singlechecker` package provides the `main` function for a command that runs one Analyzer. By convention, each Analyzer should be accompanied by a singlechecker-based command defined in its entirety as: This code calls our Analyzer. 
If we wanted our command to run multiple analyzers, we would have to use `tools/go/analysis/multichecker`.

### 1.4 Testing our code
The `analysistest` subpackage provides utilities for testing an Analyzer. Using `analysistest.Run`, it is possible to run an analyzer on a package of `testdata` files and check that it reported all the expected diagnostics.
Expectations are expressed using "// want ..." comments in the input code, such as the following:

```
package testdata  
  
func main() {  
   _ = func(a int) {  
      a = 5 // want `"a" overwrites func parameter`  
  }  
}
```

### 1.5 Implementing a code analyzer using the analysis api.   
In this section, we'll convert our `ArgsOverwrite` Analyzer from earlier to the analysis API

### 1.6 Integrating it as part of our toolchain. 
We can run our analysis in 2 ways:
1. Run it directly
2. Using `go vet` with the following command: 
```
go vet -vettool=$(which analyzer name) path/to/files
```
