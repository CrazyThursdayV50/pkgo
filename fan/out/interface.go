package out

type Fan[T any] interface {
	Close()
	Send(element T)
}
