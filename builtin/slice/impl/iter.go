package impl

import (
	"github.com/CrazyThursdayV50/pkgo/builtin"
	"github.com/CrazyThursdayV50/pkgo/builtin/wrap"
)

func (s *Slice[E]) Iter(f func(index int, element E) (bool, error)) (builtin.UnWrapper[int], error) {
	if s == nil {
		return wrap.Wrap(-1), nil
	}
	for i, e := range s.Unwrap() {
		ok, err := f(i, e)
		if !ok {
			return wrap.Wrap(i), err
		}
	}
	return wrap.Wrap(s.Len()), nil
}

func (s *Slice[E]) IterMut(f func(index int, element E, self builtin.GetSeter[int, E]) (bool, error)) (builtin.UnWrapper[int], error) {
	if s == nil {
		return wrap.Wrap(-1), nil
	}
	for i, e := range s.Unwrap() {
		ok, err := f(i, e, s)
		if !ok {
			return wrap.Wrap(i), err
		}
	}
	return wrap.Wrap(s.Len()), nil
}
