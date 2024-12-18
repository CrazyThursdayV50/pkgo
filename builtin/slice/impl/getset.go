package impl

import (
	"github.com/CrazyThursdayV50/pkgo/builtin"
	"github.com/CrazyThursdayV50/pkgo/builtin/wrap"
)

func (s *Slice[E]) Set(index int, element E) {
	if s == nil {
		return
	}
	if s.Len() < index+1 {
		return
	}
	s.slice[index] = element
}

func (s *Slice[E]) Get(index int) builtin.UnWrapper[E] {
	if s == nil {
		return wrap.Nil[E]()
	}
	if s.Len() < index+1 {
		return wrap.Nil[E]()
	}
	return wrap.Wrap(s.Unwrap()[index])
}
