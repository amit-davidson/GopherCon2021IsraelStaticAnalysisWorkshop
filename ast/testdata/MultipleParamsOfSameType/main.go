package testdata

func body(a, b int, c int) {
	a = 5
	_ = b
	c = 3
}

func main() {
	body(1, 2, 3)
}
