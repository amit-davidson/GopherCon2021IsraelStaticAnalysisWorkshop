package testdata

func body(a, b int, c int) {
	f := func(a int, b int) {
		a = 5 // want `"a" overwrites func parameter`
	}
	_ = f
}

func main() {
	body(1, 2, 3)
}
