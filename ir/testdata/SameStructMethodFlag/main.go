package testdata

type T struct {
	n int
}

func (t T) Fn1() {
	t.Fn1()
}