package impl

func (s *Slice[E]) Unwrap() []E {
	if s == nil {
		return nil
	}
	return s.slice
}

func (s *Slice[E]) IsNil() bool {
	return s == nil || s.slice == nil
}
