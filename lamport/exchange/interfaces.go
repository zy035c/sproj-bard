package exchange

type Exchange[T any] interface {
	Put(t T)
	Get() T
}
