package testdata

func main(x int) bool {
	if x > 10 {
		goto l1
	}
	return main(x + 1)
l1:
	return true
}
