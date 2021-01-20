package testdata

type T struct {
	n int
}

func (t T) Fn3() {
	if t.n == 0 {
		return
	}
	t.Fn3()
}
