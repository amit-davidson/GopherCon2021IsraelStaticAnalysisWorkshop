package testdata

func closeBody(body int) {
	body = 1 // want `"body" overwrites func parameter`
}

func main() {
	closeBody(1)
}
