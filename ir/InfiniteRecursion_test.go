package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/analysis/analysistest"
	"path/filepath"
	"testing"
)

func Test_analyzePackage(t *testing.T) {
	var testCases = []struct {
		name   string
		result string
	}{
		{name: "ConditionalReturnWithGoToNoFlag", result: ""},
		{name: "InfiniteRecursionWithGoroutineNoFlag", result: ""},
		{name: "SameStructMethodFlag", result: errMessage},
		{name: "SameTypeDifferentStructMethodFlag", result: errMessage},
		{name: "SameTypeDifferentStructMethodNoFlag", result: ""},
		{name: "SimpleRecursionFlag", result: errMessage},
		{name: "SimpleRecursionWithReturnInCallFlag", result: errMessage},
		{name: "UnreachableReturnFlag", result: errMessage},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(analysistest.TestData(), tc.name, "main.go")
			prog, err := loadPackage(path)
			require.NoError(t, err)
			funcs := getAllFunctions(prog)
			for _, fn := range funcs {
				res := checkInfiniteRecursion(fn)
				assert.Contains(t, res, tc.result)
			}
		})
	}
}
