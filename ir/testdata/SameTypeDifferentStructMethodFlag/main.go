package testdata

type T struct {
	n int
}

func (t T) Fn2() {
	x := T{}
	x.Fn2()
}
