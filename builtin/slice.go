package builtin

type SliceAPI[E any] interface {
	Cap() int
	Len() int
	Swap(i int, j int)
	Less(i int, j int) bool
	WithLessFunc(f func(a, b E) bool)
	Append(elements ...E)
	Clear()
	Chunk(len int) <-chan []E

	UnWrapper[[]E]
	GetSeter[int, E]
	Iter[int, E]
	IterMut[int, E]
}
