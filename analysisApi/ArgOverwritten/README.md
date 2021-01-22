# ArgOverwritten

[![made-with-Go](https://github.com/go-critic/go-critic/workflows/Go/badge.svg)](http://golang.org)
[![made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](http://golang.org)
[![MIT license](https://img.shields.io/badge/License-MIT-blue.svg)](https://lbesson.mit-license.org/)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat-square)](http://makeapullrequest.com)

ArgOverwritten finds function arguments being overwritten

## Quick Start:

Download the package

```
go get github.com/amit-davidson/ArgOverwritten/cmd/argoverwritten
```

Pass the entry point

```
go vet -vettool=$(which argoverwritten) ${path_to_file}
```

## Example:
```
package testdata

func body(a int) {
	_ = func() {
		a = 5
	}
}

func main() {
	body(5)
}
```
```
go vet -vettool=$(which argoverwritten) /testdata/OverwritingParamFromOuterScope
```

Output
```
/testdata/OverwritingParamFromOuterScope/main.go:5:3: "a" overwrites func parameter
```