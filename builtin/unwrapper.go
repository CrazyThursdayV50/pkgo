package builtin

type UnWrapper[T any] interface {
	Unwrap() T
}
