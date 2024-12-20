package gmap

import (
	"github.com/CrazyThursdayV50/pkgo/builtin"
	"github.com/CrazyThursdayV50/pkgo/builtin/map/impl"
)

func From[K comparable, V any](m map[K]V) builtin.MapAPI[K, V] {
	return impl.From(m)
}

func Make[K comparable, V any](cap int) builtin.MapAPI[K, V] {
	return From(make(map[K]V, cap))
}
