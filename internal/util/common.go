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
