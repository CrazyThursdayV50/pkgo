package impl

func (m *Map[K, V]) Unwrap() map[K]V {
	if m == nil {
		return nil
	}
	return m.m
}
