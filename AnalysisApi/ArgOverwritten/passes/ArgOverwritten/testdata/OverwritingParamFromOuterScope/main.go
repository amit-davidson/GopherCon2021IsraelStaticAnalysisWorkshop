package testdata

func body(a int) {
	_ = func() {
		a = 5 // want `"a" overwrites func parameter`
	}
}

func main() {
	body(5)
}