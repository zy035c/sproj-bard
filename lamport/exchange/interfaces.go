package exchange

type Exchange[T any] interface {
	Put(t T) bool
	Get() T
	C() <-chan T
	GetId() uint64
}
