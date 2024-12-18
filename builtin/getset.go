package builtin

type GetSeter[K any, V any] interface {
	Get(K) UnWrapper[V]
	Set(k K, v V)
}
