package builtin

type MapAPI[K comparable, V any] interface {
	Len() int
	Has(k K) bool
	AddSoft(k K, v V)
	Del(k K)
	Keys() SliceAPI[K]
	Values() SliceAPI[V]
	Clear()

	UnWrapper[map[K]V]
	GetSeter[K, V]
	Iter[K, V]
	IterMut[K, V]
}
