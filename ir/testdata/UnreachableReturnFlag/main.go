package testdata

func main(p *int, n int) {
	x := 0
	main(&x, n-1)
	if x != n {
		panic("stack is corrupted")
	}

}
