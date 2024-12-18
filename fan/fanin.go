package fan

import (
	"github.com/CrazyThursdayV50/pkgo/fan/in"
	"github.com/CrazyThursdayV50/pkgo/fan/in/impl"
)

func NewIn[element any](handler func(element), from ...<-chan element) in.Fan {
	return impl.New(handler, from...)
}
