package impl

func (s *Slice[E]) Unwrap() []E {
	if s == nil {
		return nil
	}
	return s.slice
}
