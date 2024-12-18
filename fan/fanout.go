package fan

import (
	"github.com/CrazyThursdayV50/pkgo/fan/out"
	"github.com/CrazyThursdayV50/pkgo/fan/out/impl"
)

var _ out.Fan[any] = (*impl.Fan[any])(nil)

func NewOut[element any](count, buffer int, handler func(element)) out.Fan[element] {
	return impl.New[element](count, buffer, handler)
}
