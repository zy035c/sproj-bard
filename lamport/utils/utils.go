package utils

import (
	"fmt"
	"lamport/option"
	"math/rand"
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

func ReflectConvert[T any](val reflect.Value) *option.Option[T] {
	if res, ok := reflect.ValueOf(val).Interface().(T); ok {
		fmt.Printf("Successfully casted")
		return option.Some[T](res)
	} else {
		return option.None[T]()
	}
}

func RandomInt64(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomFloat32(min, max float32) float32 {
	return min + rand.Float32()*(max-min)
}
