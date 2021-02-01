package testdata

import "testing"

func main() {
	var t *testing.T
	go func() {
		t.Skip()
	}()
}
