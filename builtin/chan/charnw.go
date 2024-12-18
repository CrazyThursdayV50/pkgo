package gchan

import (
	"github.com/CrazyThursdayV50/pkgo/builtin"
	impl "github.com/CrazyThursdayV50/pkgo/builtin/chan/impl/charw"
)

var _ builtin.UnWrapper[chan<- any] = (*impl.ChanW[any])(nil)
var _ builtin.ChanAPIW[any] = (*impl.ChanW[any])(nil)

func FromWrite[E any](c chan<- E) builtin.ChanAPIW[E] {
	return impl.FromChanW(c)
}
