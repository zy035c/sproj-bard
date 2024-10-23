package option

type Option[T any] struct {
	val T
	has bool
}

func Some[T any](elem T) *Option[T] {
	return &Option[T]{val: elem, has: true}
}

func None[T any]() *Option[T] {
	return &Option[T]{}
}

func (opt *Option[T]) Has() bool {
	return opt.has
}

func (opt *Option[T]) Unwrap() T {
	return opt.val
}
