package ArgOverwritten

import (
	"fmt"
	"golang.org/x/tools/go/analysis/analysistest"
	"os"
	"testing"
)

func TestAnalyzer(t *testing.T) {
	var testCases = []struct {
		name     string
	}{
		{name: "SimpleOverwriting"},
		{name: "AnonymousFunction"},
		{name: "OverwritingParamFromOuterScope"},
		{name: "AssigningParamToAVariableFirst"},
		{name: "MultipleParamsOfSameType"},
		{name: "ShadowingVariable"},
		{name: "EmptyBodyFunction"},
		{name: "NoWarnings"},
		{name: "DecrementOperator"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			analysistest.Run(t, fmt.Sprintf("%s%s%s", analysistest.TestData(), string(os.PathSeparator), tc.name), Analyzer)
		})
	}
}
