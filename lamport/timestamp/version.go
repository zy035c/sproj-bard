package timestamp

import (
	"fmt"
)

type Version[T any, K any] struct {
	data      T
	timestamp DistributedClock[K]
	id        uint64
}

func NewVersion[T any, K any](
	data T, timestamp DistributedClock[K], id uint64,
) Version[T, K] {
	return Version[T, K]{
		data:      data,
		timestamp: timestamp,
		id:        id,
	}
}

// type Message[T Payload, K any] interface {
// 	String() string
// 	GetTs() timestamp.DistributedClock[K]
// 	GetData() T
// 	GetId() uint64
// }

func (v Version[T, K]) String() string {
	return fmt.Sprintf("Version{data: %v, Ts: %v, Source: %v}", v.data, v.timestamp, v.id)
}

func (v Version[T, K]) GetTs() DistributedClock[K] {
	return v.timestamp
}

func (v Version[T, K]) GetData() T {
	return v.data
}

func (v Version[T, K]) GetId() uint64 {
	return v.id
}

func (v Version[T, K]) Clone() Version[T, K] {
	return Version[T, K]{
		data:      v.GetData(),
		timestamp: v.GetTs().Clone().(DistributedClock[K]),
		id:        v.GetId(),
	}
}
