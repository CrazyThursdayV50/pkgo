package collector

import "github.com/CrazyThursdayV50/pkgo/builtin/slice"

func Map[E any, K comparable, V any](sli []E, mapper func(int, E) (bool, K, V)) (m map[K]V) {
	slice.From(sli...).Iter(func(i int, v E) (bool, error) {
		ok, key, val := mapper(i, v)
		if !ok {
			return true, nil
		}

		if m == nil {
			m = make(map[K]V)
		}

		m[key] = val
		return true, nil
	})
	return
}
