package builtin

type Iter[K comparable, V any] interface {
	Iter(func(k K, v V) (bool, error)) (UnWrapper[K], error)
}

type IterMut[K comparable, V any] interface {
	IterMut(func(k K, v V, s GetSeter[K, V]) (bool, error)) (UnWrapper[K], error)
}
