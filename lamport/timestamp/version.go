package timestamp

import "fmt"

type Version[T any, K any, U DistributedClock[K]] struct {
	data      *T
	timestamp U
}

func (v Version[T, K, U]) String() string {
	return fmt.Sprintf("Version{data: %v, Ts: %v}", *v.data, v.timestamp)
}

func (v Version[T, K, U]) getClock() K {
	return v.timestamp.Value()
}
