package testdata

func body(a int) {
	_ = func() {
		a = 5
	}
}

func main() {
	body(5)
}