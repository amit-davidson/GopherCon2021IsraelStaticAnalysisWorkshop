package main

func retry(fn func() error, retries int) error {
	for {
		if err := fn(); err != nil {
			if retries < 1 {
				return err
			}
			retries--
			continue
		}
		return nil
	}
}
func main() {
	_ = retry(func() error {
		return nil
	}, 5)
}
