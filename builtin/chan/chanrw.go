package gchan

import (
	"github.com/CrazyThursdayV50/pkgo/builtin"
	impl "github.com/CrazyThursdayV50/pkgo/builtin/chan/impl/chanrw"
)

var _ builtin.UnWrapper[chan any] = (*impl.ChanRW[any])(nil)
var _ builtin.Iter[int, any] = (*impl.ChanRW[any])(nil)
var _ builtin.ChanAPI[any] = (*impl.ChanRW[any])(nil)

func Make[E any](buff int) builtin.ChanAPI[E] {
	return impl.From(make(chan E, buff))
}

func From[E any](c chan E) builtin.ChanAPI[E] {
	return impl.From(c)
}
