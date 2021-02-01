package testdata

func main() {
	_ = func(a int) {
		a = 5 // want `"a" overwrites func parameter`
	}
}