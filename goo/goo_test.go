package goo

import "testing"

func TestTry(t *testing.T) {
	err := TryE(func() {
		_ = []string{""}[1]
		panic("?")
	})
	t.Errorf("err: %v", err)

	err = TryE(func() {
		panic("?")
	})

	t.Errorf("err: %v", err)

	err = TryE(func() {
		a := 1
		b := 0
		_ = a / b
	})

	t.Errorf("err: %v", err)

	err = TryE(func() {
		var fn func()
		fn()
	})

	t.Errorf("err: %v", err)
}
