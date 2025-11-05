package impl

func (m *Map[K, V]) Unwrap() map[K]V {
	if m == nil {
		return nil
	}
	return m.m
}

func (m *Map[K, V]) IsNil() bool {
	return m == nil || m.m == nil
}
