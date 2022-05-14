package util

func Set[T comparable](arr []T) []T {
	res := make([]T, 0)
	for _, e := range arr {
		flag := false
		for _, r := range res {
			if r == e {
				flag = true
			}
		}

		if !flag {
			res = append(res, e)
		}
	}

	return res
}

func MergeMap[K, V comparable](m1, m2 map[K]V) map[K]V {
	res := make(map[K]V)
	for k, v := range m1 {
		res[k] = v
	}

	for k, v := range m2 {
		res[k] = v
	}

	return res
}
