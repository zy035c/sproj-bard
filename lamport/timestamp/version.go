package timestamp

import (
	"fmt"
)

type Version[T any, K any, U DistributedClock[K]] struct {
	data      T
	timestamp U
	id        uint64
}

func NewVersion[T any, K any, U DistributedClock[K]](
	data T, timestamp U, id uint64,
) *Version[T, K, U] {
	return &Version[T, K, U]{
		data:      data,
		timestamp: timestamp,
		id:        id,
	}
}

func (v Version[T, K, U]) String() string {
	return fmt.Sprintf("Version{data: %v, Ts: %v, Source: %v}", v.data, v.timestamp, v.id)
}

func (v Version[T, K, U]) GetTs() DistributedClock[K] {
	return v.timestamp
}

func (v Version[T, K, U]) GetData() T {
	return v.data
}

func (v Version[T, K, U]) GetId() uint64 {
	return v.id
}

func (v *Version[T, K, U]) Clone() *Version[T, K, U] {
	return &Version[T, K, U]{
		data:      v.GetData(),
		timestamp: v.GetTs().Clone().(U),
		id:        v.GetId(),
	}
}
