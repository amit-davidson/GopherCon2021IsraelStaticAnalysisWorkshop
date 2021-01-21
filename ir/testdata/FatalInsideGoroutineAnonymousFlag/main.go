package testdata

import "testing"

func main() {
	var t *testing.T
	fn := func() {
		t.Fatal()
	}
	go fn()
}
