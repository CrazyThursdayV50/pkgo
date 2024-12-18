package impl

import (
	"github.com/CrazyThursdayV50/pkgo/builtin"
	"github.com/CrazyThursdayV50/pkgo/builtin/wrap"
)

func (m *Map[K, V]) Iter(f func(k K, v V) (bool, error)) (builtin.UnWrapper[K], error) {
	if m == nil {
		return wrap.Nil[K](), nil
	}
	m.l.RLock()
	defer m.l.RUnlock()
	for k, v := range m.m {
		ok, err := f(k, v)
		if !ok {
			return wrap.Wrap(k), err
		}
	}
	return wrap.Nil[K](), nil
}

func (m *Map[K, V]) IterMut(f func(k K, v V, self builtin.GetSeter[K, V]) (bool, error)) (builtin.UnWrapper[K], error) {
	if m == nil {
		return wrap.Nil[K](), nil
	}

	keys := m.Keys()
	index, err := keys.Iter(func(_ int, element K) (bool, error) {
		return f(element, m.Get(element).Unwrap(), m)
	})

	if index == nil {
		return wrap.Nil[K](), err
	}

	return keys.Get(index.Unwrap()), err
}
