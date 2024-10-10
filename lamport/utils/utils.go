package utils

import "strconv"

type MyError struct {
	message string
}

// custom error
func (err MyError) Error() string {
	return err.message
}

func ParseIntOrPanic(s string) int {
	num, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return num
}

func SliceCpy[T any](src []T) []T {
	cpy := make([]T, len(src))
	copy(cpy, src)
	return cpy
}
