package impl

import (
	"github.com/CrazyThursdayV50/gotils/pkg/wrapper/wrap"
	"github.com/CrazyThursdayV50/pkgo/builtin"
)

func (m *Map[K, V]) Set(k K, v V) {
	if m == nil {
		return
	}
	m.l.Lock()
	defer m.l.Unlock()
	m.m[k] = v
}

func (m *Map[K, V]) Get(k K) builtin.UnWrapper[V] {
	if m == nil {
		return wrap.Nil[V]()
	}
	m.l.RLock()
	defer m.l.RUnlock()
	v, ok := m.m[k]
	if !ok {
		return wrap.Nil[V]()
	}
	return wrap.Wrap(v)
}
