package wrap

type unwrapper[T any] struct {
	t T
}

func Nil[T any]() *unwrapper[T] {
	var w *unwrapper[T]
	return w
}

func Wrap[T any](t T) *unwrapper[T] {
	return &unwrapper[T]{t: t}
}

func (w *unwrapper[T]) Unwrap() T {
	if w == nil {
		var zero T
		return zero
	}
	return w.t
}
