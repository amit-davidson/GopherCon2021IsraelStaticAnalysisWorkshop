package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"testing"
)

func Test_analyzePackage(t *testing.T) {
	var testCases = []struct {
		name   string
		result []string
	}{
		{name: "SimpleRecursionFlag", result: []string{"\"body\" overwrites func parameter in pos"}},
		{name: "AnonymousFunction", result: []string{"\"a\" overwrites func parameter in pos"}},
		{name: "OverwritingParamFromOuterScope", result: []string{"\"a\" overwrites func parameter in pos"}},
		{name: "AssigningParamToAVariableFirst", result: []string{}},
		{name: "MultipleParamsOfSameType", result: []string{"\"a\" overwrites func parameter in pos", "\"c\" overwrites func parameter in pos"}},
		{name: "ShadowingVariable", result: []string{"\"a\" overwrites func parameter in pos"}},
		{name: "EmptyBodyFunction", result: []string{}},
		{name: "NoWarnings", result: []string{}},
		{name: "DecrementOperator", result: []string{"\"retries\" overwrites func parameter"}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(".", "testdata", tc.name, "main.go")
			fset := token.NewFileSet()
			pkgs, err := parser.ParseDir(fset, path, nil, 0)
			require.NoError(t, err)
			require.Len(t, pkgs, 1)
			var testPkg *ast.Package
			for _, pkg := range pkgs {
				testPkg = pkg
			}

			outputs := analyzePackage(testPkg, fset)
			require.NotNil(t, outputs)
			for i := range outputs {
				assert.Contains(t, outputs[i], tc.result[i])
			}
		})
	}
}
