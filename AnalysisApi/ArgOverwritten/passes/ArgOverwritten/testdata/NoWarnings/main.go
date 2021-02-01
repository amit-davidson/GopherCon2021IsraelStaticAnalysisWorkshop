package testdata

import "fmt"

func closeBody(body int) {
	fmt.Print(body)
}

func main() {
	closeBody(1)
}
