package gchan

import (
	"github.com/CrazyThursdayV50/pkgo/builtin"
	impl "github.com/CrazyThursdayV50/pkgo/builtin/chan/impl/chanr"
)

var _ builtin.UnWrapper[<-chan any] = (*impl.ChanR[any])(nil)
var _ builtin.Iter[int, any] = (*impl.ChanR[any])(nil)
var _ builtin.ChanAPIR[any] = (*impl.ChanR[any])(nil)

func FromRead[E any](c <-chan E) builtin.ChanAPIR[E] {
	return impl.From(c)
}
