package implements

import (
	"github.com/CrazyThursdayV50/pkgo/broadcast/iface"
	"github.com/CrazyThursdayV50/pkgo/broadcast/implements/mem"
)

func NewMemoryBroadcaster[T any]() iface.Broadcaster[T] {
	return mem.New[T]()
}