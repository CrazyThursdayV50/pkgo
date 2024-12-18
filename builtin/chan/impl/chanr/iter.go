package impl

import (
	"sync/atomic"

	"github.com/CrazyThursdayV50/pkgo/builtin"
	"github.com/CrazyThursdayV50/pkgo/builtin/wrap"
)

func (c *ChanR[E]) Iter(f func(index int, element E) (bool, error)) (builtin.UnWrapper[int], error) {
	if c == nil {
		return wrap.Nil[int](), nil
	}

	for e := range c.c {
		atomic.AddInt64(&c.count, 1)
		ok, err := f(int(c.count), e)
		if !ok {
			return wrap.Wrap(int(c.count)), err
		}
	}

	return wrap.Nil[int](), nil
}
