package main

func Map[T, K any](vals []T, transform func(val T) K) []K {
	res := []K{}

	for _, elem := range vals {
		res = append(res, transform(elem))
	}

	return res
}

func Join[T any](vals []T, sep T) []T {
	res := []T{}

	for idx, val := range vals {
		res = append(res, val)
		if idx != len(vals)-1 {
			res = append(res, sep)
		}
	}

	return res
}

func Reduce[T, K any](vals []T, reduce func(prev K, curr T, idx int, arr []T) K, startVal K) K {
	res := startVal

	for idx, elem := range vals {
		res = reduce(res, elem, idx, vals)
	}

	return res
}

func Expand[T any](val T, num int) []T {
	res := []T{}

	for i := 0; i < num; i++ {
		res = append(res, val)
	}

	return res
}
