package utils

type Cloneable[T any] interface {
	Clone() T
}

func CloneTrait[T any](x Cloneable[T]) T {
	return x.Clone()
}
