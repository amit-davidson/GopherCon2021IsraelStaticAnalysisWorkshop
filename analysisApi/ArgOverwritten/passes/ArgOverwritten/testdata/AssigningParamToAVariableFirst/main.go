package testdata

func body(a int) {
	b := a
	b = 5
	_ = b
}

func main() {
	body(5)
}