package main

func main() {
	for n := 0; n <= 5; n++ {
		if n%2 == 0 {
			continue
		}
	}
}

//1. See the properties of the for loop

//2. See the properties of the if condition

//3. See how continue is represented in the code

// You can go to https://golang.org/pkg/go/ast/ for the types and their documentation. It'll give you an explanation
// about the properties of each type.
