package impl

import (
	"sync/atomic"

	"github.com/CrazyThursdayV50/pkgo/builtin"
	"github.com/CrazyThursdayV50/pkgo/builtin/wrap"
)

func (c *ChanRW[E]) Iter(f func(index int, element E) (bool, error)) (builtin.UnWrapper[int], error) {
	if c == nil {
		return wrap.Nil[int](), nil
	}

	for e := range c.c {
		atomic.AddInt64(&c.countR, 1)
		ok, err := f(int(c.countR), e)
		if !ok {
			return wrap.Wrap(int(c.countR)), err
		}
	}

	return wrap.Nil[int](), nil
}
