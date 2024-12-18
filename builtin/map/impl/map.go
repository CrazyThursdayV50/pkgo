package impl

import (
	"sync"

	"github.com/CrazyThursdayV50/pkgo/builtin"
	"github.com/CrazyThursdayV50/pkgo/builtin/slice/impl"
)

type Map[K comparable, V any] struct {
	l *sync.RWMutex
	m map[K]V
}

func Make[K comparable, V any](cap int) *Map[K, V] {
	return From[K, V](make(map[K]V, cap))
}

func From[K comparable, V any](m map[K]V) *Map[K, V] {
	if m == nil {
		return Make[K, V](0)
	}

	return &Map[K, V]{
		l: &sync.RWMutex{},
		m: m,
	}
}

func (m *Map[K, V]) Len() int {
	if m == nil {
		return 0
	}
	return len(m.m)
}
func (m *Map[K, V]) Has(k K) bool {
	if m == nil {
		return false
	}
	m.l.RLock()
	defer m.l.RUnlock()
	_, ok := m.m[k]
	return ok
}

// add if not exist
func (m *Map[K, V]) AddSoft(k K, v V) {
	if m == nil {
		return
	}

	if m.Has(k) {
		return
	}

	m.Set(k, v)
}
func (m *Map[K, V]) Del(k K) {
	if m == nil {
		return
	}
	m.l.Lock()
	defer m.l.Unlock()
	delete(m.m, k)
}

func (m *Map[K, V]) Keys() builtin.SliceAPI[K] {
	if m == nil {
		return nil
	}
	m.l.RLock()
	defer m.l.RUnlock()
	slice := impl.Make[K](0, m.Len())
	for k := range m.m {
		slice.Append(k)
	}

	return slice
}

func (m *Map[K, V]) Values() builtin.SliceAPI[V] {
	if m == nil {
		return nil
	}
	m.l.RLock()
	defer m.l.RUnlock()
	slice := impl.Make[V](0, m.Len())
	for _, v := range m.m {
		slice.Append(v)
	}

	return slice
}

func (m *Map[K, V]) Clear() {
	if m == nil {
		return
	}
	m.l.Lock()
	defer m.l.Unlock()
	clear(m.m)
}
