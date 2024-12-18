package slice

import (
	"github.com/CrazyThursdayV50/pkgo/builtin"
	"github.com/CrazyThursdayV50/pkgo/builtin/slice/impl"
)

var _ builtin.GetSeter[int, any] = (*impl.Slice[any])(nil)
var _ builtin.UnWrapper[[]any] = (*impl.Slice[any])(nil)
var _ builtin.Iter[int, any] = (*impl.Slice[any])(nil)
var _ builtin.IterMut[int, any] = (*impl.Slice[any])(nil)
var _ builtin.SliceAPI[any] = (*impl.Slice[any])(nil)

func From[E any](sli ...E) builtin.SliceAPI[E] {
	return impl.From(sli...)
}

func Make[E any](len, cap int) builtin.SliceAPI[E] {
	return impl.Make[E](len, cap)
}

func Empty(len int) builtin.SliceAPI[struct{}] {
	return impl.Make[struct{}](len, len)
}
