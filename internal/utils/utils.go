package utils

import "math/rand"

func RandFrom[T any](arr []T) T {
	return arr[rand.Intn(len(arr))]
}
