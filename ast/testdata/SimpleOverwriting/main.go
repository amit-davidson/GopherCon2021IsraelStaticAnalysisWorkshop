package testdata

func closeBody(body int) {
	body = 1
}

func main() {
	closeBody(1)
}
