package utils

import (
	"math/rand"

	"github.com/samber/lo"
)

func RandFrom[T any](arr []T) T {
	return arr[rand.Intn(len(arr))]
}

type Weighted interface {
	Weight() int
}

func WeightedRandFrom[T Weighted](arr []T) T {
	var totalWeight int
	filtered := lo.Filter(arr, weightAboveZero[T])
	for _, v := range filtered {
		totalWeight += v.Weight()
	}

	r := rand.Intn(totalWeight)
	for _, v := range filtered {
		if r < v.Weight() {
			return v
		}
		r -= v.Weight()
	}

	panic("WeightedRandFrom: unreachable")
}

func weightAboveZero[T Weighted](item T, _ int) bool {
	return item.Weight() > 0
}
