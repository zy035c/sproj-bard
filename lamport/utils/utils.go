package utils

import (
	"fmt"
	"reflect"
	"strconv"
)

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

func ReflectConvert[T any](val reflect.Value) *T {
	if res, ok := reflect.ValueOf(val).Interface().(T); ok {
		fmt.Printf("Successfully casted")
		return &res
	} else {
		return nil
	}
}
