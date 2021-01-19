package main

// Add returns the sum of a and b.
func Add(a int64, b int64) int64

func body(a int) {
	b := a
	b = 5
	_ = b
}

func main() {
	body(5)
	_ = Add(3, 4)
}
