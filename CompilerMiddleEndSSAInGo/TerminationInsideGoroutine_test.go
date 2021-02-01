package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"path/filepath"
	"testing"
)

func Test_analyzePackage(t *testing.T) {
	var testCases = []struct {
		name string

		result string
	}{
		{name: "FatalInsideGoroutineSimpleFlag", result: errMessage},
		{name: "SkipInsideGoroutineSimpleFlag", result: errMessage},
		{name: "FatalInsideGoroutineAnonymousFlag", result: errMessage},
		{name: "FatalInsideRegularFunctionNoFlag", result: ""},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(".", "testdata", tc.name, "main.go")
			prog, pkg, err := loadPackage(path)
			require.NoError(t, err)
			funcs := getAllFunctions(pkg)
			for _, fn := range funcs {
				res := checkTerminationInsideGoroutine(fn, prog.Fset)
				if fn.Name() == "main" {
					assert.Contains(t, res, tc.result)
				}
			}
		})
	}
}
