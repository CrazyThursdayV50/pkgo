package collector

import "github.com/CrazyThursdayV50/pkgo/builtin/slice"

func Slice[E any, T any](sli []E, collector func(E) (bool, T)) (list []T) {
	_, _ = slice.From(sli...).Iter(func(_ int, v E) (bool, error) {
		ok, e := collector(v)
		if !ok {
			return true, nil
		}

		list = append(list, e)
		return true, nil
	})
	return
}
