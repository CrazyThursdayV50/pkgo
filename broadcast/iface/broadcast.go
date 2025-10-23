package iface

type Broadcaster[T any] interface {
	Subscribe() chan T
	Unsubscribe(ch chan T)
	Broadcast(msg T)
	Close()
}
