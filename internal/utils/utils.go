package utils

import "math/rand"

func RandFrom[T any](arr []T) T {
	return arr[rand.Intn(len(arr))]
}

type Weighted interface {
	Weight() int
}

func WeightedRandFrom[T Weighted](arr []T) T {
	var totalWeight int
	for _, v := range arr {
		totalWeight += v.Weight()
	}

	r := rand.Intn(totalWeight)
	for _, v := range arr {
		if r < v.Weight() {
			return v
		}
		r -= v.Weight()
	}

	panic("WeightedRandFrom: unreachable")
}
