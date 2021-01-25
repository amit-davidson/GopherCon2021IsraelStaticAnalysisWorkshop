package main

type IceCreamMaker interface {
	// Hello greets a customer
	Hello()
}

type Jerry struct {
	name string
}

func (j *Jerry) Hello() {
	j.name = "Jerry"
}

func sayHello() {
	var maker IceCreamMaker
	maker = <-channel
	maker.Hello()
}

var channel = make(chan IceCreamMaker)

func main() {
	jerry := &Jerry{}
	channel <- jerry
	go sayHello()
	jerry.name = "1"
}

// 1. Focus on how values are sent to and received from the channel
// 2. In sayHello see how t1.Hello is preceded with "invoke" indicating it's a an interface method (The channel
// receives an interface type)