package collector

import "github.com/CrazyThursdayV50/pkgo/builtin/slice"

func Slice[E any, T any](sli []E, collector func(int, E) (bool, T)) (list []T) {
	_, _ = slice.From(sli...).Iter(func(i int, v E) (bool, error) {
		ok, e := collector(i, v)
		if !ok {
			return true, nil
		}

		list = append(list, e)
		return true, nil
	})
	return
}
