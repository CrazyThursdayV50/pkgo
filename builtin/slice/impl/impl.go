package impl

import (
	"github.com/CrazyThursdayV50/pkgo/goo"
)

type Slice[E any] struct {
	slice    []E
	lessFunc func(E, E) bool
}

func From[E any](slice ...E) *Slice[E] {
	return &Slice[E]{
		slice: slice,
	}
}

func Make[E any](len, cap int) *Slice[E] {
	return From(make([]E, len, cap)...)
}

func (s *Slice[E]) Cap() int {
	if s == nil {
		return 0
	}
	return cap(s.Unwrap())
}

func (s *Slice[E]) Len() int {
	if s == nil {
		return 0
	}
	return len(s.Unwrap())
}

func (s *Slice[E]) Swap(i, j int) {
	if s == nil {
		return
	}
	s.slice[i], s.slice[j] = s.Unwrap()[j], s.Unwrap()[i]
}

func (s *Slice[E]) Less(i, j int) bool {
	if s == nil {
		return false
	}
	if s.lessFunc == nil {
		return false
	}
	ie := s.Get(i)
	je := s.Get(j)
	if ie == nil || je == nil {
		return false
	}
	return s.lessFunc(ie.Unwrap(), je.Unwrap())
}

func (s *Slice[E]) WithLessFunc(f func(a, b E) bool) {
	if s == nil {
		return
	}
	s.lessFunc = f
}

func (s *Slice[E]) Append(elements ...E) {
	if s == nil {
		return
	}
	s.slice = append(s.Unwrap(), elements...)
}

func (s *Slice[E]) Clear() {
	if s == nil {
		return
	}
	clear(s.Unwrap())
}

func (s *Slice[E]) Chunk(len int) <-chan []E {
	var ch = make(chan []E)
	if s.Len() == 0 || len == 0 {
		goo.Go(func() { close(ch) })
		return ch
	}

	var count = (s.Len()-1)/len + 1
	goo.Go(func() {
		defer close(ch)
		for i := range count {
			from := i * len
			end := from + len
			if end > s.Len() {
				end = s.Len()
			}

			ch <- s.slice[from:end]
		}
	})

	return ch
}
