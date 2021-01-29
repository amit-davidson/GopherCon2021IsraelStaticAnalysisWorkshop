## 5. Conclusion and further discussion  
### 5.1 Precision and Soundness
Soundness and precision are two properties used to measure static analysis tools.

A static analysis tool is sound if it finds all of its occurrences in the code for a specific property to check. In other words, the fewer false negatives the static tool has, the more sound it is.

Given the following program:
``` go
r := 1
if rand.Int() > 0.5 {
    r = 2
}
```
An unsound analysis will result in `r=1` ignoring the possible value of `2`, whereas a sound analysis will result in `r=1` or `r=2`.

Precision is the ability of a static analysis tool to flag only the property we're interested in. As previously, the fewer false positives, the more the tool is precise.
If we return the example from above, an imprecise analysis might say `r` is between `0` and `100`, where a precise analysis will say `r` is only `1` or `2`.

Precision and Soundness are a tradeoff. Having the ability to flag more cases makes the program more sound but might result in false positives. On the other way around, limiting the number of cases to cover makes the analysis more precise.

When you'll write static analysis tools, you might encounter where this trade off comes into play, especially on IR level.
Usually, it's easier to implement a more sounder analysis then a more precise one, so in reality, most of the tools go
in this direction. 

## 5.2 Other tools
There are famous built-in tools such as `go vet` and `go fmt`, but there are many more others:

- `go fmt` - Gofmt formats Go programs. It uses tabs for indentation and blanks for alignment. Alignment assumes that an editor is using a fixed-width font.
- `go vet` - Vet examines Go source code and reports suspicious constructs, such as `Printf` calls whose arguments do not align with the format string
- `go imports` - The command go imports updates your Go import lines, adding missing ones and removing unreferenced ones.
- `go fix` - Fix finds Go programs that use old APIs and rewrites them to use newer ones. After you update to a new Go release, fix helps make the necessary changes to your programs.

There's also an awesome (pun intended) [list](https://github.com/golangci/awesome-go-linters) of Go analysis tools written by the Go community.
It contains tools you can integrate into your toolchain. 
   
[staticcheck](https://github.com/dominikh/go-tools) and [golanglint-ci](https://github.com/golangci/golangci-lint) are some of the more noteable tools. 
- staticcheck is similar to `go vet` but applies many more checks such as forgetting to `unlock` a `mutex` using the defer statement, validating JSON tags correctness, and so on.

- golanglint-ci is a fast Go linters runner. It runs linters in parallel, uses caching, supports `yaml` config, has integrations with all major IDE, and has dozens of linters included. You can look at the full list of linters [here](https://golangci-lint.run/usage/linters/).

### 5.3 Further reading:
- A deeper dive into the topics of the AST part in Go - https://github.com/golang/example/tree/master/gotypes  
- A deeper dive into SSA in GO- https://www.youtube.com/watch?v=uTMvKVma5ms&ab_channel=GopherAcademy
- A viewer used to look at the code at different phases of the compilation process - https://golang.design/gossa

### 5.4 Discussion
   Any questions?
